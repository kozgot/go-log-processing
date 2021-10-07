package models

import "time"

// StatusEntry contains data from a log entry related to the network status.
type StatusEntry struct {
	TimeStamp  time.Time // the exact time the event happened
	Message    string
	StatusByte string
}
