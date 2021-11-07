package processing

import (
	"strconv"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// ProcessError processes a log entry with ERROR log level.
func ProcessError(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	if logEntry.ErrorParams == nil {
		return nil, nil
	}

	if logEntry.ErrorParams.Source == "" {
		// We do not know the source, so we do not know what to pair it with...
		return nil, nil
	}

	smcUID := logEntry.ErrorParams.Source
	data := models.SmcData{SmcUID: smcUID}

	label := "Error type " + logEntry.ErrorParams.Message + ", severity: " + strconv.Itoa(logEntry.ErrorParams.Severity)
	event := models.SmcEvent{
		Time:            logEntry.Timestamp,
		EventType:       models.DLMSError,
		EventTypeString: models.EventTypeToString(models.DLMSError),
		Label:           label,
		SmcUID:          smcUID,
		DataPayload:     data,
	}

	return &data, &event
}
