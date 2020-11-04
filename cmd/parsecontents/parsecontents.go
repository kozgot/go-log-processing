package parsecontents

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/kozgot/go-log-processing/cmd/parsedates"
)

// ParseContents extracts the params of an error level line in the log file
func ParseContents(line parsedates.LineWithDate) *ParsedLine {
	parsedLine := ParsedLine{Level: line.Level, Timestamp: line.Timestamp}
	switch line.Level {
	case "ERROR":
		errorParams := parseError(line)
		parsedLine.ErrorParams = errorParams

	case "INFO":
		// fmt.Println("INFO")
	case "WARN":
		// fmt.Println("WARN")
	}

	return &parsedLine
}

func parseError(line parsedates.LineWithDate) *ErrorParams {
	errorCode := parseErrorCode(line.Rest)
	errorMessage := parseErrorMessage(line.Rest)
	errorSeverity := parseErrorSeverity(line.Rest)
	errorDesc := parseErrorDesc(line.Rest)
	errorSource := parseErrorSource(line.Rest)

	errorParams := ErrorParams{ErrorCode: errorCode, Message: errorMessage, Severity: errorSeverity, Description: errorDesc, Source: errorSource}

	return &errorParams
}

func parseErrorCode(line string) int {
	errorCodeFieldRegex, _ := regexp.Compile(ErrorCodeRegex)
	errorCodeField := errorCodeFieldRegex.FindString(line)

	if errorCodeField != "" {
		errorCodeRegex, _ := regexp.Compile(NumberRegex)
		errorCodeString := errorCodeRegex.FindString(line)
		errorCode, err := strconv.Atoi(errorCodeString)
		if err != nil {
			panic(err)
		}

		return errorCode
	}

	return 0
}

func parseErrorMessage(line string) string {
	errorMessageFieldRegex, _ := regexp.Compile(MessageRegex)
	errorMessageField := errorMessageFieldRegex.FindString(line)

	if errorMessageField != "" {
		// todo: match this with regex
		errorMessageString := strings.Split(errorMessageField, "[")[1]
		errorMessageString = strings.Replace(errorMessageString, "]", "", 1)

		return errorMessageString
	}

	return ""
}

func parseErrorDesc(line string) string {
	errorDescFieldRegex, _ := regexp.Compile(ErrorDescRegex)
	errorDescField := errorDescFieldRegex.FindString(line)

	if errorDescField != "" {
		// todo: match this with regex
		errorDescString := strings.Split(errorDescField, "[")[1]
		errorDescString = strings.Replace(errorDescString, "]", "", 1)

		return errorDescString
	}

	return ""
}

func parseErrorSource(line string) string {
	errorSourceFieldRegex, _ := regexp.Compile(ErrorSourceRegex)
	errorSourceField := errorSourceFieldRegex.FindString(line)

	if errorSourceField != "" {
		// todo: match this with regex
		errorSourceString := strings.Split(errorSourceField, "[")[1]
		errorSourceString = strings.Replace(errorSourceString, "]", "", 1)

		return errorSourceString
	}

	return ""
}

func parseErrorSeverity(line string) int {
	errorSeverityFieldRegex, _ := regexp.Compile(ErrorSeverityRegex)
	errorSeverityField := errorSeverityFieldRegex.FindString(line)

	if errorSeverityField != "" {
		errorSeverityRegex, _ := regexp.Compile(NumberRegex)
		errorSeverityString := errorSeverityRegex.FindString(errorSeverityField)
		errorSeverity, err := strconv.Atoi(errorSeverityString)
		if err != nil {
			panic(err)
		}

		return errorSeverity
	}

	return 0
}
