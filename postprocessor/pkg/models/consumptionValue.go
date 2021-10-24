package models

import "time"

// ConsumtionValue contains cunsumption data in a given time range.
type ConsumtionValue struct {
	ReceiveTime  time.Time
	StartTime    time.Time
	EndTime      time.Time
	Value        int
	ServiceLevel int
	SmcUID       string
}
