package contentparser

import (
	"log"
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func ParseInfo(line models.EntryWithLevelAndTimestamp) *models.InfoParams {
	infoParams := models.InfoParams{}

	routingMessage := parseRoutingTableLine(line.Rest)
	if routingMessage != nil {
		infoParams.RoutingMessage = routingMessage
		infoParams.EntryType = models.Routing
		return &infoParams
	}

	joinMessage := parseSmcJoinLine(line.Rest)
	if joinMessage != nil {
		infoParams.JoinMessage = joinMessage
		infoParams.EntryType = models.SMCJoin
		return &infoParams
	}

	statusMessage := parseStatusLine(line.Rest)
	if statusMessage != nil {
		infoParams.StatusMessage = statusMessage
		infoParams.EntryType = models.NetworkStatus
		return &infoParams
	}

	dcMessage := parseDCMessage(line.Rest)
	if dcMessage != nil {
		infoParams.DCMessage = dcMessage
		infoParams.EntryType = models.DCMessage
		return &infoParams
	}

	connectionAttempt := parseConnectionAttempt(line.Rest)
	if connectionAttempt != nil {
		infoParams.ConnectionAttempt = connectionAttempt
		infoParams.EntryType = models.ConnectionAttempt
		return &infoParams
	}

	configUpdate := parseSmcConfigUpdate(line.Rest)
	if configUpdate != nil {
		infoParams.SmcConfigUpdate = configUpdate
		infoParams.EntryType = models.SmcConfigUpdate
		return &infoParams
	}

	connectionReleased := parseConnectionReleasedEntry(line.Rest)
	if connectionReleased != nil {
		infoParams.ConnectionReleased = connectionReleased
		infoParams.EntryType = models.ConnectionReleased
		return &infoParams
	}

	initConnectionParams := parseInitConnectionLogEntry(line.Rest)
	if initConnectionParams != nil {
		infoParams.InitConnection = initConnectionParams
		infoParams.EntryType = models.InitDLMSConnection
		return &infoParams
	}

	internalDiagnosticsEntry := parseSmcInternalDiagnosticsEntry(line.Rest)
	if internalDiagnosticsEntry != nil {
		infoParams.InternalDiagnosticsData = internalDiagnosticsEntry
		infoParams.EntryType = models.InternalDiagnostics
		return &infoParams
	}

	return &infoParams
}

func parseStatusLine(line string) *models.StatusMessageParams {
	statusLine := models.StatusMessageParams{}
	statusLine.StatusByte = common.ParseFieldInBracketsAsString(line, formats.StatusByteRegex)
	if statusLine.StatusByte == "" {
		return nil
	}

	statusLine.Message = common.ParseFieldAsString(line, formats.StatusMessageRegex)

	return &statusLine
}

func parseRoutingTableLine(line string) *models.RoutingTableParams {
	isRoutingTableLine := common.ParseFieldAsString(line, formats.RoutingTableRegex) != ""
	if !isRoutingTableLine {
		return nil
	}

	routingtableLine := models.RoutingTableParams{}
	routingtableLine.Address = common.ParseFieldInBracketsAsString(line, formats.RoutingAddressRegex)
	routingtableLine.NextHopAddress = common.ParseFieldInBracketsAsString(line, formats.NextHopAddressRegex)
	routingtableLine.RouteCost = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(line, formats.RouteCostRegex))
	routingtableLine.ValidTimeMins = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(line, formats.ValidTimeRegex))
	routingtableLine.WeakLink = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(line, formats.WeakLinkRegex))
	routingtableLine.HopCount = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(line, formats.HopCountRegex))

	return &routingtableLine
}

func parseSmcJoinLine(line string) *models.SmcJoinMessageParams {
	smcJoinstring := common.ParseFieldAsString(line, formats.SmcJoinRegex)
	isSmcJoinLine := smcJoinstring != ""
	if !isSmcJoinLine {
		return nil
	}

	smcJoinLine := models.SmcJoinMessageParams{}

	lineRest := strings.Replace(line, smcJoinstring, "", 1)
	// OK [Confirmed] <-- [join_type[LBA] smc_uid[dc18-smc28] physical_address[EEBEDDFFFE6210A5]
	// logical_address[FE80::4021:FF:FE00:000e:61616] short_address[14] last_joining_date[Wed Jun 10 09:37:35 2020]]--(PLC)

	// split the string by the <-- arrow
	// join messages are always directed towards the dc
	messageParts := strings.Split(lineRest, formats.InComingArrow)
	expectedMessagePartsLength := 2
	if len(messageParts) < expectedMessagePartsLength {
		log.Fatalf("There was no direction indicator in Join message: %s", line)
	}

	responseString := messageParts[0]
	response := common.ParseFieldInBracketsAsString(responseString, formats.AnyLettersBetweenBrackets)
	smcJoinLine.Response = response

	status := parseJoinStatus(responseString)
	smcJoinLine.Ok = status == "OK"

	payloadString := strings.TrimLeft(messageParts[1], " [")
	payloadString = strings.TrimRight(payloadString, "] ")

	joinType := common.ParseFieldInBracketsAsString(payloadString, formats.JoinTypeRegex)
	if joinType != "" {
		smcJoinLine.JoinType = joinType
	}

	smcAddress := models.SmcAddressParams{}

	smcUID := common.ParseFieldInBracketsAsString(payloadString, formats.SmcUIDRegex)
	if smcUID != "" {
		smcAddress.SmcUID = smcUID
	}

	physicalAddress := common.ParseFieldInBracketsAsString(payloadString, formats.PhysicalAddressRegex)
	if physicalAddress != "" {
		smcAddress.PhysicalAddress = physicalAddress
	}

	logicalAddress := common.ParseFieldInBracketsAsString(payloadString, formats.LogicalAddressRegex)
	if logicalAddress != "" {
		smcAddress.LogicalAddress = logicalAddress
	}

	shortAddress := common.ParseFieldInBracketsAsString(payloadString, formats.ShortAddressRegex)
	if shortAddress != "" {
		smcAddress.ShortAddress = common.TryParseIntFromString(shortAddress)
	}

	smcAddress.LastJoiningDate = common.ParseDateTimeField(line, formats.LastJoiningDateRegex)
	smcJoinLine.SmcAddress = smcAddress
	return &smcJoinLine
}

func parseJoinStatus(line string) string {
	parsed := common.ParseFieldAsString(line, formats.JoinStatusResponseRegex)
	return parsed
}
