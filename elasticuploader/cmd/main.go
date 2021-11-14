package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/elastic"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/rabbit"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/service"
)

func main() {
	log.Println("Elastic Uploader starting...")

	// Load environment variables.
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

	// Setup ES client.
	esClient := elastic.NewEsClientWrapper()

	// Setup rabbitmq consumer.
	rabbitMQConsumer := rabbit.NewAmqpConsumer(
		rabbitMqURL,
		processedDataExchangeName,
		saveDataRoutingKey,
		saveDataQueueName)

	// Setup rabbitmq connection.
	err := rabbitMQConsumer.Connect()
	utils.FailOnError(err, "Could not connect ro RabbitMQ")
	defer rabbitMQConsumer.CloseConnection()

	// Setup rabbitmq channel.
	err = rabbitMQConsumer.Channel()
	utils.FailOnError(err, "Could not open channel")
	defer rabbitMQConsumer.CloseChannel()

	forever := make(chan bool)

	// Start handling messages.
	uploaderService := service.NewUploaderService(rabbitMQConsumer, esClient)
	uploaderService.HandleMessages()

	log.Printf(" [ESUPLOADER] Waiting for messages. To exit press CTRL+C")

	<-forever
}
