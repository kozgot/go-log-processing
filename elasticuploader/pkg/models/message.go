package models

// Message contains the received message bytes.
type Message struct {
	Content []byte
}

// DataUnitToSend contains the data to send to ES.
type DataUnitToSend struct {
	Content []byte
}
