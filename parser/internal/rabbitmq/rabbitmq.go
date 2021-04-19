package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/streadway/amqp"
)

// SendLinesToElastic sends the parsed log lines to the message queue.
func SendLinesToElastic(rabbitMqURL string, line models.ParsedLine) {
	byteData := serializeLine(line)
	sendData(rabbitMqURL, byteData)
}

// SendStringMessageToElastic sends a string message to the message queue.
func SendStringMessageToElastic(rabbitMqURL string, indexName string) {
	bytes := []byte(indexName)
	sendData(rabbitMqURL, bytes)
}

func serializeLine(line models.ParsedLine) []byte {
	bytes, err := json.Marshal(line)
	if err != nil {
		fmt.Println("Can't serialize", line)
	}

	return bytes
}

func sendData(rabbitMqURL string, data []byte) {
	conn, err := amqp.Dial(rabbitMqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// create the channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body := data

	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")

	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
