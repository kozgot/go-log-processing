package processing

import (
	"fmt"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// ProcessInfoEntry processes a log entry with INFO log level.
func ProcessInfoEntry(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	switch logEntry.InfoParams.EntryType {
	case parsermodels.Routing:
		// this case is not really interesting for our logic
		return nil, nil
	case parsermodels.NetworkStatus:
		// this case is not really interesting for our logic
		return nil, nil

	case parsermodels.SMCJoin:
		smcData, event := processJoinEntry(logEntry)
		if smcData != nil {
			fmt.Println(smcData.SmcUID)
		}
		return smcData, event

	case parsermodels.DCMessage:
		data, event := processDCMessageEntry(logEntry)
		return data, event

	case parsermodels.ConnectionAttempt:
		data, event := processConnectionAttempt(logEntry)
		return data, event

	case parsermodels.ConnectionReleased:
		data, event := processConnectionReleased(logEntry)
		return data, event

	case parsermodels.InitDLMSConnection:
		data, event := processInitDLMSConnection(logEntry)
		return data, event

	case parsermodels.InternalDiagnostics:
		data, event := processInternalDiagnostics(logEntry)
		return data, event

	case parsermodels.SmcConfigUpdate:
		data, event := processSmcConfigUpdate(logEntry)
		return data, event

	// Unrecognized entry type
	case parsermodels.UnknownEntryType:
		return nil, nil
	default:
		break
	}

	return nil, nil
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
	}

	return &result, &event
}

func processDCMessageEntry(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	messageType := logEntry.InfoParams.DCMessage.MessageType
	switch messageType {
	case parsermodels.Connect:
		data, event := processConnect(logEntry)
		return data, event

	case parsermodels.DLMSLogs:
		data, event := processDLMSLogsEntry(logEntry)
		return data, event

	case parsermodels.IndexHighProfileGeneric:
		data, event := processIndexHighProfileGeneric(logEntry)
		return data, event

	case parsermodels.IndexLowProfileGeneric:
		data, event := processIndexLowProfileGeneric(logEntry)
		return data, event

	case parsermodels.ReadIndexLowProfiles:
		data, event := processReadIndexLowProfiles(logEntry)
		return data, event

	case parsermodels.ReadIndexProfiles:
		data, event := processReadIndexProfiles(logEntry)
		return data, event

	case parsermodels.IndexReceived:
		data, event := processIndexReceived(logEntry)
		return data, event

	case parsermodels.NewSmc:
		data, event := processNewSmc(logEntry)
		return data, event

	case parsermodels.MessageSentToSVI:
		data, event := processSVIMessage(logEntry)
		return data, event

	case parsermodels.PodConfig:
		data, event := processPodConfig(logEntry)
		return data, event

	case parsermodels.SmcConfig:
		data, event := processSmcConfig(logEntry)
		return data, event

	case parsermodels.SmcAddress:
		data, event := processSmcAddress(logEntry)
		return data, event

	case parsermodels.ServiceLevel:
		data, event := processServicelevelEntry(logEntry)
		return data, event

	case parsermodels.Settings:
		data, event := processSettings(logEntry)
		return data, event

	case parsermodels.Statistics:
		data, event := processStatistics(logEntry)
		return data, event

	case parsermodels.UnknownDCMessage:
		return nil, nil
	default:
		return nil, nil
	}
}

func processIndexLowProfileGeneric(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processIndexHighProfileGeneric(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processDLMSLogsEntry(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processReadIndexProfiles(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processReadIndexLowProfiles(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processIndexReceived(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processNewSmc(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processSVIMessage(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processPodConfig(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processConnect(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processSettings(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processServicelevelEntry(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processSmcAddress(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processSmcConfig(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo implementation
	return nil, nil
}

func processStatistics(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	statisticsPayload := logEntry.InfoParams.DCMessage.Payload.StatisticsEntryPayload

	if statisticsPayload == nil {
		return nil, nil
	}

	smcUID := statisticsPayload.SourceID
	data := models.SmcData{
		SmcUID: smcUID,
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.StatisticsSent,
		Label:     "Statistics sent to SVI (" + statisticsPayload.Type + ")",
	}

	return &data, &event
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
	}

	return &data, &event
}
