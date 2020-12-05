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

// WarningParams contains the parsed warning parameters
type WarningParams struct {
	Name          string
	SMC_UID       string
	UID           int
	Priority      int
	Retry         int
	Creation      time.Time
	MinLaunchTime time.Time
	ErrorParams   *ErrorParams
	FileName      string
}

// ParsedLine contains a parsed line from the log file
type ParsedLine struct {
	Timestamp     time.Time
	Level         string
	ErrorParams   *ErrorParams
	WarningParams *WarningParams
}
