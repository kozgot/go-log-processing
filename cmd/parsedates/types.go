package parsedates

import (
	"time"
)

// LineWithDate contains a line of the input file with the parsed timestamp
type LineWithDate struct {
	Timestamp time.Time
	Level     string
	Rest      string
}
