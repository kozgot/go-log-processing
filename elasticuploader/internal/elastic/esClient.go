package elastic

import (
	"bytes"
	"context"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
)

type EsClient interface {
	BulkUpload(lines []models.Message, indexName string)
	CreateEsIndex(index string)
}

type EsClientWrapper struct {
	esClient *elasticsearch.Client
}

func NewEsClientWrapper() *EsClientWrapper {
	// Create the ES client
	// Use a third-party package for implementing the backoff function
	retryBackoff := backoff.NewExponentialBackOff()
	elasticSearchClient, err := elasticsearch.NewClient(elasticsearch.Config{
		// Retry on 429 TooManyRequests statuses
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 10,
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	result := EsClientWrapper{esClient: elasticSearchClient}

	return &result
}

const numWorkers = 4
const flushBytes = 1000000

func (esuploader *EsClientWrapper) BulkUpload(lines []models.Message, indexName string) {
	var (
		countSuccessful uint64
		err             error
	)

	flushInterval := 30 * time.Second

	// Create the BulkIndexer
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,           // The default index name
		Client:        esuploader.esClient, // The Elasticsearch client
		NumWorkers:    numWorkers,          // The number of worker goroutines
		FlushBytes:    flushBytes,          // The flush threshold in bytes
		FlushInterval: flushInterval,       // The periodic flush interval
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	start := time.Now().UTC()

	for _, a := range lines {
		// Add an item to the BulkIndexer
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",

				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(a.Content),

				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}

	// Close the indexer
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	dur := time.Since(start)
	biStats := bi.Stats()

	reportBulkIndexerStats(biStats, dur)
}

// CreateEsIndex creates an ES index.
func (esuploader *EsClientWrapper) CreateEsIndex(index string) {
	var (
		res *esapi.Response
		err error
	)

	log.Println("Deleting index:  ", index, "...")

	// Re-create the index
	if res, err = esuploader.esClient.Indices.Delete(
		[]string{index},
		esuploader.esClient.Indices.Delete.WithIgnoreUnavailable(true)); err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)
	}

	res.Body.Close()

	log.Println("Creating index:  ", index, "...")
	res, err = esuploader.esClient.Indices.Create(index)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}
	res.Body.Close()
}

func reportBulkIndexerStats(biStats esutil.BulkIndexerStats, dur time.Duration) {
	if biStats.NumFailed > 0 {
		// We got some errors while trying to index the documents
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		// Indexed evereything successfully
		log.Printf(
			"Successfully indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}

	log.Println(strings.Repeat("▔", 65))
}
