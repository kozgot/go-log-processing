package processing

import (
	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// ProcessError processes a log entry with ERROR log level.
func ProcessError(logEntry parsermodels.ParsedLine) *models.SmcEntry {
	result := models.SmcEntry{}
	errorParams := logEntry.ErrorParams
	result.TimeStamp = logEntry.Timestamp
	result.EventType = logEntry.Level
	result.UID = errorParams.Source

	// todo more params?
	return &result
}
