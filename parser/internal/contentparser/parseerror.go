package contentparser

import (
	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// ParseError parses an error level log entry.
func ParseError(line models.EntryWithLevelAndTimestamp) *models.ErrorParams {
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
	return common.TryParseIntFromString(common.ParseFieldInBracketsAsString(line, formats.ErrorCodeRegex))
}

func parseErrorMessage(line string) string {
	return common.ParseFieldInBracketsAsString(line, formats.MessageRegex)
}

func parseErrorDesc(line string) string {
	return common.ParseFieldInBracketsAsString(line, formats.ErrorDescRegex)
}

func parseErrorSource(line string) string {
	return common.ParseFieldInBracketsAsString(line, formats.ErrorSourceRegex)
}

func parseErrorSeverity(line string) int {
	return common.TryParseIntFromString(common.ParseFieldInBracketsAsString(line, formats.ErrorSeverityRegex))
}
