package testmodels

import (
	"encoding/json"
	"log"

	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// TestProcessedData contains processed test data.
type TestProcessedData struct {
	IndexNames   []string
	Events       []models.SmcEvent
	Consumptions []models.ConsumtionValue
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
