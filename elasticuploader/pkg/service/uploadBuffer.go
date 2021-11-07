package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/elastic"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
)

// UploadBuffer stores data by index name until the datacount reaches a treshold,
// then uploads the contents, while implementing mutual exclosure.
type UploadBuffer struct {
	mutex    sync.Mutex
	value    map[string][]models.Message
	esClient elastic.EsClient
}

// NewUploadBuffer initializes the buffer.
func NewUploadBuffer(esClient elastic.EsClient) *UploadBuffer {
	return &UploadBuffer{value: make(map[string][]models.Message), esClient: esClient}
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

		// Upload to ES.
		d.esClient.BulkUpload(d.value[key], key)

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
			d.esClient.BulkUpload(d.value[indexName], indexName)

			// Clear the buffer after uploading the contents.
			d.value[indexName] = []models.Message{}
		}
	}
}
