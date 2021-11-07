package models

import "time"

// EntryWithLevelAndTimestamp contains a line of the input file with the parsed timestamp.
type EntryWithLevelAndTimestamp struct {
	Timestamp time.Time
	Level     string
	Rest      string
}

// Equals checks equality.
func (e *EntryWithLevelAndTimestamp) Equals(other EntryWithLevelAndTimestamp) bool {
	return e.Timestamp == other.Timestamp && e.Level == other.Level && e.Rest == other.Rest
}
