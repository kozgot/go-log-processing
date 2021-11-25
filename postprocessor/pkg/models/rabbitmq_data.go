package models

import (
	"encoding/json"

	"github.com/kozgot/go-log-processing/postprocessor/internal/utils"
)

// DataUnit contains the sent data unit.
type DataUnit struct {
	DataType DataType
	Data     []byte
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

// DataType represents the type of the data unit published to RabbitMQ.
type DataType int64

const (
	// UnknownDataType is the default value of DataType.
	UnknownDataType DataType = iota
	Event
	Consumption
)
