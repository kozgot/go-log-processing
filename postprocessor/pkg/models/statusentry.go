package models

import "time"

// SmcEntry contains data from a log entry related to a specific SMC.
type StatusEntry struct {
	TimeStamp  time.Time // the exact time the event happened
	Message    string
	StatusByte string
}
