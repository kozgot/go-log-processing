package rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

type AmqpConsumer struct {
	HostDsn    string
	channel    *amqp.Channel
	queue      amqp.Queue
	connection *amqp.Connection
}

func (c *AmqpConsumer) Connect() error {
	var err error
	c.connection, err = amqp.Dial(c.HostDsn)

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

func (c *AmqpConsumer) Consume(exchangeName string, queueName string, routingKey string) (<-chan amqp.Delivery, error) {
	var err error
	var msgs <-chan amqp.Delivery
	err = c.channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	c.queue, err = c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = c.channel.QueueBind(
		c.queue.Name, // queue name
		routingKey,   // routing key
		exchangeName, // exchange
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
