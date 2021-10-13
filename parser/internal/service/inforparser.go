package service

import (
	"log"
	"strings"

	"github.com/kozgot/go-log-processing/parser/pkg/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func parseInfo(line models.LineWithDate) *models.InfoParams {
	infoParams := models.InfoParams{}

	routingMessage := parseRoutingTableLine(line.Rest) // done
	if routingMessage != nil {
		infoParams.RoutingMessage = *routingMessage
		infoParams.MessageType = models.RountingMessageType
		return &infoParams
	}

	joinMessage := parseSmcJoinLine(line.Rest) // todo: finish payload parsing
	if joinMessage != nil {
		infoParams.JoinMessage = *joinMessage
		infoParams.MessageType = models.JoinMessageType
		return &infoParams
	}

	statusMessage := parseStatusLine(line.Rest)
	if statusMessage != nil {
		infoParams.StatusMessage = *statusMessage
		infoParams.MessageType = models.StatusMessageType
		return &infoParams
	}

	dcMessage := parseDCMessage(line.Rest)
	if dcMessage != nil {
		infoParams.DCMessage = *dcMessage
		infoParams.MessageType = models.DCMessageType
		return &infoParams
	}

	connectionAttempt := parseConnectionAttempt(line.Rest)
	if connectionAttempt != nil {
		infoParams.ConnectionAttempt = *connectionAttempt
		infoParams.MessageType = models.ConnectionAttempt
		return &infoParams
	}

	if strings.Contains(line.Rest, "Update SMC configuration in DB ") {
		smcAddress := parseSmcAddressPayload(line.Rest)
		if smcAddress != nil {
			smcConfigUpdate := models.SmcConfigUpdateParams{
				PhysicalAddress: smcAddress.PhysicalAddress,
				LogicalAddress:  smcAddress.LogicalAddress,
				ShortAddress:    smcAddress.ShortAddress,
				LastJoiningDate: smcAddress.LastJoiningDate}
			infoParams.SmcConfigUpdate = smcConfigUpdate
			infoParams.MessageType = models.SmcConfigUpdate
			return &infoParams
		}
	}

	return &infoParams
}

func parseStatusLine(line string) *models.StatusMessageParams {
	statusLine := models.StatusMessageParams{}
	statusLine.StatusByte = parseFieldInBracketsAsString(line, formats.StatusByteRegex)
	if statusLine.StatusByte == "" {
		return nil
	}

	statusLine.Message = parseFieldAsString(line, formats.StatusMessageRegex)

	return &statusLine
}

func parseRoutingTableLine(line string) *models.RoutingTableParams {
	isRoutingTableLine := parseFieldAsString(line, formats.RoutingTableRegex) != ""
	if !isRoutingTableLine {
		return nil
	}

	routingtableLine := models.RoutingTableParams{}
	routingtableLine.Address = parseFieldInBracketsAsString(line, formats.RoutingAddressRegex)
	routingtableLine.NextHopAddress = parseFieldInBracketsAsString(line, formats.NextHopAddressRegex)
	routingtableLine.RouteCost = tryParseIntFromString(parseFieldInBracketsAsString(line, formats.RouteCostRegex))
	routingtableLine.ValidTimeMins = tryParseIntFromString(parseFieldInBracketsAsString(line, formats.ValidTimeRegex))
	routingtableLine.WeakLink = tryParseIntFromString(parseFieldInBracketsAsString(line, formats.WeakLinkRegex))
	routingtableLine.HopCount = tryParseIntFromString(parseFieldInBracketsAsString(line, formats.HopCountRegex))

	return &routingtableLine
}

func parseSmcJoinLine(line string) *models.SmcJoinMessageParams {
	smcJoinstring := parseFieldAsString(line, formats.SmcJoinRegex)
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
	response := parseStringBetweenBrackets(responseString)
	smcJoinLine.Response = response

	status := parseJoinStatus(responseString)
	smcJoinLine.Ok = status == "OK" // todo: is there a better way?

	payloadString := strings.TrimLeft(messageParts[1], " [") // todo: is there a better way?
	payloadString = strings.TrimRight(payloadString, "] ")   // todo: is there a better way?

	joinType := parseFieldInBracketsAsString(payloadString, formats.JoinTypeRegex)
	if joinType != "" {
		smcJoinLine.JoinType = joinType
	}

	smcAddress := models.SmcAddressParams{}

	smcUID := parseFieldInBracketsAsString(payloadString, formats.SmcUIDRegex)
	if smcUID != "" {
		smcAddress.SmcUID = smcUID
	}

	physicalAddress := parseFieldInBracketsAsString(payloadString, formats.PhysicalAddressRegex)
	if physicalAddress != "" {
		smcAddress.PhysicalAddress = physicalAddress
	}

	logicalAddress := parseFieldInBracketsAsString(payloadString, formats.LogicalAddressRegex)
	if logicalAddress != "" {
		smcAddress.LogicalAddress = logicalAddress
	}

	shortAddress := parseFieldInBracketsAsString(payloadString, formats.ShortAddressRegex)
	if shortAddress != "" {
		smcAddress.ShortAddress = tryParseIntFromString(shortAddress)
	}

	// todo last joining date

	smcJoinLine.SmcAddress = smcAddress
	return &smcJoinLine
}

func parseStringBetweenBrackets(line string) string {
	return parseFieldInBracketsAsString(line, formats.AnyLettersBetweenBrackets)
}

func parseJoinStatus(line string) string {
	parsed := parseFieldAsString(line, formats.JoinStatusResponseRegex)
	return parsed
}
