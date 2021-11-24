package contentparser

import (
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// ParseEntryContents extracts the custom contents of a log entry.
func ParseEntryContents(line models.EntryWithLevelAndTimestamp) *models.ParsedLogEntry {
	parsedLine := models.ParsedLogEntry{
		Level:     line.Level,
		Timestamp: line.Timestamp}

	infoParser := NewInfoParser(line)
	warningParser := NewWarningParser(line)
	errorParser := NewErrorParser(line)

	switch line.Level {
	case "ERROR":
		errorParams := errorParser.ParseError()
		parsedLine.ErrorParams = errorParams

	case "WARN":
		warning := warningParser.ParseWarn()
		if warning == nil {
			return nil
		}
		parsedLine.WarningParams = warning

	// Log entries with 'WARNING' log level come from a different log file,
	// and they have a completely different format, so they are handled separately.
	case "WARNING":
		warning := warningParser.ParseWarning()
		parsedLine.WarningParams = warning

	case "INFO":
		info := infoParser.ParseInfo()
		parsedLine.InfoParams = info
	}

	return &parsedLine
}
