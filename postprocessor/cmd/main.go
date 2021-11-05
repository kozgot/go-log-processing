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

	processedDataExchangeName := os.Getenv("PROCESSED_DATA_EXCHANGE")
	fmt.Println("PROCESSED_DATA_EXCHANGE:", processedDataExchangeName)
	if len(processedDataExchangeName) == 0 {
		log.Fatal("The PROCESSED_DATA_EXCHANGE environment variable is not set")
	}

	saveDataRoutingKey := os.Getenv("SAVE_DATA_ROUTING_KEY")
	fmt.Println("SAVE_DATA_ROUTING_KEY:", saveDataRoutingKey)
	if len(saveDataRoutingKey) == 0 {
		log.Fatal("The SAVE_DATA_ROUTING_KEY environment variable is not set")
	}

	logEntriesExchangeName := os.Getenv("LOG_ENTRIES_EXCHANGE")
	fmt.Println("LOG_ENTRIES_EXCHANGE:", logEntriesExchangeName)
	if len(logEntriesExchangeName) == 0 {
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

	conn, err := amqp.Dial(rabbitMqURL)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		logEntriesExchangeName, // name
		"direct",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		processingQueueName, // name
		true,                // durable
		false,               // delete when unused
		true,                // exclusive
		false,               // no-wait
		nil,                 // arguments
	)

	utils.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,                 // queue name
		processEntryRoutingKey, // routing key
		logEntriesExchangeName, // exchange
		false,
		nil,
	)

	utils.FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	utils.FailOnError(err, "Failed to register a consumer")

	esUploader := rabbitmq.EsUploadSender{
		RabbitMqURL:  rabbitMqURL,
		ExchangeName: processedDataExchangeName,
		RoutingKey:   saveDataRoutingKey}
	esUploader.OpenChannelAndConnection(rabbitMqURL)

	defer esUploader.CloseChannelAndConnection()

	forever := make(chan bool)

	// Create indices in ES.
	esUploader.CreateIndex(smcIndexName)
	esUploader.CreateIndex(consumptionIndexName)

	processor := processing.InitEntryProcessor(&esUploader)

	go func() {
		for d := range msgs {
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
