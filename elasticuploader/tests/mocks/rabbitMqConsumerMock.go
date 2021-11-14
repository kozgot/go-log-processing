package mocks

import (
	"github.com/streadway/amqp"
)

// RabbitMQConsumerMock is a mock RabbitMQ consumer used in tests.
type RabbitMQConsumerMock struct {
	Messages [][]byte
}

// ConsumeMessages creates a channel from the parsed log file of the MockMessageConsumer.
func (m *RabbitMQConsumerMock) ConsumeMessages() <-chan amqp.Delivery {
	deliveries := make(chan amqp.Delivery, 100)

	for i, messageBytes := range m.Messages {
		mockDelivery := NewMockDelivery(messageBytes, uint64(i))
		deliveries <- mockDelivery
	}

	doneDelivery := NewMockDelivery([]byte("DONE"), uint64(len(m.Messages)+1))
	deliveries <- doneDelivery

	return deliveries
}

// NewMockDelivery creates an amqp.Delivery from a parsed log entry,
// filling only the properties to used by the processor code.
func NewMockDelivery(data []byte, tag uint64) amqp.Delivery {
	acknowledger := MockAcknowledger{}
	delivery := amqp.Delivery{
		Acknowledger: &acknowledger,
		DeliveryTag:  tag,
		Body:         data,
	}
	return delivery
}
