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
	"github.com/streadway/amqp"
)

const logEntriesExchangeName = "logentries"

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
		"fanout",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,                 // queue name
		"",                     // routing key
		logEntriesExchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
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

	go func() {
		for d := range msgs {
			if strings.Contains(string(d.Body), "START") {
				fmt.Println("Start of entries...")
				continue
			} else if strings.Contains(string(d.Body), "END") {
				fmt.Println("End of entries...")
				rabbitmq.SendStringMessageToElastic("DONE", channelToSendTo)
				continue
			}

			entry := deserializeMessage(d.Body)
			processing.Process(entry, channelToSendTo)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func deserializeMessage(message []byte) parsermodels.ParsedLine {
	var data parsermodels.ParsedLine
	if err := json.Unmarshal(message, &data); err != nil {
		fmt.Println("failed to unmarshal:", err)
	}

	return data
}
