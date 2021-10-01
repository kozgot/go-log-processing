package models

import "time"

// SmcEntry contains data from a log entry related to a specific SMC.
type RoutingEntry struct {
	TimeStamp time.Time // the exact time the event happened
}
