package mocks

import (
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/tests/testmodels"
)

type MockMessageProducer struct {
	Data testmodels.TestProcessedData
	Done chan string
}

// PublishEvent is the implementation of the PublishEvent(event models.SmcEvent)
// function of the MessageProducer interface.
func (m *MockMessageProducer) PublishEvent(event models.SmcEvent) {
	m.Data.Events = append(m.Data.Events, event)
}

// PublishConsumption is the implementation
// of the PublishConsumption(cons models.ConsumtionValue)
// function of the MessageProducer interface.
func (m *MockMessageProducer) PublishConsumption(cons models.ConsumtionValue) {
	m.Data.Consumptions = append(m.Data.Consumptions, cons)
}

// Connect is the implementation of the Connect() function of the MessageProducer interface.
func (m *MockMessageProducer) Connect() {
	// NOOP

}

// CloseChannelAndConnection is the implementation of
// the CloseChannelAndConnection() function of the MessageProducer interface.
func (m *MockMessageProducer) CloseChannelAndConnection() {
	// NOOP
}

// PublishDoneMessage is the implementation of the PublishDoneMessage() function of the MessageProducer interface.
func (m *MockMessageProducer) PublishDoneMessage() {
	m.Done <- "done"
}
