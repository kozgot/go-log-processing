package mocks

import (
	"github.com/kozgot/go-log-processing/postprocessor/tests/testmodels"
	"github.com/streadway/amqp"
)

// MockMessageConsumer mocks a message consumer, implements the MessageConsumer interface.
type MockMessageConsumer struct {
	TestParsedLogFile testmodels.TestParsedLogFile
}

// ConsumeMessages creates a channel from the parsed log file of the MockMessageConsumer.
func (m *MockMessageConsumer) ConsumeMessages() <-chan amqp.Delivery {
	deliveries := make(chan amqp.Delivery, 100)

	for i, entry := range m.TestParsedLogFile.Lines {
		data := entry.Serialize()

		mockDelivery := NewMockDelivery(data, uint64(i+1))
		deliveries <- mockDelivery
	}

	doneDelivery := NewMockDelivery([]byte("END"), uint64(len(m.TestParsedLogFile.Lines)+1))
	deliveries <- doneDelivery

	return deliveries
}

func (m *MockMessageConsumer) CloseConnectionAndChannel() {}
func (m *MockMessageConsumer) Connect()                   {}

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
