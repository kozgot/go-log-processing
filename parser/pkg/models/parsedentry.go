package models

import (
	"encoding/json"
	"fmt"
	"time"
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
	if err != nil {
		fmt.Println("Can't serialize", p)
	}

	return bytes
}
