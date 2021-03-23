package parsecontents

import (
	"log"
	"strings"
)

func parseDCMessage(lin string) *DCMessageParams {
	dcMessageParams := DCMessageParams{}

	source := parseDCMessageSource(lin)
	dest := parseDCMessageDest(lin)
	if source != "" {
		dcMessageParams.IsInComing = true
		dcMessageParams.SourceOrDestName = source
		dcMessageParams.MessageType = parseFieldInBracketsAsString(lin, IncomingMessageTypeRegex)
	} else if dest != "" {
		dcMessageParams.IsInComing = false
		dcMessageParams.SourceOrDestName = dest
		dcMessageParams.MessageType = parseFieldInBracketsAsString(lin, OutGoingMessageTypeRegex)
	}

	dcMessageParams.Payload = parseDCMessagePayload(lin)

	return &dcMessageParams
}

func parseDCMessageSource(line string) string {
	inComingMessageSource := parseFieldInParenthesesAsString(line, IncomingMessageSourceRegex)
	return inComingMessageSource
}

func parseDCMessageDest(line string) string {
	outGoingMessageSource := parseFieldInParenthesesAsString(line, OutGoingMessageDestRegex)
	return outGoingMessageSource
}

func parseDCMessagePayload(line string) DcMessagePayload {
	payload := DcMessagePayload{}
	payload.SmcUID = parseFieldInBracketsAsString(line, SmcUidRegex)
	payload.PodUID = parseFieldInBracketsAsString(line, PodUidRegex)
	payload.ServiceLevelId = tryParseIntFromString(parseFieldInBracketsAsString(line, ServiceLevelIdRegex))
	payload.Value = tryParseIntFromString(parseFieldInBracketsAsString(line, ValueRegex))
	if strings.Contains(line, "--[DLMS Logs]-->(SVI) ") {
		log.Print("dlms ")
	}
	dateTime := parseDateTimeField(line, DateTimeRegexFieldRegex)
	ticksTime := tryParseIntFromString(parseFieldInBracketsAsString(line, TimeTicksRegex))
	if ticksTime != 0 {
		log.Print(ticksTime)
		// todo convert to datetime (human readable format)
		// for some reason, it can parse the long date format to int, so that needs to be handled as well
	}
	if dateTime.Year() > 1000 {
		log.Print(dateTime)
	}
	// todo rest

	return payload
}

/*
type DcMessagePayload struct {
	Time           time.Time

	TimeRange TimeRange

	ConnectOrDisconnectPayload ConnectOrDisconnectPayload
	DLMSLogPayload             DLMSLogPayload
	IndexPayload               IndexPayload
	MessagePayload             MessagePayload
	SettingsPayload            SettingsPayload
	ServiceLevelPayload        ServiceLevelPayload
	SmcAddressPayload          SmcAddressParams
	SmcConfigPayload           SmcConfigPayload
	PodConfigPayload           PodConfigPayload
}
*/
