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

	// todo: try to parse all INFO lines from both files
	// 3 kinds from plc_manager (routing table, smc join, message)
	// from dc_main: message sent/recieved lines with all kinds of payload (for now, that is enough)
	// --> collect all payload stuff into a note for starters

	// could not parse log level
	return &infoParams
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
	isSmcJoinLine := parseFieldAsString(line, SmcJoinRegex) != ""
	if !isSmcJoinLine {
		return nil
	}

	lineRest := strings.Replace(line, SmcJoinRegex, "", 1)
	// OK [Confirmed] <-- [join_type[LBA] smc_uid[dc18-smc28] physical_address[EEBEDDFFFE6210A5] logical_address[FE80::4021:FF:FE00:000e:61616] short_address[14] last_joining_date[Wed Jun 10 09:37:35 2020]]--(PLC)

	messageParts := strings.Split(lineRest, "<--")
	if len(messageParts) < 2 {
		log.Fatalf("There was no direction indicator in Join message: %s", line)
	}

	responseString := messageParts[0]
	log.Println(responseString)

	payloadString := messageParts[1]
	log.Println("payload: ", payloadString)

	smcJoinLine := SmcJoinMessageParams{}

	return &smcJoinLine
}
