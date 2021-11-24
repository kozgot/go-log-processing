package common

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/internal/utils"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// ParseFieldInBracketsAsString parses a field in a log entry surrounded by brackets.
func ParseFieldInBracketsAsString(line string, regexString string) string {
	regex, _ := regexp.Compile(regexString)
	textualField := regex.FindString(line)

	if textualField == "" {
		return ""
	}

	textualFieldValue := strings.Split(textualField, "[")[1]
	textualFieldValue = strings.Replace(textualFieldValue, "]", "", 1)

	return textualFieldValue
}

// ParseFieldInDoubleBracketsAsString parses a field in a log entry surrounded by double brackets.
func ParseFieldInDoubleBracketsAsString(line string, regexString string) string {
	regex, _ := regexp.Compile(regexString)
	textualField := regex.FindString(line)

	if textualField == "" {
		return ""
	}

	textualFieldValue := strings.Split(textualField, "[[")[1]
	textualFieldValue = strings.Replace(textualFieldValue, "]]", "", 1)

	return textualFieldValue
}

// ParseFieldInParenthesesAsString parses a field in a log entry surrounded by parentheses.
func ParseFieldInParenthesesAsString(line string, regexString string) string {
	regex, _ := regexp.Compile(regexString)
	textualField := regex.FindString(line)

	if textualField == "" {
		return ""
	}

	textualFieldValue := strings.Split(textualField, "(")[1]
	textualFieldValue = strings.Replace(textualFieldValue, ")", "", 1)

	return textualFieldValue
}

// ParseFieldAsString parses a field of a log entry as string.
func ParseFieldAsString(line string, regexString string) string {
	regex, _ := regexp.Compile(regexString)
	textualField := regex.FindString(line)

	if textualField == "" {
		// log.Println("Could not parse textual field from line: ", line, regex)
		return ""
	}

	return textualField
}

// TryParseIntFromString parses an integer value from a string representation.
func TryParseIntFromString(stringRepresentation string) int {
	if stringRepresentation != "" {
		parsedNumber, err := strconv.Atoi(stringRepresentation)
		utils.FailOnError(err, "Could not parse integer value from string representation.")
		return parsedNumber
	}

	return 0
}

// TryParseInt64FromString parses an int64 value from a string representation.
func TryParseInt64FromString(stringRepresentation string) int64 {
	if stringRepresentation != "" {
		parsedNumber, err := strconv.ParseInt(stringRepresentation, 10, 64)
		utils.FailOnError(err, "Could not parse int64 value from string representation.")
		return parsedNumber
	}

	return 0
}

// TryParseFloat64FromString parses a float value from a string representation.
func TryParseFloat64FromString(stringRepresentation string) float64 {
	if stringRepresentation != "" {
		parsedNumber, err := strconv.ParseFloat(stringRepresentation, 64)
		utils.FailOnError(err, "Could not parse float value from string representation.")
		return parsedNumber
	}

	return 0
}

// ParseDateTimeField parses a datetime field sorrounded by brackets.
func ParseDateTimeField(line string, regex string) time.Time {
	timeFieldRegex, _ := regexp.Compile(regex)
	timeField := timeFieldRegex.FindString(line)

	if timeField != "" {
		timeString := strings.Split(timeField, "[")[1]
		timeString = strings.Replace(timeString, "]", "", 1)

		dateTime := ParseDateTime(timeString)
		return dateTime
	}

	return time.Time{}
}

// ParseDateTime parses a datetime field.
func ParseDateTime(timeString string) time.Time {
	dateRegex, _ := regexp.Compile(formats.DateFormatRegex)

	dateString := dateRegex.FindString(timeString)
	if dateString != "" {
		date, err := time.ParseInLocation(formats.DateLayoutString, dateString, time.UTC)
		if err != nil {
			// Do not die here, log instead
			panic(err)
		}

		return date
	}

	return time.Time{}
}

// ParseTimeFieldFromSeconds parses a time field represented by seconds.
func ParseTimeFieldFromSeconds(line string, timeStampRegex string) time.Time {
	seconds := TryParseInt64FromString(ParseFieldInBracketsAsString(line, timeStampRegex))
	if seconds != 0 {
		dateTimeFromsSecs := time.Unix(seconds, 0)
		utcDatTime := time.Date(
			dateTimeFromsSecs.Year(),
			dateTimeFromsSecs.Month(),
			dateTimeFromsSecs.Day(),
			dateTimeFromsSecs.Hour(),
			dateTimeFromsSecs.Minute(),
			dateTimeFromsSecs.Second(),
			dateTimeFromsSecs.Nanosecond(),
			time.UTC)
		return utcDatTime
	}

	return time.Time{}
}

// ParseTimeFieldFromMilliSeconds parses a time field represented by milliseconds.
func ParseTimeFieldFromMilliSeconds(line string, timeStampRegex string) time.Time {
	milliseconds := TryParseInt64FromString(ParseFieldInBracketsAsString(line, timeStampRegex))
	if milliseconds != 0 {
		dateTimeFromsSecs := time.Unix(0, convertMillisecondsToSeconds(milliseconds))
		utcDatTime := time.Date(
			dateTimeFromsSecs.Year(),
			dateTimeFromsSecs.Month(),
			dateTimeFromsSecs.Day(),
			dateTimeFromsSecs.Hour(),
			dateTimeFromsSecs.Minute(),
			dateTimeFromsSecs.Second(),
			dateTimeFromsSecs.Nanosecond(),
			time.UTC)
		return utcDatTime
	}

	return time.Time{}
}

// ParseTimeRange parses a time range from a log entry.
func ParseTimeRange(line string) *models.TimeRange {
	from := ParseDateTimeField(line, formats.TimeRangeFromRegex)
	if !IsValidDate(from) {
		// The from field is in a different format, try using that.
		from = ParseTimeFieldFromSeconds(line, formats.TimeRangeStartTicksRegex)
	}

	to := ParseDateTimeField(line, formats.TimeRangeToRegex)
	if !IsValidDate(to) {
		// The to field is in a different format, try using that.
		to = ParseTimeFieldFromSeconds(line, formats.TimeRangeEndTicksRegex)
	}

	if IsValidDate(from) && IsValidDate(to) {
		result := models.TimeRange{From: from, To: to}
		return &result
	}

	return nil
}

// IsValidDate checks if the parsed date is valid by checking if the year value is big enough.
func IsValidDate(date time.Time) bool {
	return (date.Year() > 1500)
}

func convertMillisecondsToSeconds(milliseconds int64) int64 {
	return milliseconds * 1000 * 1000
}
