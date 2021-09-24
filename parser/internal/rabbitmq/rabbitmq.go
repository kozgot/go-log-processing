package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/streadway/amqp"
)

// DataUnit contains the sent data unit.
type DataUnit struct {
	IndexName string
	Data      []byte
}

// SendLineToElastic sends the parsed log lines to the message queue.
func SendLineToElastic(line models.ParsedLogEntry, channel *amqp.Channel, indexName string) {
	// TODO: actual index name
	dataToSend := DataUnit{IndexName: indexName, Data: serializeLine(line)}
	// byteData := serializeLine(line)
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
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
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
func SendStringMessageToElastic(indexName string, channel *amqp.Channel) {
	bytes := []byte(indexName)
	sendData(bytes, channel)
}

func serializeLine(line models.ParsedLogEntry) []byte {
	bytes, err := json.Marshal(line)
	if err != nil {
		fmt.Println("Can't serialize", line)
	}

	return bytes
}

func serializeDataUnit(data DataUnit) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Can't serialize", data)
	}

	return bytes
}

func sendData(data []byte, channel *amqp.Channel) {
	body := data

	err := channel.Publish(
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
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
