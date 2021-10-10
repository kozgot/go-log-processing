package models

// ErrorParams contains the parsed error parameters.
type ErrorParams struct {
	ErrorCode   int
	Message     string
	Severity    int
	Description string
	Source      string // If not empty, this is the UID of an smc.
}
