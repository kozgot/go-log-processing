package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/elastic"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/rabbit"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/uploader"
)

func main() {
	log.Println("Elastic Uploader starting...")

	// Load environment variables.
	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("RABBIT_URL:", rabbitMqURL)
	if len(rabbitMqURL) == 0 {
		log.Fatal("The RABBIT_URL environment variable is not set")
	}

	elasticSearchURL := os.Getenv("ELASTICSEARCH_URL")
	fmt.Println("ELASTICSEARCH_URL:", elasticSearchURL)
	if len(elasticSearchURL) == 0 {
		log.Fatal("The ELASTICSEARCH_URL environment variable is not set")
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

	eventIndexName := os.Getenv("EVENT_INDEX_NAME")
	fmt.Println("EVENT_INDEX_NAME:", eventIndexName)
	if len(eventIndexName) == 0 {
		log.Fatal("The EVENT_INDEX_NAME environment variable is not set")
	}

	consumptionIndexName := os.Getenv("CONSUMPTION_INDEX_NAME")
	fmt.Println("CONSUMPTION_INDEX_NAME:", consumptionIndexName)
	if len(consumptionIndexName) == 0 {
		log.Fatal("The CONSUMPTION_INDEX_NAME environment variable is not set")
	}

	// Setup ES client.
	esClient := elastic.NewEsClientWrapper(elasticSearchURL)

	// Setup rabbitmq consumer.
	rabbitMQConsumer := rabbit.NewAmqpConsumer(
		rabbitMqURL,
		processedDataExchangeName,
		saveDataRoutingKey,
		saveDataQueueName)

	// Setup rabbitmq connection and channel.
	rabbitMQConsumer.Connect()
	defer rabbitMQConsumer.CloseChannelAndConnection()

	forever := make(chan bool)

	// Start handling messages.
	uploaderService := uploader.NewUploaderService(rabbitMQConsumer, esClient, eventIndexName, consumptionIndexName)
	uploaderService.HandleMessages()

	log.Printf(" [ESUPLOADER] Waiting for messages. To exit press CTRL+C")

	<-forever
}
