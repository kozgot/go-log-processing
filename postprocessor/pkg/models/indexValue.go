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

// Equals checks equality.
func (i *IndexValue) Equals(other IndexValue) bool {
	return i.ReceiveTime == other.ReceiveTime &&
		i.PreviousTime == other.PreviousTime &&
		i.Time == other.Time &&
		i.PreviousValue == other.PreviousValue &&
		i.Value == other.Value &&
		i.ServiceLevel == other.ServiceLevel &&
		i.PodUID == other.PodUID &&
		i.SerialNumber == other.SerialNumber &&
		i.SmcUID == other.SmcUID
}
