package models

import "time"

// IndexValue contains index data in a given time range.
type IndexValue struct {
	ReceiveTime   time.Time
	PreviousTime  time.Time
	Time          time.Time
	PreviousValue int
	Value         int
	ServiceLevel  int
	PodUID        string
	SerialNumber  int
	SmcUID        string
}
