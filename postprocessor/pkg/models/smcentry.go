package models

import "time"

// SmcEntry contains data from a log entry related to a specific SMC.
type SmcEntry struct {
	UID       string    // the UID of the SMC.
	TimeStamp time.Time // the exact time the event happened
	EventType string    // the type of the event (warning, info, error)
}
