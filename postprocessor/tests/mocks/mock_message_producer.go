package mocks

import (
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/tests/testmodels"
)

type MockMessageProducer struct {
	Data                     testmodels.TestProcessedData
	done                     chan string
	expectedEventCount       int
	expectedConsumptionCount int
	publishedDataCount       int
}

func NewMockMessageProducer(
	data testmodels.TestProcessedData,
	done chan string,
	expectedEventCount int,
	expectedConsumptionCount int,
) *MockMessageProducer {
	return &MockMessageProducer{
		Data:                     data,
		done:                     done,
		expectedEventCount:       expectedEventCount,
		expectedConsumptionCount: expectedConsumptionCount,
		publishedDataCount:       0,
	}
}

// PublishEvent is the implementation of the PublishEvent(event models.SmcEvent)
// function of the MessageProducer interface.
func (m *MockMessageProducer) PublishEvent(event models.SmcEvent) {
	m.Data.Events = append(m.Data.Events, event)
	m.publishedDataCount++
	if m.publishedDataCount == m.expectedConsumptionCount+m.expectedEventCount {
		m.done <- "DONE"
	}
}

// PublishConsumption is the implementation
// of the PublishConsumption(cons models.ConsumtionValue)
// function of the MessageProducer interface.
func (m *MockMessageProducer) PublishConsumption(cons models.ConsumtionValue) {
	m.Data.Consumptions = append(m.Data.Consumptions, cons)
	m.publishedDataCount++
	if m.publishedDataCount == m.expectedConsumptionCount+m.expectedEventCount {
		m.done <- "DONE"
	}
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
