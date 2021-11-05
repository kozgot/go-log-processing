package rabbitmq

import (
	"fmt"

	"github.com/kozgot/go-log-processing/parser/internal/utils"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/streadway/amqp"
)

type AmqpProducer struct {
	connection   *amqp.Connection
	channel      *amqp.Channel
	routingKey   string
	exchangeName string
	rabbitMqURL  string
}

func InitProducer(routingKey string, exchangeName string, rabbitMqURL string) *AmqpProducer {
	result := AmqpProducer{routingKey: routingKey, exchangeName: exchangeName, rabbitMqURL: rabbitMqURL}
	return &result
}

// OpenChannelAndConnection opens a channel and a connection.
func (producer *AmqpProducer) OpenChannelAndConnection() {
	var err error
	producer.connection, err = amqp.Dial(producer.rabbitMqURL)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	fmt.Println("Created connection")

	// create the channel
	producer.channel, err = producer.connection.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	fmt.Println("Created channel")

	err = producer.channel.ExchangeDeclare(
		producer.exchangeName, // name
		"direct",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")
}

// CloseChannelAndConnection closes the channel and connection received in params.
func (producer *AmqpProducer) CloseChannelAndConnection() {
	producer.connection.Close()
	fmt.Println("Closed connection")
	producer.channel.Close()
	fmt.Println("Closed channel")
}

// SendStringMessageToPostProcessor sends a string message to the message queue.
func (producer *AmqpProducer) SendStringMessageToPostProcessor(indexName string) {
	bytes := []byte(indexName)
	producer.sendDataToPostprocessor(bytes)
}

// SendDataToPostProcessor sends the parsed log lines to the message queue.
func (producer *AmqpProducer) SendDataToPostProcessor(line models.ParsedLogEntry) {
	producer.sendDataToPostprocessor(line.Serialize())
}

func (producer *AmqpProducer) sendDataToPostprocessor(data []byte) {
	body := data

	err := producer.channel.Publish(
		producer.exchangeName, // exchange
		producer.routingKey,   // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	utils.FailOnError(err, "Failed to publish a message")
}
