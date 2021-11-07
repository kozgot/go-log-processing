package mocks

import "github.com/kozgot/go-log-processing/parser/pkg/models"

// MessageProducerMock mocks a rabbitmq message producer, implements the MessageProducer interface.
type MessageProducerMock struct {
	Entries []models.ParsedLogEntry
}

func (m *MessageProducerMock) PublishStringMessage(indexName string) {
	// NOOP
}

func (m *MessageProducerMock) PublishEntry(line models.ParsedLogEntry) {
	// Save sent entries to be able to validate them in the test.
	m.Entries = append(m.Entries, line)
}

func (m *MessageProducerMock) OpenChannelAndConnection() {
	// NOOP
}

func (m *MessageProducerMock) CloseChannelAndConnection() {
	// NOOP
}
