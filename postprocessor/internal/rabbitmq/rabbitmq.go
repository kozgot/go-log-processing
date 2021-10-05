package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/streadway/amqp"
)

// todo: maybe set this as an environment variable
const processedDataExchangeName = "processeddata"

// SendEntryToElasticUploader sends the parsed log lines to the message queue.
func SendEntryToElasticUploader(entry models.SmcEntry, channel *amqp.Channel, indexName string) {
	dataToSend := models.DataUnit{IndexName: indexName, Data: serializeLine(entry)}
	sendData(serializeDataUnit((dataToSend)), channel)
}

// SendLinesToElastic sends the parsed log lines to the message queue.
func OpenChannelAndConnection(rabbitMqURL string) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial(rabbitMqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	fmt.Println("Created connection")

	// create the channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	fmt.Println("Created channel")

	err = ch.ExchangeDeclare(
		processedDataExchangeName, // name
		"fanout",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	return ch, conn
}

func CloseChannelAndConnection(channel *amqp.Channel, connection *amqp.Connection) {
	connection.Close()
	fmt.Println("Closed connection")
	channel.Close()
	fmt.Println("Closed channel")
}

// SendStringMessageToElastic sends a string message to the message queue.
func SendStringMessageToElastic(message string, channel *amqp.Channel) {
	bytes := []byte(message)
	sendData(bytes, channel)
}

func serializeLine(entry models.SmcEntry) []byte {
	bytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Can't serialize", entry)
	}

	return bytes
}

func serializeDataUnit(data models.DataUnit) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Can't serialize", data)
	}

	return bytes
}

func sendData(data []byte, channel *amqp.Channel) {
	body := data

	err := channel.Publish(
		processedDataExchangeName, // exchange
		"",                        // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
