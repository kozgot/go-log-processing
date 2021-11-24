package contentparser

import (
	"log"
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type SmcJoinEntryParser struct {
	line models.EntryWithLevelAndTimestamp
}

func (s *SmcJoinEntryParser) Parse() *models.SmcJoinMessageParams {
	smcJoinstring := common.ParseFieldAsString(s.line.Rest, formats.SmcJoinRegex)
	isSmcJoinLine := smcJoinstring != ""
	if !isSmcJoinLine {
		return nil
	}

	smcJoinLine := models.SmcJoinMessageParams{}

	lineRest := strings.Replace(s.line.Rest, smcJoinstring, "", 1)
	// OK [Confirmed] <-- [join_type[LBA] smc_uid[dc18-smc28] physical_address[EEBEDDFFFE6210A5]
	// logical_address[FE80::4021:FF:FE00:000e:61616] short_address[14] last_joining_date[Wed Jun 10 09:37:35 2020]]--(PLC)

	// split the string by the <-- arrow
	// join messages are always directed towards the dc
	messageParts := strings.Split(lineRest, formats.InComingArrow)
	expectedMessagePartsLength := 2
	if len(messageParts) < expectedMessagePartsLength {
		log.Fatalf("There was no direction indicator in Join message: %s", s.line.Rest)
	}

	responseString := messageParts[0]
	response := common.ParseFieldInBracketsAsString(responseString, formats.AnyLettersBetweenBrackets)
	smcJoinLine.Response = response

	status := common.ParseFieldAsString(responseString, formats.JoinStatusResponseRegex)
	smcJoinLine.Ok = status == "OK"

	payloadString := strings.TrimLeft(messageParts[1], " [")
	payloadString = strings.TrimRight(payloadString, "] ")

	smcJoinLine.JoinType = common.ParseFieldInBracketsAsString(payloadString, formats.JoinTypeRegex)
	smcAddress := models.SmcAddressParams{}

	smcAddress.SmcUID = common.ParseFieldInBracketsAsString(payloadString, formats.SmcUIDRegex)
	smcAddress.PhysicalAddress = common.ParseFieldInBracketsAsString(payloadString, formats.PhysicalAddressRegex)
	smcAddress.LogicalAddress = common.ParseFieldInBracketsAsString(payloadString, formats.LogicalAddressRegex)

	smcAddress.ShortAddress = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(payloadString, formats.ShortAddressRegex))

	smcAddress.LastJoiningDate = common.ParseDateTimeField(s.line.Rest, formats.LastJoiningDateRegex)
	smcJoinLine.SmcAddress = smcAddress
	return &smcJoinLine
}
