package parsecontents

import (
	"regexp"
	"strconv"
	"strings"
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
		parsedNumber, err := strconv.ParseInt(stringRepresentation, 10, 64)
		if err != nil {
			panic(err)
		}

		return parsedNumber
	}

	return 0
}

func tryParseFloat64FromString(stringRepresentation string) float64 {
	if stringRepresentation != "" {
		parsedNumber, err := strconv.ParseFloat(stringRepresentation, 64)
		if err != nil {
			panic(err)
		}

		return parsedNumber
	}

	return 0
}
