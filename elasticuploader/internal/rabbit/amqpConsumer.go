package rabbit

import (
	"log"

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

func (c *AmqpConsumer) Connect() error {
	var err error
	c.connection, err = amqp.Dial(c.rabbitMqURL)

	return err
}

func (c *AmqpConsumer) CloseConnection() {
	c.connection.Close()
}

func (c *AmqpConsumer) Channel() error {
	var err error
	c.channel, err = c.connection.Channel()

	return err
}

func (c *AmqpConsumer) CloseChannel() {
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
	failOnError(err, "Failed to declare an exchange")

	c.queue, err = c.channel.QueueDeclare(
		c.queueName, // name
		true,        // durable
		false,       // delete when unused
		true,        // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = c.channel.QueueBind(
		c.queue.Name,   // queue name
		c.routingKey,   // routing key
		c.exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err = c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	return msgs, err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
