package rabbitmq

import (
	"log"

	"github.com/kozgot/go-log-processing/postprocessor/internal/utils"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/streadway/amqp"
)

// AmqpProducer implements the MessageProducer interface.
type AmqpProducer struct {
	rabbitMqURL  string
	connection   *amqp.Connection
	channel      *amqp.Channel
	routingKey   string
	exchangeName string
}

// NewAmqpProducer creates a new message producer that publishes messages to rabbitmq.
func NewAmqpProducer(
	rabbitMqURL string,
	exchangeName string,
	routingKey string) *AmqpProducer {
	producer := AmqpProducer{
		rabbitMqURL:  rabbitMqURL,
		exchangeName: exchangeName,
		routingKey:   routingKey}

	return &producer
}

// PublishEvent sends an SMC event to the uploader service.
func (producer *AmqpProducer) PublishEvent(event models.SmcEvent, eventIndexName string) {
	dataToSend := models.DataUnit{IndexName: eventIndexName, Data: event.Serialize()}
	producer.publishData(dataToSend.Serialize())
}

// PublishConsumption sends a consumption data item to the uploader service.
func (producer *AmqpProducer) PublishConsumption(cons models.ConsumtionValue, consumptionIndexName string) {
	dataToSend := models.DataUnit{IndexName: consumptionIndexName, Data: cons.Serialize()}
	producer.publishData(dataToSend.Serialize())
}

// Connect opens a channel and a connection.
func (producer *AmqpProducer) Connect() {
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
func (producer *AmqpProducer) CloseChannelAndConnection() {
	producer.connection.Close()
	log.Println(" [AMQP PRODUCER] Closed connection")
	producer.channel.Close()
	log.Println(" [AMQP PRODUCER] Closed channel")
}

// PublishRecreateIndexMessage sends a string message to the message queue.
func (producer *AmqpProducer) PublishRecreateIndexMessage(indexName string) {
	bytes := []byte("RECREATEINDEX|" + indexName)
	producer.publishData(bytes)
}

// PublishDoneMessage sends a string message to the message queue.
func (producer *AmqpProducer) PublishDoneMessage() {
	bytes := []byte("DONE")
	producer.publishData(bytes)
}

func (producer *AmqpProducer) publishData(data []byte) {
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
