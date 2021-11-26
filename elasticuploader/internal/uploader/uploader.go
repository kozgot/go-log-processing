package uploader

import (
	"log"
	"strings"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/elastic"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/rabbit"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	postprocmodels "github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// UploaderService encapsultes the data and logic of the uploading service.
type UploaderService struct {
	rabbitMQConsumer        rabbit.MessageConsumer
	esClient                elastic.EsClient
	eventIndexName          string
	consumptionIndexName    string
	indexRecreationTimeSpec string
}

// NewUploaderService creates a new uploader service instance.
func NewUploaderService(
	messageConsumer rabbit.MessageConsumer,
	esClient elastic.EsClient,
	eventIndexName string,
	consumptionIndexName string,
	indexRecreationTimeSpec string,
) *UploaderService {
	service := UploaderService{
		rabbitMQConsumer:        messageConsumer,
		esClient:                esClient,
		eventIndexName:          eventIndexName,
		consumptionIndexName:    consumptionIndexName,
		indexRecreationTimeSpec: indexRecreationTimeSpec,
	}
	return &service
}

// HandleMessages consumes messages from rabbitMQ and uploads them to ES.
func (service *UploaderService) HandleMessages() {
	uploadBuffer := NewUploadBuffer(
		service.esClient,
		1000,
		service.eventIndexName,
		service.consumptionIndexName,
		service.indexRecreationTimeSpec,
	)

	msgs, err := service.rabbitMQConsumer.Consume()
	utils.FailOnError(err, " [UPLOADER SERVICE] Failed to register a consumer")

	go func() {
		for delivery := range msgs {
			msgParts := strings.Split(string(delivery.Body), "|")
			msgPrefix := msgParts[0]
			switch msgPrefix {
			case "DONE":
				log.Println(" [UPLOADER SERVICE] Received DONE from Postprocessor")
			default:
				data := postprocmodels.DataUnit{}
				data.Deserialize(delivery.Body)
				uploadBuffer.AppendAndUploadIfNeeded(
					models.DataUnit{Content: data.Data},
					data.DataType,
				)
			}

			// Acknowledge message
			err := delivery.Ack(false)
			utils.FailOnError(err, " [UPLOADER SERVICE] Could not acknowledge message")
		}
	}()
}
