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
		smcData := processJoinEntry(logEntry)
		if smcData != nil {
			fmt.Println(smcData.SmcUID)
		}
		return smcData, nil

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

	case parsermodels.UnknownEntryType:
		return nil, nil
	default:
		break
	}

	return nil, nil
}

func processJoinEntry(logEntry parsermodels.ParsedLogEntry) *models.SmcData {
	address := models.AddressDetails{
		PhysicalAddress: logEntry.InfoParams.JoinMessage.SmcAddress.PhysicalAddress,
		LogicalAddress:  logEntry.InfoParams.JoinMessage.SmcAddress.LogicalAddress,

		ShortAddress: logEntry.InfoParams.JoinMessage.SmcAddress.ShortAddress,
		URL:          "", // this is filled in later with another log entry

	}
	result := models.SmcData{
		Address:              address,
		SmcUID:               logEntry.InfoParams.JoinMessage.SmcAddress.SmcUID,
		CustomerSerialNumber: "", // this is filled in later with another log entry
		LastJoiningDate:      logEntry.InfoParams.JoinMessage.SmcAddress.LastJoiningDate,
	}

	return &result
}

func processDCMessageEntry(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	address := models.AddressDetails{}
	data := models.SmcData{
		Address: address,
	}

	messageType := logEntry.InfoParams.DCMessage.MessageType
	switch messageType {
	case parsermodels.Connect:
		break

	case parsermodels.DLMSLogs:
		break

	case parsermodels.IndexHighProfileGeneric:
		break

	case parsermodels.IndexLowProfileGeneric:
		break

	case parsermodels.ReadIndexLowProfiles:
		break

	case parsermodels.ReadIndexProfiles:
		break

	case parsermodels.IndexReceived:
		break

	case parsermodels.NewSmc:
		break

	case parsermodels.MessageSentToSVI:
		break

	case parsermodels.PodConfig:
		break

	case parsermodels.SmcConfig:
		break

	case parsermodels.SmcAddress:
		break

	case parsermodels.ServiceLevel:
		break

	case parsermodels.Settings:
		break

	case parsermodels.Statistics:
		break

	case parsermodels.UnknownDCMessage:
		break

	default:
		break
	}

	event := models.SmcEvent{}

	// todo
	return &data, &event
}

func processConnectionAttempt(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	address := models.AddressDetails{
		URL: logEntry.InfoParams.ConnectionAttempt.URL,
	}
	data := models.SmcData{
		Address: address,
		SmcUID:  logEntry.InfoParams.ConnectionAttempt.SmcUID,
	}

	// todo
	event := models.SmcEvent{}
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

	// todo
	event := models.SmcEvent{}
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

	// todo
	event := models.SmcEvent{}
	return &data, &event
}

func processInternalDiagnostics(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	data := models.SmcData{
		SmcUID:                    logEntry.InfoParams.InternalDiagnosticsData.SmcUID,
		LastSuccesfulDlmsResponse: logEntry.InfoParams.InternalDiagnosticsData.LastSuccessfulDlmsResponseDate,
	}

	// todo
	event := models.SmcEvent{}
	return &data, &event
}

func processSmcConfigUpdate(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	address := models.AddressDetails{
		PhysicalAddress: logEntry.InfoParams.SmcConfigUpdate.PhysicalAddress,
		LogicalAddress:  logEntry.InfoParams.SmcConfigUpdate.LogicalAddress,
		ShortAddress:    logEntry.InfoParams.SmcConfigUpdate.ShortAddress,
	}

	data := models.SmcData{
		Address:         address,
		SmcUID:          logEntry.InfoParams.SmcConfigUpdate.SmcUID,
		LastJoiningDate: logEntry.InfoParams.SmcConfigUpdate.LastJoiningDate,
	}

	// todo
	event := models.SmcEvent{}
	return &data, &event
}
