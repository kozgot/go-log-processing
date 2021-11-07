package models

// ErrorParams contains the parsed error parameters.
type ErrorParams struct {
	ErrorCode   int
	Message     string
	Severity    int
	Description string
	Source      string // If not empty, this is the UID of an smc.
}

// Equals checks equality.
func (e *ErrorParams) Equals(other ErrorParams) bool {
	if e.ErrorCode == other.ErrorCode &&
		e.Description == other.Description &&
		e.Message == other.Message &&
		e.Source == other.Source &&
		e.Severity == other.Severity {
		return true
	}

	return false
}
