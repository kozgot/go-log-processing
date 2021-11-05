package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/service"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

const processedDataExchangeName = "processeddata_direct_durable"

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

	conn, err := amqp.Dial(rabbitMqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		processedDataExchangeName, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"save_data_queue_durable", // name
		true,                      // durable
		false,                     // delete when unused
		true,                      // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// TODO: extract routing key to a single place, eg.: env variables
	err = ch.QueueBind(
		q.Name,                    // queue name
		"save-data",               // routing key
		processedDataExchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	linesByIndexNames := make(map[string][]models.Message)
	forever := make(chan bool)
	dataCountTreshold := 1000
	documentID := 1
	uploadTicker := time.NewTicker(10 * time.Second)
	mutex := sync.Mutex{}

	go func() {
		for d := range msgs {
			msgParts := strings.Split(string(d.Body), "|")
			msgPrefix := msgParts[0]
			switch msgPrefix {
			case "CREATEINDEX":
				indexName := strings.Split(string(d.Body), "|")[1]
				service.CreateEsIndex(indexName, elasticSearchClient)
			default:
				// save messages until we hit the 1000 line treshold
				data := models.DeserializeDataUnit(d.Body)

				mutex.Lock()
				_, ok := linesByIndexNames[data.IndexName]
				if !ok {
					linesByIndexNames[data.IndexName] = []models.Message{}
				}

				linesByIndexNames[data.IndexName] = append(linesByIndexNames[data.IndexName], models.Message{Content: data.Data})

				// If we hit the treshold, we upload to ES.
				if len(linesByIndexNames[data.IndexName]) >= dataCountTreshold {
					uploadTicker.Reset(10 * time.Second)
					fmt.Println("Resetting ticker")

					documentID = service.BulkIndexerUpload(
						linesByIndexNames[data.IndexName],
						documentID,
						data.IndexName,
						elasticSearchClient)

					linesByIndexNames[data.IndexName] = []models.Message{} // clear the buffer after uploading the contents
				}

				mutex.Unlock()
			}

			// Acknowledge message
			err := d.Ack(false)
			failOnError(err, "Could not acknowledge message")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// Periodically check if we have anything left to upload.
	go func() {
		for range uploadTicker.C {
			mutex.Lock()
			for indexName := range linesByIndexNames {
				if len(linesByIndexNames[indexName]) > 0 {
					fmt.Println("Uploading leftovers after timeout into index " + indexName)
					documentID = service.BulkIndexerUpload(
						linesByIndexNames[indexName],
						documentID,
						indexName,
						elasticSearchClient)

					linesByIndexNames[indexName] = []models.Message{} // clear the buffer after uploading the contents
				}
			}
			mutex.Unlock()
		}
	}()

	<-forever
}
