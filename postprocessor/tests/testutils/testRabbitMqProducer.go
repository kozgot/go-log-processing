package testutils

import (
	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"
	"github.com/streadway/amqp"
)

// TestRabbitMqProducer implements the MessageProducer interface.
type TestRabbitMqProducer struct {
	rabbitMqURL  string
	connection   *amqp.Connection
	channel      *amqp.Channel
	routingKey   string
	exchangeName string
}

// NewTestProducer creates a new message producer that publishes messages to rabbitmq.
func NewTestProducer(
	rabbitMqURL string,
	exchangeName string,
	routingKey string) *TestRabbitMqProducer {
	testMessageProducer := TestRabbitMqProducer{
		rabbitMqURL:  rabbitMqURL,
		exchangeName: exchangeName,
		routingKey:   routingKey}

	return &testMessageProducer
}

// OpenChannelAndConnection opens a channel and a connection.
func (producer *TestRabbitMqProducer) Connect() {
	var err error
	producer.connection, err = amqp.Dial(producer.rabbitMqURL)
	utils.FailOnError(err, "  [TEST PRODUCER] Failed to connect to RabbitMQ")

	// create the channel
	producer.channel, err = producer.connection.Channel()
	utils.FailOnError(err, "  [TEST PRODUCER] Failed to open a channel")

	err = producer.channel.ExchangeDeclare(
		producer.exchangeName, // name
		"direct",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	utils.FailOnError(err, "  [TEST PRODUCER] Failed to declare an exchange")
}

// CloseChannelAndConnection closes the channel and connection received in params.
func (producer *TestRabbitMqProducer) CloseChannelAndConnection() {
	producer.connection.Close()
	producer.channel.Close()
}

// PublishStringMessage sends a string message to the message queue.
func (producer *TestRabbitMqProducer) PublishStringMessage(indexName string) {
	bytes := []byte(indexName)
	producer.sendDataToPostprocessor(bytes)
}

// PublishEntry sends the parsed log lines to the message queue.
func (producer *TestRabbitMqProducer) PublishEntry(line parsermodels.ParsedLogEntry) {
	producer.sendDataToPostprocessor(line.Serialize())
}

func (producer *TestRabbitMqProducer) sendDataToPostprocessor(data []byte) {
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
	utils.FailOnError(err, "  [TEST PRODUCER] Failed to publish a message")
}
