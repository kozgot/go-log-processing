package contentparser

import (
	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type ErrorParser struct {
	line models.EntryWithLevelAndTimestamp
}

func NewErrorParser(line models.EntryWithLevelAndTimestamp) *ErrorParser {
	return &ErrorParser{line: line}
}

// ParseError parses an error level log entry.
func (e *ErrorParser) ParseError() *models.ErrorParams {
	errorCode := e.parseErrorCode()
	errorMessage := e.parseErrorMessage()
	errorSeverity := e.parseErrorSeverity()
	errorDesc := e.parseErrorDesc()
	errorSource := e.parseErrorSource()

	errorParams := models.ErrorParams{
		ErrorCode:   errorCode,
		Message:     errorMessage,
		Severity:    errorSeverity,
		Description: errorDesc,
		Source:      errorSource}

	return &errorParams
}

func (e *ErrorParser) parseErrorCode() int {
	return common.TryParseIntFromString(common.ParseFieldInBracketsAsString(e.line.Rest, formats.ErrorCodeRegex))
}

func (e *ErrorParser) parseErrorMessage() string {
	return common.ParseFieldInBracketsAsString(e.line.Rest, formats.MessageRegex)
}

func (e *ErrorParser) parseErrorDesc() string {
	return common.ParseFieldInBracketsAsString(e.line.Rest, formats.ErrorDescRegex)
}

func (e *ErrorParser) parseErrorSource() string {
	return common.ParseFieldInBracketsAsString(e.line.Rest, formats.ErrorSourceRegex)
}

func (e *ErrorParser) parseErrorSeverity() int {
	return common.TryParseIntFromString(common.ParseFieldInBracketsAsString(e.line.Rest, formats.ErrorSeverityRegex))
}
