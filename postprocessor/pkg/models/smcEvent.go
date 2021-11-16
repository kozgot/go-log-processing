package models

import (
	"encoding/json"
	"time"

	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"
)

// SmcEvent is an event happening to a specific smc at a specific time.
type SmcEvent struct {
	Time            time.Time
	EventType       EventType
	EventTypeString string
	Label           string
	SmcUID          string
	DataPayload     SmcData
}

// Serialize serializes an smc event and returns a byte array.
func (e *SmcEvent) Serialize() []byte {
	bytes, err := json.Marshal(e)
	utils.FailOnError(err, "Can't serialize smc event")

	return bytes
}

// Deserialize deserializes an smc event.
func (e *SmcEvent) Deserialize(bytes []byte) {
	err := json.Unmarshal(bytes, e)
	utils.FailOnError(err, "Cannot deserialize smc event.")
}
