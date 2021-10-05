package processing

import (
	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// ProcessWarn processes a log entry with WARN log level.
func ProcessWarn(logEntry parsermodels.ParsedLine) *models.SmcEntry {
	result := models.SmcEntry{}
	result.TimeStamp = logEntry.Timestamp
	result.EventType = logEntry.Level
	result.UID = logEntry.WarningParams.SmcUID

	// todo more params?
	return &result
}

// ProcessWarning processes a log entry with WARNING log level.
func ProcessWarning(logEntry parsermodels.ParsedLine) *models.SmcEntry {
	result := models.SmcEntry{}
	result.TimeStamp = logEntry.Timestamp
	result.EventType = logEntry.Level
	result.UID = logEntry.WarningParams.JoinMessageParams.SmcAddress.SmcUID

	// todo
	return &result
}
