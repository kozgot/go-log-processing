package mocks

import (
	"github.com/kozgot/go-log-processing/elasticuploader/tests/testmodels"
	postprocmodels "github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/streadway/amqp"
)

// RabbitMQConsumerMock is a mock RabbitMQ consumer used in tests.
type RabbitMQConsumerMock struct {
	TestData     testmodels.TestProcessedData
	acknowledger *MockAcknowledger
}

// NewRabbitMQConsumerMock creates a new mock consumer
// for providing processed data messages for the ulpoader service.
// The testData will be used to create the processed data messages from.
// The done channel is used to signal after all messages have been acknowledged.
func NewRabbitMQConsumerMock(
	testData testmodels.TestProcessedData,
	done chan bool,
) *RabbitMQConsumerMock {
	acknowledger := NewMockAcknowleder(len(testData.Consumptions)+len(testData.Events), done)
	mock := RabbitMQConsumerMock{
		TestData:     testData,
		acknowledger: acknowledger,
	}

	return &mock
}

// ConsumeMessages creates a channel from the parsed log file of the MockMessageConsumer.
func (m *RabbitMQConsumerMock) Consume() (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery, 100)

	for i, cons := range m.TestData.Consumptions {
		message := postprocmodels.DataUnit{DataType: postprocmodels.Consumption, Data: cons.Serialize()}
		mockDelivery := NewMockDelivery(message.Serialize(), uint64(i+1), m.acknowledger)
		deliveries <- mockDelivery
	}

	for i, event := range m.TestData.Events {
		message := postprocmodels.DataUnit{DataType: postprocmodels.Event, Data: event.Serialize()}
		mockDelivery := NewMockDelivery(message.Serialize(), uint64(i), m.acknowledger)
		deliveries <- mockDelivery
	}

	return deliveries, nil
}

// NewMockDelivery creates an amqp.Delivery from a parsed log entry,
// filling only the properties to used by the processor code.
func NewMockDelivery(data []byte, tag uint64, acknowledger *MockAcknowledger) amqp.Delivery {
	delivery := amqp.Delivery{
		Acknowledger: acknowledger,
		DeliveryTag:  tag,
		Body:         data,
	}
	return delivery
}
