package models

// EntryWithLogLevel contains a line of the input file with the parsed timestamp and log level.
type EntryWithLogLevel struct {
	Level string
	Rest  string
}

func (e *EntryWithLogLevel) Equals(other EntryWithLogLevel) bool {
	return e.Level == other.Level && e.Rest == other.Rest
}
