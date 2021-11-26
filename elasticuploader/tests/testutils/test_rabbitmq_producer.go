package testutils

import (
	"log"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/testmodels"
	"github.com/streadway/amqp"
)

// TestRabbitMQProducer is a rabbitMQ producer used in tests
// to publish test input messages into a queue.
type TestRabbitMqProducer struct {
	rabbitMqURL  string
	connection   *amqp.Connection
	channel      *amqp.Channel
	routingKey   string
	exchangeName string
}

// NewTestRabbitMqProducer creates a new test producer that publishes messages to rabbitmq.
func NewTestRabbitMqProducer(
	rabbitMqURL string,
	exchangeName string,
	routingKey string) *TestRabbitMqProducer {
	producer := TestRabbitMqProducer{
		rabbitMqURL:  rabbitMqURL,
		exchangeName: exchangeName,
		routingKey:   routingKey}

	return &producer
}

// Connect opens a channel and a connection.
func (producer *TestRabbitMqProducer) Connect() {
	var err error
	producer.connection, err = amqp.Dial(producer.rabbitMqURL)
	utils.FailOnError(err, " [AMQP PRODUCER] Failed to connect to RabbitMQ")
	log.Println(" [AMQP PRODUCER] Created connection")

	// create the channel
	producer.channel, err = producer.connection.Channel()
	utils.FailOnError(err, " [AMQP PRODUCER] Failed to open a channel")
	log.Println(" [AMQP PRODUCER] Created channel")

	err = producer.channel.ExchangeDeclare(
		producer.exchangeName, // name
		"direct",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	utils.FailOnError(err, " [AMQP PRODUCER] Failed to declare an exchange")
}

// CloseChannelAndConnection closes the channel and connection received in parameter.
func (producer *TestRabbitMqProducer) CloseChannelAndConnection() {
	producer.connection.Close()
	log.Println(" [AMQP PRODUCER] Closed connection")
	producer.channel.Close()
	log.Println(" [AMQP PRODUCER] Closed channel")
}

func (producer *TestRabbitMqProducer) PublishTestInput(testData testmodels.TestProcessedData, testIndexName string) {
	for _, event := range testData.Events {
		dataToSend := models.ReceivedDataUnit{IndexName: testIndexName, Data: event.Serialize()}
		producer.publishData(dataToSend.ToJSON())
	}

	for _, event := range testData.Consumptions {
		dataToSend := models.ReceivedDataUnit{IndexName: testIndexName, Data: event.Serialize()}
		producer.publishData(dataToSend.ToJSON())
	}
}

// PublishDoneMessage sends a string message to the message queue.
func (producer *TestRabbitMqProducer) PublishDoneMessage() {
	bytes := []byte("DONE")
	producer.publishData(bytes)
}

func (producer *TestRabbitMqProducer) publishData(data []byte) {
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

	utils.FailOnError(err, " [AMQP PRODUCER] Failed to publish a message")
}
