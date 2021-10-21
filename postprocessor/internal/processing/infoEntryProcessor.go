package processing

import (
	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// ProcessInfoEntry processes a log entry with INFO log level.
func ProcessInfoEntry(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent, *models.ConsumtionValue, *models.IndexValue) {
	switch logEntry.InfoParams.EntryType {
	case parsermodels.Routing:
		// this case is not really interesting for our logic
		return nil, nil, nil, nil
	case parsermodels.NetworkStatus:
		// this case is not really interesting for our logic
		return nil, nil, nil, nil

	case parsermodels.SMCJoin:
		smcData, event := processJoinEntry(logEntry)
		return smcData, event, nil, nil

	case parsermodels.DCMessage:
		result := processDCMessageEntry(logEntry)
		return result.SmcData, result.SmcEvent, result.ConsumtionValue, result.IndexValue

	case parsermodels.ConnectionAttempt:
		data, event := processConnectionAttempt(logEntry)
		return data, event, nil, nil

	case parsermodels.ConnectionReleased:
		data, event := processConnectionReleased(logEntry)
		return data, event, nil, nil

	case parsermodels.InitDLMSConnection:
		data, event := processInitDLMSConnection(logEntry)
		return data, event, nil, nil

	case parsermodels.InternalDiagnostics:
		data, event := processInternalDiagnostics(logEntry)
		return data, event, nil, nil

	case parsermodels.SmcConfigUpdate:
		data, event := processSmcConfigUpdate(logEntry)
		return data, event, nil, nil

	// Unrecognized entry type
	case parsermodels.UnknownEntryType:
		return nil, nil, nil, nil
	default:
		break
	}

	return nil, nil, nil, nil
}

func processJoinEntry(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	smcUID := logEntry.InfoParams.JoinMessage.SmcAddress.SmcUID
	address := models.AddressDetails{
		PhysicalAddress: logEntry.InfoParams.JoinMessage.SmcAddress.PhysicalAddress,
		LogicalAddress:  logEntry.InfoParams.JoinMessage.SmcAddress.LogicalAddress,

		ShortAddress: logEntry.InfoParams.JoinMessage.SmcAddress.ShortAddress,
		URL:          "", // this is filled in later with another log entry

	}
	result := models.SmcData{
		Address:              address,
		SmcUID:               smcUID,
		CustomerSerialNumber: "", // this is filled in later with another log entry
		LastJoiningDate:      logEntry.InfoParams.JoinMessage.SmcAddress.LastJoiningDate,
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.SmcJoined,
		Label:     "Smc " + smcUID + " has joined",
		SmcUID:    smcUID,
	}

	return &result, &event
}

func processConnectionAttempt(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	smcUID := logEntry.InfoParams.ConnectionAttempt.SmcUID
	address := models.AddressDetails{
		URL: logEntry.InfoParams.ConnectionAttempt.URL,
	}
	data := models.SmcData{
		Address: address,
		SmcUID:  smcUID,
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.ConnectionAttempt,
		Label:     "Connection attempt to " + smcUID,
		SmcUID:    smcUID,
	}

	return &data, &event
}

func processConnectionReleased(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	address := models.AddressDetails{
		URL: logEntry.InfoParams.ConnectionReleased.URL,
	}
	data := models.SmcData{
		Address: address,
	}

	// todo: there is no smc uid here, should save URL for SMC id to be able to look it up

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.ConnectionReleased,
		Label:     "Released connection, URL: " + address.URL,
	}

	return &data, &event
}

func processInitDLMSConnection(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	address := models.AddressDetails{
		URL: logEntry.InfoParams.InitConnection.URL,
	}
	data := models.SmcData{
		Address: address,
	}

	// todo: there is no smc uid here, should save URL for SMC id to be able to look it up

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.InitConnection,
		Label:     "Initialize DLMS connection, URL: " + address.URL,
	}

	return &data, &event
}

func processInternalDiagnostics(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	smcUID := logEntry.InfoParams.InternalDiagnosticsData.SmcUID
	data := models.SmcData{
		SmcUID:                    smcUID,
		LastSuccesfulDlmsResponse: logEntry.InfoParams.InternalDiagnosticsData.LastSuccessfulDlmsResponseDate,
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.InternalDiagnostics,
		Label:     "Internal diagnostics...",
		SmcUID:    smcUID,
	}

	return &data, &event
}

func processSmcConfigUpdate(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	smcUID := logEntry.InfoParams.SmcConfigUpdate.SmcUID
	address := models.AddressDetails{
		PhysicalAddress: logEntry.InfoParams.SmcConfigUpdate.PhysicalAddress,
		LogicalAddress:  logEntry.InfoParams.SmcConfigUpdate.LogicalAddress,
		ShortAddress:    logEntry.InfoParams.SmcConfigUpdate.ShortAddress,
	}

	data := models.SmcData{
		Address:         address,
		SmcUID:          smcUID,
		LastJoiningDate: logEntry.InfoParams.SmcConfigUpdate.LastJoiningDate,
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.ConfigurationChanged,
		Label:     "Configuration update for " + smcUID,
		SmcUID:    smcUID,
	}

	return &data, &event
}
