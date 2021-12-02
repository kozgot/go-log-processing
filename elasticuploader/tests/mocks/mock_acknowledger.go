package mocks

// MockAcknowledger implements the Acknowledger interface of amqp.
type MockAcknowledger struct {
	expectedAckCount        int
	ackCount                int
	allMessagesAcknowledged chan bool
}

func NewMockAcknowleder(expectedAckCount int, allMessagesAcknowledged chan bool) *MockAcknowledger {
	mockAcknowledger := MockAcknowledger{
		expectedAckCount:        expectedAckCount,
		ackCount:                0,
		allMessagesAcknowledged: allMessagesAcknowledged,
	}
	return &mockAcknowledger
}

// Ack is the implementation of the Ack() function of the Acknowledger interface.
func (m *MockAcknowledger) Ack(tag uint64, multiple bool) error {
	m.ackCount++
	if m.expectedAckCount == m.ackCount {
		m.allMessagesAcknowledged <- true
	}
	return nil
}

// Nack is the implementation of the Nack() function of the Acknowledger interface.
func (m *MockAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	return nil
}

// Reject is the implementation of the Reject() function of the Acknowledger interface.
func (m *MockAcknowledger) Reject(tag uint64, requeue bool) error {
	return nil
}
