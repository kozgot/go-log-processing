package rabbitmq

import (
	"fmt"

	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"

	"github.com/streadway/amqp"
)

// AmqpProducer implements the ESUploader interface.
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
	esUploader := AmqpProducer{
		rabbitMqURL:  rabbitMqURL,
		exchangeName: exchangeName,
		routingKey:   routingKey}

	return &esUploader
}

// PublishEvent sends an SMC event to the uploader service.
func (uploader *AmqpProducer) PublishEvent(event models.SmcEvent, eventIndexName string) {
	dataToSend := models.DataUnit{IndexName: eventIndexName, Data: event.Serialize()}
	uploader.sendData(dataToSend.Serialize())
}

// PublishConsumption sends a consumption data item to the uploader service.
func (uploader *AmqpProducer) PublishConsumption(cons models.ConsumtionValue, consumptionIndexName string) {
	dataToSend := models.DataUnit{IndexName: consumptionIndexName, Data: cons.Serialize()}
	uploader.sendData(dataToSend.Serialize())
}

// Connect opens a channel and a connection.
func (uploader *AmqpProducer) Connect() {
	var err error
	uploader.connection, err = amqp.Dial(uploader.rabbitMqURL)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	fmt.Println("Created connection")

	// create the channel
	uploader.channel, err = uploader.connection.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	fmt.Println("Created channel")

	err = uploader.channel.ExchangeDeclare(
		uploader.exchangeName, // name
		"direct",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")
}

// CloseChannelAndConnection closes the channel and connection received in parameter.
func (uploader *AmqpProducer) CloseChannelAndConnection() {
	uploader.connection.Close()
	fmt.Println("Closed connection")
	uploader.channel.Close()
	fmt.Println("Closed channel")
}

// PublishCreateIndexMessage sends a string message to the message queue.
func (uploader *AmqpProducer) PublishCreateIndexMessage(indexName string) {
	bytes := []byte("CREATEINDEX|" + indexName)
	uploader.sendData(bytes)
}

// PublishDoneMessage sends a string message to the message queue.
func (uploader *AmqpProducer) PublishDoneMessage() {
	bytes := []byte("DONE")
	uploader.sendData(bytes)
}

func (uploader *AmqpProducer) sendData(data []byte) {
	body := data

	err := uploader.channel.Publish(
		uploader.exchangeName, // exchange
		uploader.routingKey,   // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})

	utils.FailOnError(err, "Failed to publish a message")
}
