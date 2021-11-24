package rabbit

import "github.com/streadway/amqp"

// MessageConsumer encapsulates messages needed to consume rabbitmq messages.
type MessageConsumer interface {
	Consume() (<-chan amqp.Delivery, error)
}
