package uploader

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/elastic"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	postprocmodels "github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/robfig/cron/v3"
)

// UploadBuffer stores data by index name until the datacount reaches a treshold,
// then uploads the contents, while implementing mutual exclosure.
type UploadBuffer struct {
	mutex                sync.Mutex
	value                map[string][]models.ESDocument
	esClient             elastic.EsClient
	ticker               *time.Ticker
	bufferSize           int
	eventIndexName       string
	consumptionIndexName string
	indexPostFix         string
	backupBuffer         *BackupBuffer
}

// NewUploadBuffer initializes the buffer.
func NewUploadBuffer(
	esClient elastic.EsClient,
	size int,
	eventIndexName string,
	consumptionIndexName string,
	indexRecreationTimeSpec string,
) *UploadBuffer {
	ticker := time.NewTicker(5 * time.Second)
	backupBuffer := NewBackupBuffer()
	uploadBuffer := UploadBuffer{
		value:                make(map[string][]models.ESDocument),
		esClient:             esClient,
		ticker:               ticker,
		bufferSize:           size,
		eventIndexName:       eventIndexName,
		consumptionIndexName: consumptionIndexName,
		indexPostFix:         createIndexPostFix(),
		backupBuffer:         backupBuffer,
	}

	uploadBuffer.uploadBackupIfNeeded()

	// Create ES indexes for the day.
	// This takes care of the index creation just after the service is started (it might not happen at midnight exactly).
	currentEventIndexName := uploadBuffer.postfixIndexName(eventIndexName)
	currentConsIndexName := uploadBuffer.postfixIndexName(consumptionIndexName)
	log.Println(" [UPLOADER SERVICE] Events index name: " + currentEventIndexName)
	log.Println(" [UPLOADER SERVICE] Consumptions index name: " + currentConsIndexName)
	uploadBuffer.esClient.CreateEsIndex(currentEventIndexName)
	uploadBuffer.esClient.CreateEsIndex(currentConsIndexName)
	log.Println(" [UPLOADER SERVICE] Created new indexes at startup")

	// Save initial index name to the backup file.
	uploadBuffer.backupBuffer.SetIndexNames(
		uploadBuffer.postfixIndexName(eventIndexName),
		uploadBuffer.postfixIndexName(consumptionIndexName))

	cronHandler := cron.New(cron.WithLocation(time.Local))
	// the 0/24th hour and 0th minute of every day
	_, err := cronHandler.AddFunc(indexRecreationTimeSpec, func() {
		uploadBuffer.mutex.Lock()

		log.Println(" [UPLOADER SERVICE] Switching index names for the next day...")

		// Upload remainings from the previous day.
		uploadBuffer.uploadAndClearBuffer()

		// Update index postfix.
		uploadBuffer.indexPostFix = createIndexPostFix()

		// Create ES indexes for the day, every day at midnight.
		currentEventIndexName := uploadBuffer.postfixIndexName(eventIndexName)
		currentConsIndexName := uploadBuffer.postfixIndexName(consumptionIndexName)
		log.Println(" [UPLOADER SERVICE] New events index name: " + currentEventIndexName)
		log.Println(" [UPLOADER SERVICE] New consumptions index name: " + currentConsIndexName)
		uploadBuffer.esClient.CreateEsIndex(currentConsIndexName)
		uploadBuffer.esClient.CreateEsIndex(currentEventIndexName)
		log.Println(" [UPLOADER SERVICE] Created new indexes")

		// Save current index name to the backup file.
		uploadBuffer.backupBuffer.SetIndexNames(
			uploadBuffer.postfixIndexName(eventIndexName),
			uploadBuffer.postfixIndexName(consumptionIndexName))

		uploadBuffer.mutex.Unlock()
	})
	utils.FailOnError(err, " [UPLOADER SERVICE] Failed to register cronhandler function")

	cronHandler.Start()

	// Periodically check if we have anything left to upload.
	go func() {
		for range ticker.C {
			uploadBuffer.UploadRemaining()
		}
	}()

	return &uploadBuffer
}

