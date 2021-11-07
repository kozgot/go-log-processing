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
	ErrorParams   *ErrorParams
	WarningParams *WarningParams
	InfoParams    *InfoParams
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

func (p *ParsedLogEntry) Equals(other ParsedLogEntry) bool {
	if p.Level != other.Level || p.Timestamp != other.Timestamp {
		return false
	}

	if p.ErrorParams == nil && other.ErrorParams != nil || p.ErrorParams != nil && other.ErrorParams == nil {
		return false
	}

	if p.WarningParams == nil && other.WarningParams != nil || p.WarningParams != nil && other.WarningParams == nil {
		return false
	}

	if p.InfoParams == nil && other.InfoParams != nil || p.InfoParams != nil && other.InfoParams == nil {
		return false
	}

	if p.ErrorParams != nil && !p.ErrorParams.Equals(*other.ErrorParams) {
		return false
	}

	if p.WarningParams != nil && !p.WarningParams.Equals(*other.WarningParams) {
		return false
	}

	if p.InfoParams != nil && !p.InfoParams.Equals(*other.InfoParams) {
		return false
	}

	return true
}
