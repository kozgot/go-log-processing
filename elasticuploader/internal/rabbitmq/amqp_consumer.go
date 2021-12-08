package rabbitmq

import (
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/streadway/amqp"
)

// AmqpConsumer encapsulates data related to consuming rabbitmq messages.
type AmqpConsumer struct {
	rabbitMqURL  string
	channel      *amqp.Channel
	queue        amqp.Queue
	connection   *amqp.Connection
	exchangeName string
	queueName    string
	routingKey   string
}

// NewAmqpConsumer creates a new AmqpConsumer.
func NewAmqpConsumer(url string, exchangeName string, routingKey string, queueName string) *AmqpConsumer {
	consumer := AmqpConsumer{
		rabbitMqURL:  url,
		exchangeName: exchangeName,
		routingKey:   routingKey,
		queueName:    queueName,
	}
	return &consumer
}

func (c *AmqpConsumer) Connect() {
	var err error
	c.connection, err = amqp.Dial(c.rabbitMqURL)
	utils.FailOnError(err, "Could not connect to rabbitMQ.")
	c.channel, err = c.connection.Channel()
	utils.FailOnError(err, "Could not create a channel.")
}

func (c *AmqpConsumer) CloseChannelAndConnection() {
	c.connection.Close()
	c.channel.Close()
}

func (c *AmqpConsumer) Consume() (<-chan amqp.Delivery, error) {
	var err error
	var msgs <-chan amqp.Delivery
	err = c.channel.ExchangeDeclare(
		c.exchangeName, // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")

	c.queue, err = c.channel.QueueDeclare(
		c.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	err = c.channel.QueueBind(
		c.queue.Name,   // queue name
		c.routingKey,   // routing key
		c.exchangeName, // exchange
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to bind a queue")

	msgs, err = c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	return msgs, err
}
