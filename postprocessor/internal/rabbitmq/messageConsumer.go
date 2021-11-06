package rabbitmq

import "github.com/streadway/amqp"

// MessageConsumer encapsulates messages needed to consume messages.
type MessageConsumer interface {
	ConsumeMessages() (<-chan amqp.Delivery, error)
}
