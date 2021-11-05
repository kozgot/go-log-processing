package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/service"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"
)

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

	consumptionIndexName := os.Getenv("CONSUMPTION_INDEX_NAME")
	fmt.Println("CONSUMPTION_INDEX_NAME:", consumptionIndexName)
	if len(consumptionIndexName) == 0 {
		log.Fatal("The CONSUMPTION_INDEX_NAME environment variable is not set")
	}

	eventsIndexName := os.Getenv("EVENTS_INDEX_NAME")
	fmt.Println("EVENTS_INDEX_NAME:", eventsIndexName)
	if len(eventsIndexName) == 0 {
		log.Fatal("The EVENTS_INDEX_NAME environment variable is not set")
	}

	rabbitMQConsumer := rabbitmq.NewAmqpConsumer(
		rabbitMqURL,
		processEntryRoutingKey,
		processEntriesExchangeName,
		processingQueueName)

	err := rabbitMQConsumer.Connect()
	utils.FailOnError(err, "Could not connect ro RabbitMQ")
	defer rabbitMQConsumer.CloseConnection()

	err = rabbitMQConsumer.Channel()
	utils.FailOnError(err, "Could not open channel")
	defer rabbitMQConsumer.CloseChannel()

	forever := make(chan bool)

	esUploader := rabbitmq.NewEsUploadSender(
		rabbitMqURL,
		saveDataExchangeName,
		saveDataRoutingKey,
		eventsIndexName,
		consumptionIndexName)

	// Open rabbitmq channel and connection.
	esUploader.Connect(rabbitMqURL)
	defer esUploader.CloseChannelAndConnection()

	// Create indices in ES.
	esUploader.CreateIndex(eventsIndexName)
	esUploader.CreateIndex(consumptionIndexName)

	service.HandleEntries(rabbitMQConsumer, esUploader)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
