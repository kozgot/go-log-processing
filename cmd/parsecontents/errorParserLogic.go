package parsecontents

import "github.com/kozgot/go-log-processing/cmd/parsedates"

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
	return tryParseIntFromString(parseFieldInBracketsAsString(line, ErrorCodeRegex))
}

func parseErrorMessage(line string) string {
	return parseFieldInBracketsAsString(line, MessageRegex)
}

func parseErrorDesc(line string) string {
	return parseFieldInBracketsAsString(line, ErrorDescRegex)
}

func parseErrorSource(line string) string {
	return parseFieldInBracketsAsString(line, ErrorSourceRegex)
}

func parseErrorSeverity(line string) int {
	return tryParseIntFromString(parseFieldInBracketsAsString(line, ErrorSeverityRegex))
}
