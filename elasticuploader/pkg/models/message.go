package models

// Message contains the received message bytes.
type Message struct {
	Content []byte
}

// ReceivedDataUnit contains the received data unit,
// that has a string property to indicate the index it belogns to, and the data content.
type ReceivedDataUnit struct {
	IndexName string
	Data      []byte
}

// DataUnitToSend contains the data to send to ES.
type DataUnitToSend struct {
	Content []byte
}
