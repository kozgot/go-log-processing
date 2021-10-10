package models

import "time"

// ParsedLogEntry contains a parsed line from the log file.
type ParsedLogEntry struct {
	Timestamp     time.Time
	Level         string
	ErrorParams   ErrorParams
	WarningParams WarningParams
	InfoParams    InfoParams
}
