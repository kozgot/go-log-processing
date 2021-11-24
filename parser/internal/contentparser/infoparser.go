package contentparser

import (
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type InfoParser struct {
	line                           models.EntryWithLevelAndTimestamp
	routingEntryParser             RoutingMessageParser
	smcJoinEntryParser             SmcJoinEntryParser
	statusEntryParser              StatusEntryParser
	messageEntryParser             MessageEntryParser
	connectionAttemptParser        ConnectionAttemptParser
	configUpdateParser             SmcConfigUpdateParser
	connectionReleasedEntryParser  ConnectionReleasedEntryParser
	initConnectionEntryParser      InitConnectionEntryParser
	internalDiagnosticsEntryParser InternalDiagnosticsEntryParser
}

func NewInfoParser(line models.EntryWithLevelAndTimestamp) *InfoParser {
	inforParser := InfoParser{
		line:                           line,
		routingEntryParser:             RoutingMessageParser{line: line},
		smcJoinEntryParser:             SmcJoinEntryParser{line: line},
		statusEntryParser:              StatusEntryParser{line: line},
		messageEntryParser:             MessageEntryParser{line: line},
		connectionAttemptParser:        ConnectionAttemptParser{line: line},
		configUpdateParser:             SmcConfigUpdateParser{line: line},
		connectionReleasedEntryParser:  ConnectionReleasedEntryParser{line: line},
		initConnectionEntryParser:      InitConnectionEntryParser{line: line},
		internalDiagnosticsEntryParser: InternalDiagnosticsEntryParser{line: line},
	}

	return &inforParser
}

func (infoParser *InfoParser) ParseInfo() *models.InfoParams {
	infoParams := models.InfoParams{}

	routingMessage := infoParser.routingEntryParser.Parse()
	if routingMessage != nil {
		infoParams.RoutingMessage = routingMessage
		infoParams.EntryType = models.Routing
		return &infoParams
	}

	joinMessage := infoParser.smcJoinEntryParser.Parse()
	if joinMessage != nil {
		infoParams.JoinMessage = joinMessage
		infoParams.EntryType = models.SMCJoin
		return &infoParams
	}

	statusMessage := infoParser.statusEntryParser.Parse()
	if statusMessage != nil {
		infoParams.StatusMessage = statusMessage
		infoParams.EntryType = models.NetworkStatus
		return &infoParams
	}

	dcMessage := infoParser.messageEntryParser.Parse()
	if dcMessage != nil {
		infoParams.DCMessage = dcMessage
		infoParams.EntryType = models.DCMessage
		return &infoParams
	}

	connectionAttempt := infoParser.connectionAttemptParser.Parse()
	if connectionAttempt != nil {
		infoParams.ConnectionAttempt = connectionAttempt
		infoParams.EntryType = models.ConnectionAttempt
		return &infoParams
	}

	configUpdate := infoParser.configUpdateParser.Parse()
	if configUpdate != nil {
		infoParams.SmcConfigUpdate = configUpdate
		infoParams.EntryType = models.SmcConfigUpdate
		return &infoParams
	}

	connectionReleased := infoParser.connectionReleasedEntryParser.Parse()
	if connectionReleased != nil {
		infoParams.ConnectionReleased = connectionReleased
		infoParams.EntryType = models.ConnectionReleased
		return &infoParams
	}

	initConnectionParams := infoParser.initConnectionEntryParser.Parse()
	if initConnectionParams != nil {
		infoParams.InitConnection = initConnectionParams
		infoParams.EntryType = models.InitDLMSConnection
		return &infoParams
	}

	internalDiagnosticsEntry := infoParser.internalDiagnosticsEntryParser.Parse()
	if internalDiagnosticsEntry != nil {
		infoParams.InternalDiagnosticsData = internalDiagnosticsEntry
		infoParams.EntryType = models.InternalDiagnostics
		return &infoParams
	}

	return &infoParams
}
