package service

import (
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// ParseContents extracts the params of an error level line in the log file
func ParseContents(line models.LineWithDate) *models.ParsedLogEntry {
	parsedLine := models.ParsedLogEntry{
		Level:     line.Level,
		Timestamp: line.Timestamp}

	switch line.Level {
	case "ERROR":
		// todo: add these methods to the appropriate types,
		// eg.: ErrorParams should have a method parse() that parses error params
		errorParams := parseError(line)
		parsedLine.ErrorParams = *errorParams

	case "WARN":
		warning := parseWarn(line)
		if warning == nil {
			return nil
		}
		parsedLine.WarningParams = *warning

	// Log entries with 'WARNING' log level come from a different log file,
	// and they have a completely different format, so they are handled separately.
	case "WARNING":
		warning := parseWarning(line)
		parsedLine.WarningParams = *warning

	case "INFO":
		info := parseInfo(line)
		if info == nil {
			return nil
		}
		parsedLine.InfoParams = *info
	}

	return &parsedLine
}
