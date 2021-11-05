package models

import (
	"encoding/json"
	"fmt"
)

// ReceivedDataUnit contains the received data unit,
// that has a string property to indicate the index it belogns to, and the data content.
type ReceivedDataUnit struct {
	IndexName string
	Data      []byte
}

// DeserializeDataUnit deserializes a received data unit.
func DeserializeDataUnit(dataBytes []byte) ReceivedDataUnit {
	var data ReceivedDataUnit
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		fmt.Println("Failed to unmarshal:", err)
	}

	return data
}
