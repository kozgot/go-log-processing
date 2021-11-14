package elastic

import (
	"bytes"
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
)

// EsClientWrapper implements the Esclient interface used by the uploader service.
type EsClientWrapper struct {
	esClient *elasticsearch.Client
}

// NewEsClientWrapper creates a new EsClientWrapper.
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
	utils.FailOnError(err, " [ESClient] Error creating the client")
	clientWrapper := EsClientWrapper{esClient: elasticSearchClient}

	return &clientWrapper
}

// BulkUpload performs a bulk indexing for the given array of data units.
func (esuploader *EsClientWrapper) BulkUpload(dataUnits []models.DataUnit, indexName string) {
	var (
		countSuccessful uint64
		err             error
	)

	// Create the BulkIndexer.
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,           // The default index name
		Client:        esuploader.esClient, // The Elasticsearch client
		NumWorkers:    4,                   // The number of worker goroutines
		FlushBytes:    1000000,             // The flush threshold in bytes
		FlushInterval: 30 * time.Second,    // The periodic flush interval
	})
	utils.FailOnError(err, " [ESClient] Error creating the indexer")

	start := time.Now().UTC()

	for _, dataUnit := range dataUnits {
		// Add a data unit to the BulkIndexer.
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "index",
				Body:   bytes.NewReader(dataUnit.Content),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		utils.FailOnError(err, " [ESClient] Unexpected error")
	}

	// Close the indexer
	err = bi.Close(context.Background())
	utils.FailOnError(err, " [ESClient] Unexpected error")

	dur := time.Since(start)
	biStats := bi.Stats()

	logBulkIndexerStats(biStats, dur)
}

// RecreateEsIndex deletes an ES index if it exists, and then recreates it.
func (esuploader *EsClientWrapper) RecreateEsIndex(index string) {
	var (
		res *esapi.Response
		err error
	)

	log.Println(" [ESClient] Deleting index:  ", index, "...")

	// Re-create the index
	if res, err = esuploader.esClient.Indices.Delete(
		[]string{index},
		esuploader.esClient.Indices.Delete.WithIgnoreUnavailable(true)); err != nil || res.IsError() {
		log.Fatalf(" [ESClient] Cannot delete index: %s", err)
	}

	res.Body.Close()

	log.Println(" [ESClient] Creating index:  ", index, "...")
	res, err = esuploader.esClient.Indices.Create(index)
	if err != nil {
		log.Fatalf(" [ESClient] Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf(" [ESClient] Cannot create index: %s", res)
	}
	res.Body.Close()
}

func logBulkIndexerStats(biStats esutil.BulkIndexerStats, dur time.Duration) {
	if biStats.NumFailed > 0 {
		// We got some errors while trying to index the documents
		log.Fatalf(
			" [ESClient] Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		// Indexed evereything successfully
		log.Printf(
			" [ESClient] Successfully indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}
}
