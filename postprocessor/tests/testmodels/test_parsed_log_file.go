package testmodels

import (
	"encoding/json"
	"log"

	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type TestParsedLogFile struct {
	Lines []models.ParsedLogEntry
}

func (t *TestParsedLogFile) ToJSON() []byte {
	bytes, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		log.Fatal("Can't serialize test entry", t)
	}

	return bytes
}

func (t *TestParsedLogFile) FromJSON(bytes []byte) {
	err := json.Unmarshal(bytes, t)
	if err != nil {
		log.Fatal("Can't deserialize", bytes)
	}
}
