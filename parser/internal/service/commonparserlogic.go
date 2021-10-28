package service

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/parser/pkg/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func parseFieldInBracketsAsString(line string, regexString string) string {
	regex, _ := regexp.Compile(regexString)
	textualField := regex.FindString(line)

	if textualField == "" {
		// log.Println("Could not parse textual field from line: ", line, regex)
		return ""
	}

	textualFieldValue := strings.Split(textualField, "[")[1]
	textualFieldValue = strings.Replace(textualFieldValue, "]", "", 1)

	return textualFieldValue
}

func parseFieldInDoubleBracketsAsString(line string, regexString string) string {
	regex, _ := regexp.Compile(regexString)
	textualField := regex.FindString(line)

	if textualField == "" {
		return ""
	}

	textualFieldValue := strings.Split(textualField, "[[")[1]
	textualFieldValue = strings.Replace(textualFieldValue, "]]", "", 1)

	return textualFieldValue
}

func parseFieldInParenthesesAsString(line string, regexString string) string {
	regex, _ := regexp.Compile(regexString)
	textualField := regex.FindString(line)

	if textualField == "" {
		// log.Println("Could not parse textual field from line: ", line, regex)
		return ""
	}

	textualFieldValue := strings.Split(textualField, "(")[1]
	textualFieldValue = strings.Replace(textualFieldValue, ")", "", 1)

	return textualFieldValue
}

func parseFieldAsString(line string, regexString string) string {
	regex, _ := regexp.Compile(regexString)
	textualField := regex.FindString(line)

	if textualField == "" {
		// log.Println("Could not parse textual field from line: ", line, regex)
		return ""
	}

	return textualField
}

func tryParseIntFromString(stringRepresentation string) int {
	if stringRepresentation != "" {
		parsedNumber, err := strconv.Atoi(stringRepresentation)
		if err != nil {
			panic(err)
		}

		return parsedNumber
	}

	return 0
}

func tryParseInt64FromString(stringRepresentation string) int64 {
	if stringRepresentation != "" {
		base := 10
		bitSize := 64
		parsedNumber, err := strconv.ParseInt(stringRepresentation, base, bitSize)
		if err != nil {
			panic(err)
		}

		return parsedNumber
	}

	return 0
}

func tryParseFloat64FromString(stringRepresentation string) float64 {
	if stringRepresentation != "" {
		bitSize := 64
		parsedNumber, err := strconv.ParseFloat(stringRepresentation, bitSize)
		if err != nil {
			panic(err)
		}

		return parsedNumber
	}

	return 0
}

func parseDateTimeField(line string, regex string) time.Time {
	timeFieldRegex, _ := regexp.Compile(regex)
	timeField := timeFieldRegex.FindString(line)

	if timeField != "" {
		timeString := strings.Split(timeField, "[")[1]
		timeString = strings.Replace(timeString, "]", "", 1)

		dateTime := parseDateTime(timeString)
		return dateTime
	}

	return time.Time{}
}

func parseDateTime(timeString string) time.Time {
	dateRegex, _ := regexp.Compile(formats.DateFormatRegex)

	dateString := dateRegex.FindString(timeString)
	if dateString != "" {
		date, err := time.Parse(formats.DateLayoutString, dateString)
		if err != nil {
			// Do not die here, log instead
			panic(err)
		}

		return date
	}

	return time.Time{}
}

func parseTimeFieldFromSeconds(line string, timeStampRegex string) time.Time {
	seconds := tryParseInt64FromString(parseFieldInBracketsAsString(line, timeStampRegex))
	if seconds != 0 {
		dateTimeFromsSecs := time.Unix(seconds, 0)
		return dateTimeFromsSecs
	}

	return time.Time{}
}

func parseTimeFieldFromMilliSeconds(line string, timeStampRegex string) time.Time {
	milliseconds := tryParseInt64FromString(parseFieldInBracketsAsString(line, timeStampRegex))
	if milliseconds != 0 {
		dateTimeFromsSecs := time.Unix(0, convertMillisecondsToSeconds(milliseconds))
		return dateTimeFromsSecs
	}

	return time.Time{}
}

func parseTimeRange(line string) *models.TimeRange {
	from := parseDateTimeField(line, formats.TimeRangeFromRegex)
	if !isValidDate(from) {
		// The from field is in a different format, try using that.
		from = parseTimeFieldFromSeconds(line, formats.TimeRangeStartTicksRegex)
	}

	to := parseDateTimeField(line, formats.TimeRangeToRegex)
	if !isValidDate(to) {
		// The to field is in a different format, try using that.
		to = parseTimeFieldFromSeconds(line, formats.TimeRangeEndTicksRegex)
	}

	if isValidDate(from) && isValidDate(to) {
		result := models.TimeRange{From: from, To: to}
		return &result
	}

	return nil
}

func isValidDate(date time.Time) bool {
	return (date.Year() > 1500)
}

func convertMillisecondsToSeconds(milliseconds int64) int64 {
	return milliseconds * 1000 * 1000
}
