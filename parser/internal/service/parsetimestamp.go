package service

import (
	"regexp"
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/internal/utils"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// ParseTimestamp returns a date parsed from the input (a line of the currently processed log file).
func ParseTimestamp(line models.Line) (*models.EntryWithLevelAndTimestamp, bool) {
	dateRegex, _ := regexp.Compile(formats.DateFormatRegex)
	dateRegexshort, _ := regexp.Compile(formats.DateFormatRegexShort)

	dateString := dateRegex.FindString(line.Rest)
	if dateString != "" {
		date, err := time.ParseInLocation(formats.DateLayoutString, dateString, time.UTC)
		utils.FailOnError(err, "Could not parse long date format")

		restOfLine := removeParsedParts(line.Rest, dateString)
		return &models.EntryWithLevelAndTimestamp{Timestamp: date, Rest: restOfLine, Level: line.Level}, true
	}

	dateString = dateRegexshort.FindString(line.Rest)
	if dateString != "" {
		date, err := time.ParseInLocation(formats.DateLayoutStringShort, dateString, time.UTC)
		utils.FailOnError(err, "Could not parse short date format")

		restOfLine := removeParsedParts(line.Rest, dateString)
		return &models.EntryWithLevelAndTimestamp{Timestamp: date, Rest: restOfLine, Level: line.Level}, true
	}

	return nil, false
}

func removeParsedParts(line string, parsedPart string) (rest string) {
	surroundingRegex, _ := regexp.Compile(formats.DateSurroundingRegex)
	restOfLine := strings.Replace(line, parsedPart, "", 1)

	// remove all tailing leftover square brackets, whitespaces and colons if present
	sourroundingString := surroundingRegex.FindString(restOfLine)
	if sourroundingString != "" {
		restOfLine = strings.TrimLeft(restOfLine, sourroundingString)
	}
	return restOfLine
}