func (d *UploadBuffer) uploadBackupIfNeeded() {
	log.Println(" [UPLOADER SERVICE] Checking backup documents to upload...")
	unsavedEvents, unsavedConsumptions := d.backupBuffer.Load()
	eventsIndexName, consumptionsIndexName := d.backupBuffer.GetBackupIndexNames()
	foundBackupData := false
	if len(unsavedEvents) > 0 {
		log.Println(" [UPLOADER SERVICE] Found backup events to upload")
		d.esClient.BulkUpload(unsavedEvents, eventsIndexName)
		foundBackupData = true
	}

	if len(unsavedConsumptions) > 0 {
		log.Println(" [UPLOADER SERVICE] Found backup consumptions to upload")
		d.esClient.BulkUpload(unsavedConsumptions, consumptionsIndexName)
		foundBackupData = true
	}

	if !foundBackupData {
		log.Println(" [UPLOADER SERVICE] No backup documents to upload")
		return
	}

	d.backupBuffer.Reset()
}

// AppendAndUploadIfNeeded appends a message for the given key.
func (d *UploadBuffer) AppendAndUploadIfNeeded(m models.ESDocument, dataType postprocmodels.DataType) {
	// Lock so only one goroutine at a time can access the map.
	d.mutex.Lock()
	defer d.mutex.Unlock()

	indexName := d.eventIndexName
	if dataType == postprocmodels.Consumption {
		indexName = d.consumptionIndexName
	}

	// Check if the key is already present.
	_, ok := d.value[indexName]
	if !ok {
		d.value[indexName] = []models.ESDocument{}
	}

	// Append value to the buffer.
	d.value[indexName] = append(d.value[indexName], m)

	// Append value to the backup buffer until final upload.
	d.backupBuffer.Add(m, dataType)

	// If we hit the treshold, we upload to ES.
	if len(d.value[indexName]) >= d.bufferSize {
		d.ticker.Reset(5 * time.Second)

		// Upload to ES.
		d.esClient.BulkUpload(d.value[indexName], d.postfixIndexName(indexName))

		// Clear the buffer.
		d.value[indexName] = []models.ESDocument{}

		// Clear the backup after upload.
		d.backupBuffer.Clear(dataType)
	}
}

func (d *UploadBuffer) postfixIndexName(indexName string) string {
	return indexName + "_" + d.indexPostFix
}

func createIndexPostFix() string {
	currentDateTime := time.Now()

	return fmt.Sprint(currentDateTime.Year()) +
		fmt.Sprint(int(currentDateTime.Month())) +
		fmt.Sprint(currentDateTime.Day()) + "_" +
		fmt.Sprint(currentDateTime.Hour()) +
		fmt.Sprint(currentDateTime.Minute()) +
		fmt.Sprint(currentDateTime.Second())
}

// GetCurrentMessages returns the current messages for a given key.
func (d *UploadBuffer) GetCurrentMessages(key string) []models.ESDocument {
	// Lock so only one goroutine at a time can access the map.
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.value[key]
}

// UploadRemaining uploads the data left in the buffer and clears the buffer.
func (d *UploadBuffer) UploadRemaining() {
	// Lock so only one goroutine at a time can access the map.
	d.mutex.Lock()

	d.uploadAndClearBuffer()

	d.mutex.Unlock()
}

func (d *UploadBuffer) uploadAndClearBuffer() {
	resetNeeded := false

	// locking is handled from the outside
	for indexName := range d.value {
		if len(d.value[indexName]) > 0 {
			log.Println(" [UPLOADER SERVICE] Uploading leftovers into index " + d.postfixIndexName(indexName))
			d.esClient.BulkUpload(d.value[indexName], d.postfixIndexName(indexName))

			// Clear the buffer after uploading the contents.
			d.value[indexName] = []models.ESDocument{}
			resetNeeded = true
		}
	}

	if resetNeeded {
		// Clear the local backup of the documents after they have been uploaded.
		d.backupBuffer.Reset()
	}
}
