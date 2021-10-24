package rabbitmq

import (
	"fmt"

	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"

	"github.com/streadway/amqp"
)

// todo: maybe set this as an environment variable
const processedDataExchangeName = "processeddata_direct_durable"

// SendEventToElasticUploader sends an SMC event to the uploader service.
func SendEventToElasticUploader(event models.SmcEvent, channel *amqp.Channel, indexName string) {
	dataToSend := models.DataUnit{IndexName: indexName, Data: event.Serialize()}
	sendData(dataToSend.Serialize(), channel)
}

// SendConsumptionToElasticUploader sends a consumption data item to the uploader service.
func SendConsumptionToElasticUploader(cons models.ConsumtionValue, channel *amqp.Channel, indexName string) {
	dataToSend := models.DataUnit{IndexName: indexName, Data: cons.Serialize()}
	sendData(dataToSend.Serialize(), channel)
}

// OpenChannelAndConnection opens a channel and a connection and returns pointers to them.
func OpenChannelAndConnection(rabbitMqURL string) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial(rabbitMqURL)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	fmt.Println("Created connection")

	// create the channel
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	fmt.Println("Created channel")

	err = ch.ExchangeDeclare(
		processedDataExchangeName, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")

	return ch, conn
}

// CloseChannelAndConnection closes the channel and connection received in parameter.
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

func sendData(data []byte, channel *amqp.Channel) {
	body := data

	// TODO: extract routing key to a single place, eg.: env variables
	err := channel.Publish(
		processedDataExchangeName, // exchange
		"save-data",               // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	utils.FailOnError(err, "Failed to publish a message")
}
