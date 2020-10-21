package parsedates

import (
	"regexp"
	"strings"
	"time"

	filter "github.com/kozgot/go-log-processing/cmd/filterlines"
)

// ParseDate returns a date parsed from the input (a line of the currently processed log file).
func ParseDate(line filter.Line) (*LineWithDate, bool) {
	dateRegex, _ := regexp.Compile(DateFormatRegex)

	dateString := dateRegex.FindString(line.Rest)
	if dateString != "" {
		date, err := time.Parse(DateLayoutString, dateString)
		if err != nil {
			// Do not die here, log instead
			panic(err)
		}
		restOfLine := strings.Replace(line.Rest, dateString, "", 1)
		return &LineWithDate{Timestamp: date, Rest: restOfLine, Level: line.Level}, true
	}

	// could not parse date, should log this event
	return nil, false
}
