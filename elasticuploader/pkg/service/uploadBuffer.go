package service

import (
	"log"
	"sync"
	"time"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/elastic"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
)

// UploadBuffer stores data by index name until the datacount reaches a treshold,
// then uploads the contents, while implementing mutual exclosure.
type UploadBuffer struct {
	mutex      sync.Mutex
	value      map[string][]models.DataUnit
	esClient   elastic.EsClient
	ticker     *time.Ticker
	bufferSize int
}

// NewUploadBuffer initializes the buffer.
func NewUploadBuffer(esClient elastic.EsClient, size int) *UploadBuffer {
	ticker := time.NewTicker(5 * time.Second)
	uploadBuffer := UploadBuffer{
		value:      make(map[string][]models.DataUnit),
		esClient:   esClient,
		ticker:     ticker,
		bufferSize: size,
	}

	// Periodically check if we have anything left to upload.
	go func() {
		for range ticker.C {
			uploadBuffer.UploadRemaining()
		}
	}()

	return &uploadBuffer
}

// AppendAndUploadIfNeeded appends a message for the given key.
func (d *UploadBuffer) AppendAndUploadIfNeeded(m models.DataUnit, key string) {
	// Lock so only one goroutine at a time can access the map.
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if the key is already present.
	_, ok := d.value[key]
	if !ok {
		d.value[key] = []models.DataUnit{}
	}

	d.value[key] = append(d.value[key], m)

	// If we hit the treshold, we upload to ES.
	if len(d.value[key]) >= d.bufferSize {
		d.ticker.Reset(5 * time.Second)

		// Upload to ES.
		d.esClient.BulkUpload(d.value[key], key)

		// Clear
		d.value[key] = []models.DataUnit{}
	}
}

// GetCurrentMessages returns the current messages for a given key.
func (d *UploadBuffer) GetCurrentMessages(key string) []models.DataUnit {
	// Lock so only one goroutine at a time can access the map.
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.value[key]
}

// UploadRemaining uploads the data left in the buffer and clears the buffer.
func (d *UploadBuffer) UploadRemaining() {
	// Lock so only one goroutine at a time can access the map.
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for indexName := range d.value {
		if len(d.value[indexName]) > 0 {
			log.Println("Uploading leftovers after timeout into index " + indexName)
			d.esClient.BulkUpload(d.value[indexName], indexName)

			// Clear the buffer after uploading the contents.
			d.value[indexName] = []models.DataUnit{}
		}
	}
}
