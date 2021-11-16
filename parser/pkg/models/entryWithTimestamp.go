package models

import "time"

// EntryWithLevelAndTimestamp contains a line of the input file with the parsed timestamp.
type EntryWithLevelAndTimestamp struct {
	Timestamp time.Time
	Level     string
	Rest      string
}
