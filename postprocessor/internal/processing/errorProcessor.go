package processing

import (
	"strconv"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// ProcessError processes a log entry with ERROR log level.
func ProcessError(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	if logEntry.ErrorParams.Source == "" {
		// We do not know the source, so we do not know what to pair it with...
		return nil, nil
	}

	smcUID := logEntry.ErrorParams.Source
	label := "Error type " + logEntry.ErrorParams.Message + ", severity: " + strconv.Itoa(logEntry.ErrorParams.Severity)
	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.DLMSError,
		Label:     label,
	}

	data := models.SmcData{SmcUID: smcUID}
	return &data, &event
}
