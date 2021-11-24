package models

import (
	"encoding/json"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
)

// ReceivedDataUnit contains the received data unit,
// that has a string property to indicate the index it belogns to, and the data content.
type ReceivedDataUnit struct {
	IndexName string
	Data      []byte
}

func (d *ReceivedDataUnit) ToJSON() []byte {
	bytes, err := json.MarshalIndent(d, "", " ")
	utils.FailOnError(err, "Cannot serialize data unit")

	return bytes
}

func (d *ReceivedDataUnit) FromJSON(bytes []byte) {
	err := json.Unmarshal(bytes, d)
	utils.FailOnError(err, "Failed to unmarshal data unit")
}
