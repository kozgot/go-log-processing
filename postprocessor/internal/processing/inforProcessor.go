package processing

import (
	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// ProcessInfo processes a log entry with INFO log level.
func ProcessInfo(logEntry parsermodels.ParsedLogEntry) (*models.SmcEntry, *models.RoutingEntry, *models.StatusEntry) {
	// one of 'ROUTING', 'JOIN', 'STATUS', or 'DC'
	switch logEntry.InfoParams.EntryType {
	case parsermodels.Routing:
		routingEntry := processRoutingMessage(logEntry)
		return nil, routingEntry, nil
	case parsermodels.SMCJoin:
		joinEntry := processJoinMessage(logEntry)
		return joinEntry, nil, nil
	case parsermodels.NetworkStatus:
		statusEntry := processStatusMessage(logEntry)
		return nil, nil, statusEntry
	case parsermodels.DCMessage:
		dcMessage := processDCMessage(logEntry)
		return dcMessage, nil, nil
	default:
		break
	}

	return nil, nil, nil
}

func processDCMessage(logEntry parsermodels.ParsedLogEntry) *models.SmcEntry {
	result := models.SmcEntry{}
	result.TimeStamp = logEntry.Timestamp
	result.EventType = logEntry.Level
	result.UID = logEntry.InfoParams.DCMessage.Payload.SmcUID

	return &result
}

func processJoinMessage(logEntry parsermodels.ParsedLogEntry) *models.SmcEntry {
	result := models.SmcEntry{}
	result.TimeStamp = logEntry.Timestamp
	result.EventType = logEntry.Level
	result.UID = logEntry.InfoParams.JoinMessage.SmcAddress.SmcUID

	return &result
}

func processStatusMessage(logEntry parsermodels.ParsedLogEntry) *models.StatusEntry {
	result := models.StatusEntry{}
	result.TimeStamp = logEntry.Timestamp
	result.Message = logEntry.InfoParams.StatusMessage.Message
	result.StatusByte = logEntry.InfoParams.StatusMessage.StatusByte

	return &result
}

func processRoutingMessage(logEntry parsermodels.ParsedLogEntry) *models.RoutingEntry {
	result := models.RoutingEntry{}
	result.TimeStamp = logEntry.Timestamp
	result.Address = logEntry.InfoParams.RoutingMessage.Address
	result.NextHopAddress = logEntry.InfoParams.RoutingMessage.NextHopAddress
	result.HopCount = logEntry.InfoParams.RoutingMessage.HopCount
	result.RouteCost = logEntry.InfoParams.RoutingMessage.RouteCost
	result.ValidTimeMins = logEntry.InfoParams.RoutingMessage.ValidTimeMins
	result.WeakLink = logEntry.InfoParams.RoutingMessage.WeakLink

	return &result
}
