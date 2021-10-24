package models

import (
	"encoding/json"
	"fmt"
)

// DataUnit contains the sent data unit.
type DataUnit struct {
	IndexName string
	Data      []byte
}

// Serialize serlializes a data unit to JSON format and returns a byte array.
func (d DataUnit) Serialize() []byte {
	bytes, err := json.Marshal(d)
	if err != nil {
		fmt.Println("Can't serialize", d)
	}

	return bytes
}
