package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
)

// UploadBuffer stores data by index name until the datacount reaches a treshold,
// then uploads the contents, while implementing mutual exclosure.
type UploadBuffer struct {
	mutex    sync.Mutex
	value    map[string][]models.Message
	esClient *elasticsearch.Client
}

// InitUploadBuffer initializes the buffer.
func InitUploadBuffer(esClient *elasticsearch.Client) *UploadBuffer {
	return &UploadBuffer{value: make(map[string][]models.Message)}
}

// AppendAndUploadIfNeeded appends a message for the given key.
func (d *UploadBuffer) AppendAndUploadIfNeeded(m models.Message, key string, uploadTicker *time.Ticker) {
	d.mutex.Lock() // Lock so only one goroutine at a time can access the map.
	defer d.mutex.Unlock()

	// Check if the key is already present.
	_, ok := d.value[key]
	if !ok {
		d.value[key] = []models.Message{}
	}

	d.value[key] = append(d.value[key], m)

	// If we hit the treshold, we upload to ES.
	if len(d.value[key]) >= 1000 {
		uploadTicker.Reset(10 * time.Second)
		fmt.Println("Resetting ticker")

		bulkIndexerUpload(
			d.value[key],
			key,
			d.esClient)

		// Clear
		d.value[key] = []models.Message{}
	}
}

// GetCurrentMessages returns the current messages for a given key.
func (d *UploadBuffer) GetCurrentMessages(key string) []models.Message {
	d.mutex.Lock() // Lock so only one goroutine at a time can access the map.
	defer d.mutex.Unlock()
	return d.value[key]
}

// UploadRemaining uploads the data left in the buffer and clears the buffer.
func (d *UploadBuffer) UploadRemaining() {
	d.mutex.Lock() // Lock so only one goroutine at a time can access the map.
	defer d.mutex.Unlock()

	for indexName := range d.value {
		if len(d.value[indexName]) > 0 {
			fmt.Println("Uploading leftovers after timeout into index " + indexName)
			bulkIndexerUpload(
				d.value[indexName],
				indexName,
				d.esClient)

			// Clear the buffer after uploading the contents.
			d.value[indexName] = []models.Message{}
		}
	}
}

const numWorkers = 4
const flushBytes = 1000000

func bulkIndexerUpload(lines []models.Message, indexName string, esClient *elasticsearch.Client) {
	var (
		countSuccessful uint64
		err             error
	)

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
