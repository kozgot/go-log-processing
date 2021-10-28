package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/kozgot/go-log-processing/parser/internal/utils"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/streadway/amqp"
)

const logEntriesExchangeName = "logentries_direct_durable"

// SendLineToPostProcessor sends the parsed log lines to the message queue.
func SendLineToPostProcessor(line models.ParsedLogEntry, channel *amqp.Channel) {
	sendDataToPostprocessor(serializeLine(line), channel)
}

// OpenChannelAndConnection opens a channel and a connection.
func OpenChannelAndConnection(rabbitMqURL string) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial(rabbitMqURL)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	fmt.Println("Created connection")

	// create the channel
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	fmt.Println("Created channel")

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

	return ch, conn
}

// CloseChannelAndConnection closes the channel and connection received in params.
func CloseChannelAndConnection(channel *amqp.Channel, connection *amqp.Connection) {
	connection.Close()
	fmt.Println("Closed connection")
	channel.Close()
	fmt.Println("Closed channel")
}

// SendStringMessageToPostProcessor sends a string message to the message queue.
func SendStringMessageToPostProcessor(indexName string, channel *amqp.Channel) {
	bytes := []byte(indexName)
	sendDataToPostprocessor(bytes, channel)
}

func serializeLine(line models.ParsedLogEntry) []byte {
	bytes, err := json.Marshal(line)
	if err != nil {
		fmt.Println("Can't serialize", line)
	}

	return bytes
}

func sendDataToPostprocessor(data []byte, channel *amqp.Channel) {
	body := data

	// TODO: extract routing key to a single place, eg.: env variables
	err := channel.Publish(
		logEntriesExchangeName, // exchange
		"process-entry",        // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	utils.FailOnError(err, "Failed to publish a message")
}
