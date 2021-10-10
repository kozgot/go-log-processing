package service

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
