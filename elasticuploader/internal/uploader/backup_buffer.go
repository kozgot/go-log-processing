package uploader

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	postprocmodels "github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

type Backup struct {
	EventDocuments       []models.ESDocument
	ConsumptionDocuments []models.ESDocument

	EventIndexName       string
	ConsumptionIndexName string
}

const backupFileName = "unsaved_docs_backup.json"

func (backup *Backup) ToJSON() []byte {
	bytes, err := json.MarshalIndent(backup, "", " ")
	utils.FailOnError(err, "Can't serialize backup documents")

	return bytes
}

func (backup *Backup) FromJSON(bytes []byte) {
	err := json.Unmarshal(bytes, backup)
	utils.FailOnError(err, "Failed to unmarshal backup documents")
}

type BackupBuffer struct {
	backup             Backup
	consumptionsBuffer []models.ESDocument
	eventsBuffer       []models.ESDocument
}

func NewBackupBuffer() *BackupBuffer {
	buffer := BackupBuffer{
		backup: Backup{
			EventDocuments:       []models.ESDocument{},
			ConsumptionDocuments: []models.ESDocument{},
		},
		consumptionsBuffer: []models.ESDocument{},
		eventsBuffer:       []models.ESDocument{},
	}

	return &buffer
}

func (b *BackupBuffer) Add(unSavedDoc models.ESDocument, dataType postprocmodels.DataType) {
	switch dataType {
	case postprocmodels.Consumption:
		b.consumptionsBuffer = append(b.consumptionsBuffer, unSavedDoc)
	case postprocmodels.Event:
		b.eventsBuffer = append(b.eventsBuffer, unSavedDoc)
	case postprocmodels.UnknownDataType:
		break
	default:
		break
	}

	if len(b.consumptionsBuffer)+len(b.eventsBuffer) == 10 {
		b.save()
		b.consumptionsBuffer = []models.ESDocument{}
		b.eventsBuffer = []models.ESDocument{}
	}
}

func (b *BackupBuffer) save() {
	b.backup.EventDocuments = append(b.backup.EventDocuments, b.eventsBuffer...)
	b.backup.ConsumptionDocuments = append(b.backup.ConsumptionDocuments, b.consumptionsBuffer...)

	serialized := b.backup.ToJSON()
	_ = ioutil.WriteFile(backupFileName, serialized, 0600)
}

// Reset resets the backup buffer and clears the backup file contents.
func (b *BackupBuffer) Reset() {
	b.backup.EventDocuments = []models.ESDocument{}
	b.backup.ConsumptionDocuments = []models.ESDocument{}

	b.consumptionsBuffer = []models.ESDocument{}
	b.eventsBuffer = []models.ESDocument{}

	serialized := b.backup.ToJSON()
	_ = ioutil.WriteFile(backupFileName, serialized, 0600)
}

// Clear clears the documents of the given type from the backup buffer.
func (b *BackupBuffer) Clear(dataType postprocmodels.DataType) {
	switch dataType {
	case postprocmodels.Consumption:
		b.backup.ConsumptionDocuments = []models.ESDocument{}
		b.consumptionsBuffer = []models.ESDocument{}
	case postprocmodels.Event:
		b.backup.EventDocuments = []models.ESDocument{}
		b.eventsBuffer = []models.ESDocument{}
	case postprocmodels.UnknownDataType:
		break
	default:
		break
	}

	serialized := b.backup.ToJSON()
	_ = ioutil.WriteFile(backupFileName, serialized, 0600)
}

func (b *BackupBuffer) Load() ([]models.ESDocument, []models.ESDocument) {
	if _, err := os.Stat(backupFileName); errors.Is(err, os.ErrNotExist) {
		// no backup file
		return []models.ESDocument{}, []models.ESDocument{}
	}

	serializedBackup, err := ioutil.ReadFile(backupFileName)
	utils.FailOnError(err, "Could not read unsaved data backup file")
	b.backup.FromJSON(serializedBackup)
	return b.backup.EventDocuments, b.backup.ConsumptionDocuments
}

func (b *BackupBuffer) GetBackupIndexNames() (string, string) {
	return b.backup.EventIndexName, b.backup.ConsumptionIndexName
}

func (b *BackupBuffer) SetIndexNames(eventIndex string, consumptionIndex string) {
	b.backup.EventIndexName = eventIndex
	b.backup.ConsumptionIndexName = consumptionIndex
	b.save()
}
