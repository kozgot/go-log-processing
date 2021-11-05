package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/rabbit"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/service"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	log.Println("Elastic Uploader starting...")

	// Create the Elasticsearch client
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
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("RABBIT_URL:", rabbitMqURL)
	if len(rabbitMqURL) == 0 {
		log.Fatal("The RABBIT_URL environment variable is not set")
	}

	processedDataExchangeName := os.Getenv("PROCESSED_DATA_EXCHANGE")
	fmt.Println("PROCESSED_DATA_EXCHANGE:", processedDataExchangeName)
	if len(processedDataExchangeName) == 0 {
		log.Fatal("The PROCESSED_DATA_EXCHANGE environment variable is not set")
	}

	saveDataQueueName := os.Getenv("SAVE_DATA_QUEUE")
	fmt.Println("SAVE_DATA_QUEUE:", saveDataQueueName)
	if len(saveDataQueueName) == 0 {
		log.Fatal("The SAVE_DATA_QUEUE environment variable is not set")
	}

	saveDataRoutingKey := os.Getenv("SAVE_DATA_ROUTING_KEY")
	fmt.Println("SAVE_DATA_ROUTING_KEY:", saveDataRoutingKey)
	if len(saveDataRoutingKey) == 0 {
		log.Fatal("The SAVE_DATA_ROUTING_KEY environment variable is not set")
	}

	rabbitMQConsumer := rabbit.AmqpConsumer{HostDsn: rabbitMqURL}
	err = rabbitMQConsumer.Connect()
	failOnError(err, "Could not connect ro RabbitMQ")
	defer rabbitMQConsumer.CloseConnection()

	err = rabbitMQConsumer.Channel()
	failOnError(err, "Could not open channel")
	defer rabbitMQConsumer.CloseChannel()

	msgs, err := rabbitMQConsumer.Consume(processedDataExchangeName, saveDataQueueName, saveDataRoutingKey)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	handleMessages(msgs, elasticSearchClient)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}

func handleMessages(deliveries <-chan amqp.Delivery, esClient *elasticsearch.Client) {
	uploadTicker := time.NewTicker(10 * time.Second)

	uploadBuffer := service.InitUploadBuffer(esClient)

	// Periodically check if we have anything left to upload.
	go func() {
		for range uploadTicker.C {
			uploadBuffer.UploadRemaining()
		}
	}()

	go func() {
		for delivery := range deliveries {
			msgParts := strings.Split(string(delivery.Body), "|")
			msgPrefix := msgParts[0]
			switch msgPrefix {
			case "CREATEINDEX":
				indexName := strings.Split(string(delivery.Body), "|")[1]
				service.CreateEsIndex(indexName, esClient)
			default:
				data := models.DeserializeDataUnit(delivery.Body)
				uploadBuffer.AppendAndUploadIfNeeded(models.Message{Content: data.Data}, data.IndexName, uploadTicker)
			}

			// Acknowledge message
			err := delivery.Ack(false)
			failOnError(err, "Could not acknowledge message")
		}
	}()
}
