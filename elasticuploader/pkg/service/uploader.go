package service

import (
	"bytes"
	"context"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
)

const numWorkers = 4
const flushBytes = 1000000

// BulkIndexerUpload uploads data to Elasticsearch using BulkIndexer from go-elasticsearch.
func BulkIndexerUpload(lines []models.Message,
	currentDocumentID int,
	indexName string,
	esClient *elasticsearch.Client) int {
	var (
		countSuccessful uint64

		err error
	)

	documentID := currentDocumentID
	flushInterval := 30 * time.Second

	// Create the BulkIndexer
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,     // The default index name
		Client:        esClient,      // The Elasticsearch client
		NumWorkers:    numWorkers,    // The number of worker goroutines
		FlushBytes:    flushBytes,    // The flush threshold in bytes
		FlushInterval: flushInterval, // The periodic flush interval
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

				// DocumentID is the (optional) document ID
				DocumentID: strconv.Itoa(documentID),

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

		documentID++

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

	return documentID
}

// CreateEsIndex creates an ES index.
func CreateEsIndex(index string, esClient *elasticsearch.Client) {
	var (
		res *esapi.Response
		err error
	)

	log.Println("Deleting index:  ", index, "...")

	// Re-create the index
	if res, err = esClient.Indices.Delete(
		[]string{index},
		esClient.Indices.Delete.WithIgnoreUnavailable(true)); err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)
	}

	res.Body.Close()

	log.Println("Creating index:  ", index, "...")
	res, err = esClient.Indices.Create(index)
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
