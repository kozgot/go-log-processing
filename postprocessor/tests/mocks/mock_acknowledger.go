package mocks

// MockAcknowledger implements the Acknowledger interface of amqp.
type MockAcknowledger struct {
}

// Ack is the implementation of the Ack() function of the Acknowledger interface.
func (m *MockAcknowledger) Ack(tag uint64, multiple bool) error {
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
