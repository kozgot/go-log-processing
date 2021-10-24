package processing

import (
	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// ProcessWarn processes a log entry with WARN log level.
func ProcessWarn(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	if logEntry.WarningParams.TimeoutParams.Protocol != "" && logEntry.WarningParams.TimeoutParams.URL != "" {
		address := models.AddressDetails{
			URL: logEntry.WarningParams.TimeoutParams.URL,
		}

		data := models.SmcData{
			Address: address,
		}

		event := models.SmcEvent{
			Time:            logEntry.Timestamp,
			EventType:       models.TimeoutWarning,
			EventTypeString: models.EventTypeToString(models.TimeoutWarning),
			Label:           "Timeout for URL " + logEntry.WarningParams.TimeoutParams.URL,
			DataPayload:     data,
		}

		return &data, &event
	}

	// the other warn type entries are not interesting for us right now
	return nil, nil
}

// ProcessWarning processes a log entry with WARNING log level.
func ProcessWarning(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	address := models.AddressDetails{
		PhysicalAddress: logEntry.WarningParams.JoinMessageParams.SmcAddress.PhysicalAddress,
		LogicalAddress:  logEntry.WarningParams.JoinMessageParams.SmcAddress.LogicalAddress,
		ShortAddress:    logEntry.WarningParams.JoinMessageParams.SmcAddress.ShortAddress,
	}
	smcUID := logEntry.WarningParams.JoinMessageParams.SmcAddress.SmcUID

	data := models.SmcData{
		SmcUID:  smcUID,
		Address: address,
	}

	event := models.SmcEvent{
		Time:            logEntry.Timestamp,
		EventTypeString: models.EventTypeToString(models.JoinRejectedWarning),
		EventType:       models.JoinRejectedWarning,
		Label:           "SMC join rejected for " + smcUID,
		SmcUID:          smcUID,
		DataPayload:     data,
	}

	return &data, &event
}
