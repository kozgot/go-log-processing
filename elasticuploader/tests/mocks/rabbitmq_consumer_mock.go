package mocks

import (
	"log"
	"time"

	"github.com/kozgot/go-log-processing/elasticuploader/tests/testmodels"
	postprocmodels "github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/streadway/amqp"
)

// RabbitMQConsumerMock is a mock RabbitMQ consumer used in tests.
type RabbitMQConsumerMock struct {
	TestData     testmodels.TestProcessedData
	acknowledger *MockAcknowledger
	// this is needed to test the timed index-recreation functionality in the uploader service
	// ignored, if set to zero
	deliveryDelaySeconds int
}

// NewRabbitMQConsumerMock creates a new mock consumer
// for providing processed data messages for the ulpoader service.
// The testData will be used to create the processed data messages from.
// The done channel is used to signal after all messages have been acknowledged.
// The deliveryDelaySeconds is used to add an artificial delay after the first message,
// to test the timed index recreation behaviour, it is ignored if set to zero.
func NewRabbitMQConsumerMock(
	testData testmodels.TestProcessedData,
	done chan bool,
	deliveryDelaySeconds int,
	expectedDocCount int,
) *RabbitMQConsumerMock {
	acknowledger := NewMockAcknowleder(expectedDocCount, done)
	mock := RabbitMQConsumerMock{
		TestData:             testData,
		acknowledger:         acknowledger,
		deliveryDelaySeconds: deliveryDelaySeconds,
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

	if m.deliveryDelaySeconds > 0 {
		go func() {
			log.Printf("Waiting %d seconds ...", m.deliveryDelaySeconds)
			time.Sleep(time.Duration(m.deliveryDelaySeconds) * time.Second)
			log.Printf("%d seconds passed, continue consuming messages...", m.deliveryDelaySeconds)
			for i, event := range m.TestData.Events {
				message := postprocmodels.DataUnit{DataType: postprocmodels.Event, Data: event.Serialize()}
				mockDelivery := NewMockDelivery(message.Serialize(), uint64(i), m.acknowledger)
				deliveries <- mockDelivery
			}
		}()
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
