package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/internal/processing"
	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"
	"github.com/streadway/amqp"
)

const consumptionIndexName = "consumption"
const smcIndexName = "smc"

func main() {
	log.Println("PostProcessor service starting...")
	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("RABBIT_URL:", rabbitMqURL)
	if len(rabbitMqURL) == 0 {
		log.Fatal("The RABBIT_URL environment variable is not set")
	}

	saveDataExchangeName := os.Getenv("PROCESSED_DATA_EXCHANGE")
	fmt.Println("PROCESSED_DATA_EXCHANGE:", saveDataExchangeName)
	if len(saveDataExchangeName) == 0 {
		log.Fatal("The PROCESSED_DATA_EXCHANGE environment variable is not set")
	}

	saveDataRoutingKey := os.Getenv("SAVE_DATA_ROUTING_KEY")
	fmt.Println("SAVE_DATA_ROUTING_KEY:", saveDataRoutingKey)
	if len(saveDataRoutingKey) == 0 {
		log.Fatal("The SAVE_DATA_ROUTING_KEY environment variable is not set")
	}

	processEntriesExchangeName := os.Getenv("LOG_ENTRIES_EXCHANGE")
	fmt.Println("LOG_ENTRIES_EXCHANGE:", processEntriesExchangeName)
	if len(processEntriesExchangeName) == 0 {
		log.Fatal("The LOG_ENTRIES_EXCHANGE environment variable is not set")
	}

	processingQueueName := os.Getenv("PROCESSING_QUEUE")
	fmt.Println("PROCESSING_QUEUE:", processingQueueName)
	if len(processingQueueName) == 0 {
		log.Fatal("The PROCESSING_QUEUE environment variable is not set")
	}

	processEntryRoutingKey := os.Getenv("PROCESS_ENTRY_ROUTING_KEY")
	fmt.Println("PROCESS_ENTRY_ROUTING_KEY:", processEntryRoutingKey)
	if len(processEntryRoutingKey) == 0 {
		log.Fatal("The PROCESS_ENTRY_ROUTING_KEY environment variable is not set")
	}

	rabbitMQConsumer := rabbitmq.AmqpConsumer{HostDsn: rabbitMqURL}
	err := rabbitMQConsumer.Connect()
	utils.FailOnError(err, "Could not connect ro RabbitMQ")
	defer rabbitMQConsumer.CloseConnection()

	err = rabbitMQConsumer.Channel()
	utils.FailOnError(err, "Could not open channel")
	defer rabbitMQConsumer.CloseChannel()

	msgs, err := rabbitMQConsumer.Consume(processEntriesExchangeName, processingQueueName, processEntryRoutingKey)
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	handleEntries(msgs, rabbitMqURL, saveDataExchangeName, saveDataRoutingKey)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func deserializeMessage(message []byte) parsermodels.ParsedLogEntry {
	var data parsermodels.ParsedLogEntry
	if err := json.Unmarshal(message, &data); err != nil {
		fmt.Println("Failed to unmarshal: ", err)
	}

	return data
}

func handleEntries(
	deliveries <-chan amqp.Delivery,
	rabbitMqURL string,
	saveDataExchangeName string,
	saveDataRoutingKey string) {
	go func() {
		esUploader := rabbitmq.EsUploadSender{
			RabbitMqURL:  rabbitMqURL,
			ExchangeName: saveDataExchangeName,
			RoutingKey:   saveDataRoutingKey}

		// Open rabbitmq channel and connection.
		esUploader.OpenChannelAndConnection(rabbitMqURL)
		defer esUploader.CloseChannelAndConnection()

		// Create indices in ES.
		esUploader.CreateIndex(smcIndexName)
		esUploader.CreateIndex(consumptionIndexName)

		processor := processing.InitEntryProcessor(&esUploader)

		for d := range deliveries {
			if strings.Contains(string(d.Body), "END") {
				fmt.Println("End of entries...")

				// Further processing to get consumption and index info.
				processor.ProcessConsumptionAndIndexValues(consumptionIndexName)

				// Acknowledge the message after it has been processed.
				err := d.Ack(false)
				utils.FailOnError(err, "Could not acknowledge END message")
				continue
			}

			entry := deserializeMessage(d.Body)
			processor.Process(entry, smcIndexName)

			// Acknowledge the message after it has been processed.
			err := d.Ack(false)
			utils.FailOnError(err,
				"Could not acknowledge message with timestamp: "+entry.Timestamp.Format("2 Jan 2006 15:04:05"))
		}
	}()
}
