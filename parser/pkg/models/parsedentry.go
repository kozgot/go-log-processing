package models

import (
	"encoding/json"
	"time"

	"github.com/kozgot/go-log-processing/parser/internal/utils"
)

// ParsedLogEntry contains a parsed line from the log file.
type ParsedLogEntry struct {
	Timestamp     time.Time
	Level         string
	ErrorParams   ErrorParams
	WarningParams WarningParams
	InfoParams    InfoParams
}

// Serialize serialzes a parsed log enrty.
func (p *ParsedLogEntry) Serialize() []byte {
	bytes, err := json.Marshal(p)
	utils.FailOnError(err, "Can't serialize parsed log entry")
	return bytes
}

// Serialize serialzes a parsed log enrty.
func (p *ParsedLogEntry) FromJSON(bytes []byte) {
	err := json.Unmarshal(bytes, p)
	utils.FailOnError(err, "Failed to unmarshal log entry")
}
