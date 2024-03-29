package models

import (
	"encoding/json"
	"time"

	"github.com/kozgot/go-log-processing/postprocessor/internal/utils"
)

// ConsumtionValue contains cunsumption data in a given time range.
type ConsumtionValue struct {
	ReceiveTime  time.Time
	StartTime    time.Time
	EndTime      time.Time
	Value        int
	ServiceLevel int
	SmcUID       string
}

// Serialize serlializes a consumption value to JSON format and returns a byte array.
func (c *ConsumtionValue) Serialize() []byte {
	bytes, err := json.Marshal(c)
	utils.FailOnError(err, "Can't serialize consumption value.")
	return bytes
}

// Deserialize deserializes a consumption value.
func (c *ConsumtionValue) Deserialize(bytes []byte) {
	err := json.Unmarshal(bytes, c)
	utils.FailOnError(err, "Cannot deserialize consumption value.")
}
