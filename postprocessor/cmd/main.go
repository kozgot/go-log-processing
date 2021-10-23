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
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/streadway/amqp"
)

const logEntriesExchangeName = "logentries_direct_durable"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	log.Println("PostProcessor service starting...")
	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("RABBIT_URL:", rabbitMqURL)

	conn, err := amqp.Dial(rabbitMqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
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
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"processing_queue_durable", // name
		true,                       // durable
		false,                      // delete when unused
		true,                       // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// TODO: extract routing key to a single place, eg.: env variables
	err = ch.QueueBind(
		q.Name,                 // queue name
		"process-entry",        // routing key
		logEntriesExchangeName, // exchange
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

	forever := make(chan bool)

	channelToSendTo, connectionToSendTo := rabbitmq.OpenChannelAndConnection(rabbitMqURL)
	defer rabbitmq.CloseChannelAndConnection(channelToSendTo, connectionToSendTo)

	rabbitmq.SendStringMessageToElastic("CREATEINDEX|"+"smc", channelToSendTo)
	rabbitmq.SendStringMessageToElastic("CREATEINDEX|"+"routing", channelToSendTo)
	rabbitmq.SendStringMessageToElastic("CREATEINDEX|"+"status", channelToSendTo)

	eventsBySmcUID := make(map[string][]models.SmcEvent)
	smcDataBySmcUID := make(map[string]models.SmcData)
	smcUIDsByURL := make(map[string]string)
	consumptionValues := []models.ConsumtionValue{}
	indexValues := []models.IndexValue{}

	go func() {
		for d := range msgs {
			if strings.Contains(string(d.Body), "START") {
				fmt.Println("Start of entries...")

				// Acknowledge the message after it has been processed.
				err := d.Ack(false)
				failOnError(err, "Could not acknowledge START message")

				continue
			} else if strings.Contains(string(d.Body), "END") {
				fmt.Println("End of entries...")
				rabbitmq.SendStringMessageToElastic("DONE", channelToSendTo)

				// Acknowledge the message after it has been processed.
				err := d.Ack(false)
				failOnError(err, "Could not acknowledge END message")

				continue
			}

			entry := deserializeMessage(d.Body)
			consumption, index := processing.Process(entry, channelToSendTo, eventsBySmcUID, smcDataBySmcUID, smcUIDsByURL)
			if index != nil {
				indexValues = append(indexValues, *index)
			}
			if consumption != nil {
				consumptionValues = append(consumptionValues, *consumption)
			}

			// Acknowledge the message after it has been processed.
			err := d.Ack(false)
			failOnError(err, "Could not acknowledge message with timestamp: "+entry.Timestamp.Format("2 Jan 2006 15:04:05"))
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func deserializeMessage(message []byte) parsermodels.ParsedLogEntry {
	var data parsermodels.ParsedLogEntry
	if err := json.Unmarshal(message, &data); err != nil {
		fmt.Println("failed to unmarshal:", err)
	}

	return data
}
