package service

import (
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func parseError(line models.LineWithDate) *models.ErrorParams {
	errorCode := parseErrorCode(line.Rest)
	errorMessage := parseErrorMessage(line.Rest)
	errorSeverity := parseErrorSeverity(line.Rest)
	errorDesc := parseErrorDesc(line.Rest)
	errorSource := parseErrorSource(line.Rest)

	errorParams := models.ErrorParams{
		ErrorCode:   errorCode,
		Message:     errorMessage,
		Severity:    errorSeverity,
		Description: errorDesc,
		Source:      errorSource}

	return &errorParams
}

func parseErrorCode(line string) int {
	return tryParseIntFromString(parseFieldInBracketsAsString(line, formats.ErrorCodeRegex))
}

func parseErrorMessage(line string) string {
	return parseFieldInBracketsAsString(line, formats.MessageRegex)
}

func parseErrorDesc(line string) string {
	return parseFieldInBracketsAsString(line, formats.ErrorDescRegex)
}

func parseErrorSource(line string) string {
	return parseFieldInBracketsAsString(line, formats.ErrorSourceRegex)
}

func parseErrorSeverity(line string) int {
	return tryParseIntFromString(parseFieldInBracketsAsString(line, formats.ErrorSeverityRegex))
}
