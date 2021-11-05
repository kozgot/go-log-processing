package rabbitmq

import (
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"
	"github.com/streadway/amqp"
)

// AmqpConsumer encapsulates data related to consuming messages from rabbitmq.
type AmqpConsumer struct {
	hostURL      string
	channel      *amqp.Channel
	queue        amqp.Queue
	connection   *amqp.Connection
	exchangeName string
	queueName    string
	routingKey   string
}

// Creates a new AmqpConsumer.
func NewAmqpConsumer(hostURL string, routingKey string, exchangeName string, queueName string) *AmqpConsumer {
	rabbitMQConsumer := AmqpConsumer{
		hostURL:      hostURL,
		routingKey:   routingKey,
		exchangeName: exchangeName,
		queueName:    queueName}

	return &rabbitMQConsumer
}

// Connect initializes a connection.
func (c *AmqpConsumer) Connect() error {
	var err error
	c.connection, err = amqp.Dial(c.hostURL)

	return err
}

// CloseConnection closes the connection.
func (c *AmqpConsumer) CloseConnection() {
	c.connection.Close()
}

// Channel initializes a channel.
func (c *AmqpConsumer) Channel() error {
	var err error
	c.channel, err = c.connection.Channel()

	return err
}

// CloseChannel closes the channel.
func (c *AmqpConsumer) CloseChannel() {
	c.channel.Close()
}

// Consume consumes messages from rabbitmq, returns the deliveries.
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
		true,        // exclusive
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
