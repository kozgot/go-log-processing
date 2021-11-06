package service

import (
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/elastic"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	"github.com/streadway/amqp"
)

// MessageConsumer encapsulates messages needed to consume rabbitmq messages.
type MessageConsumer interface {
	Consume() (<-chan amqp.Delivery, error)
}

// UploaderService encapsultes the data and logic of the uploading service.
type UploaderService struct {
	rabbitMQConsumer MessageConsumer
	esClient         elastic.EsClient
}

// NewUploaderService creates a new uploader service instance.
func NewUploaderService(messageConsumer MessageConsumer, esClient elastic.EsClient) *UploaderService {
	service := UploaderService{
		rabbitMQConsumer: messageConsumer,
		esClient:         esClient,
	}
	return &service
}

// HandleMessages consumes messages from rabbitMQ and uploads them to ES.
func (service *UploaderService) HandleMessages() {
	uploadTicker := time.NewTicker(10 * time.Second)

	uploadBuffer := NewUploadBuffer(service.esClient)

	msgs, err := service.rabbitMQConsumer.Consume()
	utils.FailOnError(err, "Failed to register a consumer")

	// Periodically check if we have anything left to upload.
	go func() {
		for range uploadTicker.C {
			uploadBuffer.UploadRemaining()
		}
	}()

	go func() {
		for delivery := range msgs {
			msgParts := strings.Split(string(delivery.Body), "|")
			msgPrefix := msgParts[0]
			switch msgPrefix {
			case "CREATEINDEX":
				indexName := strings.Split(string(delivery.Body), "|")[1]
				service.esClient.CreateEsIndex(indexName)
			default:
				data := models.DeserializeDataUnit(delivery.Body)
				uploadBuffer.AppendAndUploadIfNeeded(models.Message{Content: data.Data}, data.IndexName, uploadTicker)
			}

			// Acknowledge message
			err := delivery.Ack(false)
			utils.FailOnError(err, "Could not acknowledge message")
		}
	}()
}
