package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kozgot/go-log-processing/postprocessor/internal/processing"
	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
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

	// Init message consumer.
	rabbitMQConsumer := rabbitmq.NewAmqpConsumer(
		rabbitMqURL,
		processEntryRoutingKey,
		processEntriesExchangeName,
		processingQueueName)

	// Open consumer channel and connection.
	rabbitMQConsumer.Connect()
	defer rabbitMQConsumer.CloseConnectionAndChannel()

	// Init message producer.
	rabbitMqProducer := rabbitmq.NewAmqpProducer(
		rabbitMqURL,
		saveDataExchangeName,
		saveDataRoutingKey)

	// Open producer channel and connection.
	rabbitMqProducer.Connect()
	defer rabbitMqProducer.CloseChannelAndConnection()

	forever := make(chan bool)

	processor := processing.NewEntryProcessor(rabbitMqProducer, rabbitMQConsumer, eventsIndexName, consumptionIndexName)
	processor.HandleEntries()

	log.Printf(" [POSTPROCESSOR] Waiting for messages. To exit press CTRL+C...")
	<-forever
}
