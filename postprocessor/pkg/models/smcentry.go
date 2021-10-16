package models

import "time"

// SmcEntry contains data from a log entry related to a specific SMC.
type SmcEntry struct {
	UID       string    // the UID of the SMC.
	TimeStamp time.Time // the exact time the event happened
	EventType string    // the type of the event (warning, info, error)
	Message   string
}

// ErrorDetails contains details regarding an error level log entry.
type ErrorDetails struct {
	ErrorMessage string
	ErrorCode    int
	Severity     int
}

// InfoDetails contains details regarding an info level log entry.
type InfoDetails struct {
	Type string // DC or SMC join
}

// WarningDetails contains details regarding a warning level log entry.
type WarningDetails struct {
	Priority     int
	Name         string
	ErrorMessage string
	ErrorCode    int
	Severity     int
}
