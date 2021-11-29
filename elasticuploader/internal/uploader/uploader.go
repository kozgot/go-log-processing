package uploader

import (
	"github.com/kozgot/go-log-processing/elasticuploader/internal/elastic"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	postprocmodels "github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// UploaderService encapsultes the data and logic of the uploading service.
type UploaderService struct {
	rabbitMQConsumer        rabbitmq.MessageConsumer
	esClient                elastic.EsClient
	eventIndexName          string
	consumptionIndexName    string
	indexRecreationTimeSpec string
}

// NewUploaderService creates a new uploader service instance.
// The indexRecreationTimeSpec is used to time the creation of new ES indexes,
// see the docs of github.com/robfig/cron/v3 for the syntax.
func NewUploaderService(
	messageConsumer rabbitmq.MessageConsumer,
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
			// Deserialize the received data.
			data := postprocmodels.DataUnit{}
			data.Deserialize(delivery.Body)

			// Append it to the buffer.
			uploadBuffer.AppendAndUploadIfNeeded(
				models.ESDocument{Content: data.Data},
				data.DataType,
			)

			// Acknowledge message.
			err := delivery.Ack(false)
			utils.FailOnError(err, " [UPLOADER SERVICE] Could not acknowledge message")
		}
	}()
}
