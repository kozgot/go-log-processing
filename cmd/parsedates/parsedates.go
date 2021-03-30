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
	dateRegexshort, _ := regexp.Compile(DateFormatRegexShort)

	dateString := dateRegex.FindString(line.Rest)
	if dateString != "" {
		date, err := time.Parse(DateLayoutString, dateString)
		if err != nil {
			// Do not die here, log instead
			panic(err)
		}
		restOfLine := removeParsedParts(line.Rest, dateString)
		return &LineWithDate{Timestamp: date, Rest: restOfLine, Level: line.Level}, true
	}

	if dateString == "" {
		dateString = dateRegexshort.FindString(line.Rest)
		if dateString != "" {
			date, err := time.Parse(DateLayoutStringShort, dateString)
			if err != nil {
				// Do not die here, log instead
				panic(err)
			}

			restOfLine := removeParsedParts(line.Rest, dateString)
			return &LineWithDate{Timestamp: date, Rest: restOfLine, Level: line.Level}, true
		}
	}

	// could not parse date, should log this event
	return nil, false
}

func removeParsedParts(line string, parsedPart string) (rest string) {
	surroundingRegex, _ := regexp.Compile(DateSurroundingRegex)
	restOfLine := strings.Replace(line, parsedPart, "", 1)

	// remove all tailing leftover square brackets, whitespaces and colons if present
	sourroundingString := surroundingRegex.FindString(restOfLine)
	if sourroundingString != "" {
		restOfLine = strings.TrimLeft(restOfLine, sourroundingString)
	}
	return restOfLine
}
