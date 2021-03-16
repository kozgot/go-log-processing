package parsecontents

import (
	"log"
	"strings"

	"github.com/kozgot/go-log-processing/cmd/parsedates"
)

func parseInfo(line parsedates.LineWithDate) *InfoParams {
	infoParams := InfoParams{}

	routingMessage := parseRoutingTableLine(line.Rest)
	if routingMessage != nil {
		infoParams.RoutingMessage = *routingMessage
		infoParams.MessageType = RountingMessageType
		return &infoParams
	}

	joinMessage := parseSmcJoinLine(line.Rest)
	if joinMessage != nil {
		infoParams.JoinMessage = *joinMessage
		infoParams.MessageType = JoinMessageType
		return &infoParams
	}

	statusMessage := parseStatusLine(line.Rest)
	if statusMessage != nil {
		infoParams.StatusMessage = *statusMessage
		infoParams.MessageType = StatusMessageType
		return &infoParams
	}

	dcMessage := parseDCMessage(line.Rest)
	if dcMessage != nil {
		infoParams.DCMessage = *dcMessage
		infoParams.MessageType = DCMessageType
		return &infoParams
	}

	return &infoParams
}

func parseStatusLine(line string) *StatusMessageParams {
	statusLine := StatusMessageParams{}
	statusLine.StatusByte = parseFieldInBracketsAsString(line, StatusByteRegex)
	if statusLine.StatusByte == "" {
		return nil
	}

	statusLine.Message = parseFieldAsString(line, StatusMessageRegex)

	return &statusLine
}

func parseRoutingTableLine(line string) *RoutingTableParams {
	isRoutingTableLine := parseFieldAsString(line, RoutingTableRegex) != ""
	if !isRoutingTableLine {
		return nil
	}

	routingtableLine := RoutingTableParams{}
	routingtableLine.Address = parseFieldInBracketsAsString(line, RoutingAddressRegex)
	routingtableLine.NextHopAddress = parseFieldInBracketsAsString(line, NextHopAddressRegex)
	routingtableLine.RouteCost = tryParseIntFromString(parseFieldInBracketsAsString(line, RouteCostRegex))
	routingtableLine.ValidTimeMins = tryParseIntFromString(parseFieldInBracketsAsString(line, ValidTimeRegex))
	routingtableLine.WeakLink = tryParseIntFromString(parseFieldInBracketsAsString(line, WeakLinkRegex))
	routingtableLine.HopCount = tryParseIntFromString(parseFieldInBracketsAsString(line, HopCountRegex))

	return &routingtableLine
}

func parseSmcJoinLine(line string) *SmcJoinMessageParams {
	smcJoinstring := parseFieldAsString(line, SmcJoinRegex)
	isSmcJoinLine := smcJoinstring != ""
	if !isSmcJoinLine {
		return nil
	}

	smcJoinLine := SmcJoinMessageParams{}

	lineRest := strings.Replace(line, smcJoinstring, "", 1)
	// OK [Confirmed] <-- [join_type[LBA] smc_uid[dc18-smc28] physical_address[EEBEDDFFFE6210A5] logical_address[FE80::4021:FF:FE00:000e:61616] short_address[14] last_joining_date[Wed Jun 10 09:37:35 2020]]--(PLC)

	// todo: rest of the implementation
	messageParts := strings.Split(lineRest, "<--") // todo: match indicator with regex, or at least constant
	if len(messageParts) < 2 {
		log.Fatalf("There was no direction indicator in Join message: %s", line)
	}

	responseString := messageParts[0]
	response := parseStringBetweenBrackets(responseString)
	smcJoinLine.Response = response

	status := parseJoinStatus(responseString)
	smcJoinLine.Ok = status == "OK" // todo: is there a better way?

	payloadString := strings.TrimLeft(messageParts[1], " [") // todo: is there a better way?
	payloadString = strings.TrimRight(payloadString, "] ")   // todo: is there a better way?
	//log.Println("payload:", payloadString)

	return &smcJoinLine
}

func parseStringBetweenBrackets(line string) string {
	return parseFieldInBracketsAsString(line, AnyLettersBetweenBrackets)
}

func parseJoinStatus(line string) string {
	parsed := parseFieldAsString(line, JoinStatusResponseRegex)
	return parsed
}
