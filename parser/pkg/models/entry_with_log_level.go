package models

// EntryWithLogLevel contains a line of the input file with the parsed timestamp and log level.
type EntryWithLogLevel struct {
	Level string
	Rest  string
}
