package rabbitmq

import (
	"fmt"

	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"

	"github.com/streadway/amqp"
)

// EsUploader implements the ESUploader interface.
type EsUploader struct {
	rabbitMqURL          string
	connection           *amqp.Connection
	channel              *amqp.Channel
	routingKey           string
	exchangeName         string
	eventIndexName       string
	consumptionIndexName string
}

// NewEsUploadSender creates a new EsUploadSender.
func NewEsUploadSender(
	rabbitMqURL string,
	exchangeName string,
	routingKey string,
	eventIndexName string,
	consumptionIndexName string) *EsUploader {
	esUploader := EsUploader{
		rabbitMqURL:          rabbitMqURL,
		exchangeName:         exchangeName,
		routingKey:           routingKey,
		eventIndexName:       eventIndexName,
		consumptionIndexName: consumptionIndexName}

	return &esUploader
}

// SendEventToElasticUploader sends an SMC event to the uploader service.
func (uploader *EsUploader) SendEventToElasticUploader(event models.SmcEvent) {
	dataToSend := models.DataUnit{IndexName: uploader.eventIndexName, Data: event.Serialize()}
	uploader.sendData(dataToSend.Serialize())
}

// SendConsumptionToElasticUploader sends a consumption data item to the uploader service.
func (uploader *EsUploader) SendConsumptionToElasticUploader(cons models.ConsumtionValue) {
	dataToSend := models.DataUnit{IndexName: uploader.consumptionIndexName, Data: cons.Serialize()}
	uploader.sendData(dataToSend.Serialize())
}

// Connect opens a channel and a connection.
func (uploader *EsUploader) Connect(rabbitMqURL string) {
	var err error
	uploader.connection, err = amqp.Dial(rabbitMqURL)
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
func (uploader *EsUploader) CloseChannelAndConnection() {
	uploader.connection.Close()
	fmt.Println("Closed connection")
	uploader.channel.Close()
	fmt.Println("Closed channel")
}

// CreateIndex sends a string message to the message queue.
func (uploader *EsUploader) CreateIndex(indexName string) {
	bytes := []byte("CREATEINDEX|" + indexName)
	uploader.sendData(bytes)
}

func (uploader *EsUploader) sendData(data []byte) {
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
