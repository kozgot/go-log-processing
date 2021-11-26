package testmodels

import (
	"encoding/json"
	"log"

	postprocmodels "github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// TestProcessedData contains processed test data.
type TestProcessedData struct {
	Events       []postprocmodels.SmcEvent
	Consumptions []postprocmodels.ConsumtionValue
}

// ToJSON converts a TestProcessedData to json.
func (t *TestProcessedData) ToJSON() []byte {
	bytes, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		log.Fatal("Can't serialize test processed data", t)
	}

	return bytes
}

// FromJSON creates a TestProcessedData from json.
func (t *TestProcessedData) FromJSON(bytes []byte) {
	err := json.Unmarshal(bytes, t)
	if err != nil {
		log.Fatal("Can't deserialize", bytes)
	}
}
