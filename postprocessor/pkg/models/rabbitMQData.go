package models

import (
	"encoding/json"

	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"
)

// DataUnit contains the sent data unit.
type DataUnit struct {
	IndexName string
	Data      []byte
}

// Serialize serlializes a data unit to JSON format and returns a byte array.
func (d DataUnit) Serialize() []byte {
	bytes, err := json.Marshal(d)
	utils.FailOnError(err, "  [RABBITMQ DATA] Can't serialize DatatUnit")

	return bytes
}

// Deserialize deserializes a data unit.
func (d *DataUnit) Deserialize(bytes []byte) {
	err := json.Unmarshal(bytes, d)
	utils.FailOnError(err, "Cannot deserialize data unit.")
}
