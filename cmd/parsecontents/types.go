package parsecontents

import "time"

// ErrorParams contains the parsed error parameters
type ErrorParams struct {
	ErrorCode   int
	Message     string
	Severity    int
	Description string
	Source      string
}

// ParsedLine contains a parsed line from the log file
type ParsedLine struct {
	Timestamp   time.Time
	Level       string
	ErrorParams *ErrorParams
}
