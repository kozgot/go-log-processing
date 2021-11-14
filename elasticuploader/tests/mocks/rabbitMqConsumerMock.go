package mocks

import (
	"fmt"

	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/testmodels"
	"github.com/streadway/amqp"
)

// RabbitMQConsumerMock is a mock RabbitMQ consumer used in tests.
type RabbitMQConsumerMock struct {
	TestData      testmodels.TestProcessedData
	acknowledger  *MockAcknowledger
	testIndexName string
}

func NewRabbitMQConsumerMock(
	testData testmodels.TestProcessedData,
	done chan bool,
	testIndexName string,
) *RabbitMQConsumerMock {
	acknowledger := NewMockAcknowleder(len(testData.Consumptions)+len(testData.Events)+2, done)
	mock := RabbitMQConsumerMock{
		TestData:      testData,
		acknowledger:  acknowledger,
		testIndexName: testIndexName,
	}

	return &mock
}

// ConsumeMessages creates a channel from the parsed log file of the MockMessageConsumer.
func (m *RabbitMQConsumerMock) Consume() (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery, 100)

	createIndexDelivery := NewMockDelivery([]byte("RECREATEINDEX|"+m.testIndexName), uint64(0), m.acknowledger)
	deliveries <- createIndexDelivery
	fmt.Println("Created index " + m.testIndexName)

	for i, cons := range m.TestData.Consumptions {
		messageBytes := cons.Serialize()
		mockDelivery := NewMockDelivery(messageBytes, uint64(i+1), m.acknowledger)
		deliveries <- mockDelivery
	}

	for i, event := range m.TestData.Events {
		messageBytes := models.ReceivedDataUnit{IndexName: m.testIndexName, Data: event.Serialize()}
		mockDelivery := NewMockDelivery(messageBytes.ToJSON(), uint64(i), m.acknowledger)
		deliveries <- mockDelivery
	}

	doneDelivery := NewMockDelivery(
		[]byte("DONE"),
		uint64(len(m.TestData.Consumptions)+len(m.TestData.Events)+1),
		m.acknowledger)
	deliveries <- doneDelivery

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
