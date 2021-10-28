package models

// Line contains a line of the input file with the parsed timestamp and log level.
type Line struct {
	Level string
	Rest  string
}
