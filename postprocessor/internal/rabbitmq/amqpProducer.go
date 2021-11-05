package rabbitmq

import (
	"fmt"

	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"

	"github.com/streadway/amqp"
)

type EsUploadSender struct {
	RabbitMqURL  string
	connection   *amqp.Connection
	channel      *amqp.Channel
	RoutingKey   string
	ExchangeName string
}

// SendEventToElasticUploader sends an SMC event to the uploader service.
func (uploader *EsUploadSender) SendEventToElasticUploader(event models.SmcEvent, indexName string) {
	dataToSend := models.DataUnit{IndexName: indexName, Data: event.Serialize()}
	uploader.sendData(dataToSend.Serialize())
}

// SendConsumptionToElasticUploader sends a consumption data item to the uploader service.
func (uploader *EsUploadSender) SendConsumptionToElasticUploader(cons models.ConsumtionValue, indexName string) {
	dataToSend := models.DataUnit{IndexName: indexName, Data: cons.Serialize()}
	uploader.sendData(dataToSend.Serialize())
}

// OpenChannelAndConnection opens a channel and a connection and returns pointers to them.
func (uploader *EsUploadSender) OpenChannelAndConnection(rabbitMqURL string) {
	var err error
	uploader.connection, err = amqp.Dial(rabbitMqURL)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	fmt.Println("Created connection")

	// create the channel
	uploader.channel, err = uploader.connection.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	fmt.Println("Created channel")

	err = uploader.channel.ExchangeDeclare(
		uploader.ExchangeName, // name
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
func (uploader *EsUploadSender) CloseChannelAndConnection() {
	uploader.connection.Close()
	fmt.Println("Closed connection")
	uploader.channel.Close()
	fmt.Println("Closed channel")
}

// CreateIndex sends a string message to the message queue.
func (uploader *EsUploadSender) CreateIndex(indexName string) {
	bytes := []byte("CREATEINDEX|" + indexName)
	uploader.sendData(bytes)
}

func (uploader *EsUploadSender) sendData(data []byte) {
	body := data

	err := uploader.channel.Publish(
		uploader.ExchangeName, // exchange
		uploader.RoutingKey,   // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	utils.FailOnError(err, "Failed to publish a message")
}
