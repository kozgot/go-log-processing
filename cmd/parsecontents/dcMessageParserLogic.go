package parsecontents

import (
	"log"
	"strings"
	"time"
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
	} else {
		return nil
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

func parseDCMessagePayload(line string) *DcMessagePayload {
	payload := DcMessagePayload{}
	payload.SmcUID = parseFieldInBracketsAsString(line, SmcUidRegex)
	payload.PodUID = parseFieldInBracketsAsString(line, PodUidRegex)
	payload.ServiceLevelId = tryParseIntFromString(parseFieldInBracketsAsString(line, ServiceLevelIdRegex))
	payload.Value = tryParseIntFromString(parseFieldInBracketsAsString(line, ValueRegex))

	// Parse the time[] field of the message. It can be a formatted date or in a date represented by a timestamp in seconds.
	dateTime := parseDateTimeField(line, DateTimeFieldRegex)
	if dateTime.Year() > 1000 {
		payload.Time = dateTime // for some reason, it can parse the long date format to int, so that needs to be handled as well (hence the if-else)
	} else {
		datefromSeconds := parseTimeFieldFromSeconds(line, TimeTicksRegex)
		if datefromSeconds.Year() > 1000 {
			payload.Time = datefromSeconds
		}
	}

	payload.TimeRange = parseTimeRange(line)

	payload.ConnectOrDisconnectPayload = parseConnectOrDisconnectPayload(line)
	if payload.ConnectOrDisconnectPayload != nil {
		return &payload
	}

	payload.DLMSLogPayload = parseDLMSLogPayload(line)
	if payload.DLMSLogPayload != nil {
		return &payload
	}
	payload.IndexPayload = parseIndexPayload(line)
	if payload.IndexPayload != nil {
		return &payload
	}
	payload.MessagePayload = parseMessagePayload(line)
	if payload.MessagePayload != nil {
		return &payload
	}
	payload.PodConfigPayload = parsePodConfigPayload(line)
	if payload.PodConfigPayload != nil {
		return &payload
	}
	payload.SmcConfigPayload = parseSmcConfigPayload(line)
	if payload.SmcConfigPayload != nil {
		return &payload
	}
	payload.SmcAddressPayload = parseSmcAddressPayload(line)
	if payload.SmcAddressPayload != nil {
		return &payload
	}
	payload.SettingsPayload = parseSettingsPayload(line)
	if payload.SettingsPayload != nil {
		return &payload
	}
	payload.ServiceLevelPayload = parseServiceLevelPayload(line)
	if payload.ServiceLevelPayload != nil {
		return &payload
	}

	return &payload
}

func parseTimeRange(line string) *TimeRange {
	from := parseDateTimeField(line, TimeRangeFromRegex)
	if from.Year() < 1500 {
		from = parseTimeFieldFromSeconds(line, TimeRangeStartTicksRegex)
	}

	to := parseDateTimeField(line, TimeRangeToRegex)
	if to.Year() < 1500 {
		to = parseTimeFieldFromSeconds(line, TimeRangeEndTicksRegex)
	}

	if from.Year() > 1500 && to.Year() > 1500 {
		result := TimeRange{From: from, To: to}
		return &result
	}

	return nil
}

func parseTimeFieldFromSeconds(line string, timeStampRegex string) time.Time {
	seconds := tryParseInt64FromString(parseFieldInBracketsAsString(line, timeStampRegex))
	if seconds != 0 {
		dateTimeFromsSecs := time.Unix(seconds, 0)
		return dateTimeFromsSecs
	}

	return time.Time{}
}

func parseTimeFieldFromMilliSeconds(line string, timeStampRegex string) time.Time {
	milliseconds := tryParseInt64FromString(parseFieldInBracketsAsString(line, timeStampRegex))
	if milliseconds != 0 {
		dateTimeFromsSecs := time.Unix(0, milliseconds*1000*1000)
		return dateTimeFromsSecs
	}

	return time.Time{}
}

func parseConnectOrDisconnectPayload(line string) *ConnectOrDisconnectPayload {
	resultType := tryParseIntFromString(parseFieldInBracketsAsString(line, ConnectOrDisconnectTypeRegex))
	clientId := parseFieldInBracketsAsString(line, ClientIdRegex)
	URL := parseFieldInBracketsAsString(line, URLRegex)
	topic := parseFieldInBracketsAsString(line, TopicRegex)
	timeout := tryParseIntFromString(parseFieldInBracketsAsString(line, TimeoutRegex))
	connected := tryParseIntFromString(parseFieldInBracketsAsString(line, ConnectedRegex)) == 1

	if clientId != "" || resultType != 0 || URL != "" || topic != "" || timeout != 0 {
		result := ConnectOrDisconnectPayload{ClientId: clientId, Type: resultType, URL: URL, Topic: topic, Timeout: timeout, Connected: connected}
		return &result
	}

	// todo: ellenőrzés
	return nil
}

func parseDLMSLogPayload(line string) *DLMSLogPayload {
	requestTimeFromSeconds := parseTimeFieldFromMilliSeconds(line, DLMSRequestTimeRegex)
	responseTimeFromSeconds := parseTimeFieldFromMilliSeconds(line, DLMSResponseTimeRegex)
	DLMSError := parseFieldInBracketsAsString(line, DLMSErrorRegex)

	if requestTimeFromSeconds.Year() > 1500 || responseTimeFromSeconds.Year() > 1500 || DLMSError != "" {
		result := DLMSLogPayload{DLMSRequestTime: requestTimeFromSeconds, DLMSResponseTime: responseTimeFromSeconds, DLMSError: DLMSError}
		return &result
	}
	// todo: ellenőrzés
	return nil
}

func parseIndexPayload(line string) *IndexPayload {
	result := IndexPayload{}
	previousTimeFromSeconds := parseTimeFieldFromSeconds(line, PreviousTimeRegex)
	if previousTimeFromSeconds.Year() > 1000 {
		result.PreviousTime = previousTimeFromSeconds
	}

	result.PreviousValue = tryParseIntFromString(parseFieldInBracketsAsString(line, PreviousValueRegex))
	result.SerialNumber = tryParseIntFromString(parseFieldInBracketsAsString(line, SerailNumberRegex))

	if result.PreviousValue != 0 || result.PreviousTime.Year() > 1500 || result.SerialNumber != 0 {
		return &result
	}

	// todo: ellenőrzés
	return nil
}

func parseMessagePayload(line string) *MessagePayload {
	result := MessagePayload{}
	result.Current = tryParseFloat64FromString(parseFieldInBracketsAsString(line, CurrentRegex))
	result.Total = tryParseFloat64FromString(parseFieldInBracketsAsString(line, TotalRegex))
	result.URL = parseFieldInBracketsAsString(line, URLRegex)
	result.Topic = parseFieldInBracketsAsString(line, TopicRegex)

	if result.Current != 0 || result.URL != "" || result.Total != 0 || result.Topic != "" {
		return &result
	}

	// todo: ellenőrzés
	return nil
}

func parseSettingsPayload(line string) *SettingsPayload {
	result := SettingsPayload{}
	result.DcUID = parseFieldInBracketsAsString(line, DcUidRegex)
	result.Locality = parseFieldInBracketsAsString(line, LocalityRegex)
	result.Region = parseFieldInBracketsAsString(line, RegionRegex)
	result.Timezone = parseFieldInBracketsAsString(line, TimezoneRegex)
	result.GlobalFtpAddress = parseFieldInBracketsAsString(line, GlobalFtpAddressRegex)
	result.TargetFirmwareVersion = parseFieldInBracketsAsString(line, TargetFirmwareVersionRegex)

	result.IndexCollection = tryParseIntFromString(parseFieldInBracketsAsString(line, IndexCollectionRegex))
	result.DataPublish = tryParseIntFromString(parseFieldInBracketsAsString(line, DataPublishRegex))
	lastServerCommTimeFromSeconds := parseTimeFieldFromSeconds(line, LastServerCommunicationTimeRegex)
	if lastServerCommTimeFromSeconds.Year() > 1000 {
		result.LastServerCommunicationTime = lastServerCommTimeFromSeconds
	}

	result.DcDistroTargetFirmwareVersion = parseFieldInBracketsAsString(line, DcDistroTargetFirmwareVersionRegex)

	lastDcStartTimeFromSeconds := parseTimeFieldFromSeconds(line, LastDcStartTimeRegex)
	if lastDcStartTimeFromSeconds.Year() > 1000 {
		result.LastDcStartTime = lastDcStartTimeFromSeconds
	}

	result.FrequencyBandChanged = tryParseIntFromString(parseFieldInBracketsAsString(line, FrequencyBandChangedRegex)) == 1
	result.FrequencyBandRollBackDone = tryParseIntFromString(parseFieldInBracketsAsString(line, FrequencyBandRollbackDoneRegex)) == 1

	if result.IndexCollection != 0 || result.LastServerCommunicationTime.Year() > 1500 || result.LastDcStartTime.Year() > 1500 || result.DataPublish != 0 || result.Timezone != "" { // todo
		return &result
	}

	// todo: ellenőrzés
	return nil
}

func parseServiceLevelPayload(line string) *ServiceLevelPayload {
	result := ServiceLevelPayload{}
	result.MeterMode = tryParseIntFromString(parseFieldInBracketsAsString(line, MeterModeRegex))

	result.StartHourDailyCycle = parseFieldInBracketsAsString(line, StartHourDailyCycleRegex)

	result.LoadSheddingDailyEnergyBudget = tryParseIntFromString(parseFieldInBracketsAsString(line, LoadSheddingDailyEnergyBudgetRegex))
	result.LocalSheddingDailyEnergyBudget = tryParseIntFromString(parseFieldInBracketsAsString(line, LocalSheddingDailyEnergyBudgetRegex))
	result.MaxActivePower = tryParseIntFromString(parseFieldInBracketsAsString(line, MaxActivePowerRegex))

	result.InService = tryParseIntFromString(parseFieldInBracketsAsString(line, InServiceRegex)) == 1

	result.Name = parseFieldInBracketsAsString(line, NameRegex)

	//result.HourlyEnergyLimits = parseHourlyEnergyLimits(line, HourlyEnergyLimitsRegex)
	//result.LocalHourlyEnergyLimits = parseHourlyEnergyLimits(line, LocalHourlyEnergyLimitsRegex)

	if result.MeterMode != 0 || result.LoadSheddingDailyEnergyBudget != 0 || result.LocalSheddingDailyEnergyBudget != 0 || result.Name != "" { // todo
		return &result
	}

	// todo: ellenőrzés
	return nil
}

func parseHourlyEnergyLimits(line string, energyLimitRegex string) [24]HourlyEnergyLimit {
	var result [24]HourlyEnergyLimit
	hourlyLimitsString := parseFieldInBracketsAsString(line, energyLimitRegex)
	if hourlyLimitsString != "" {
		hourlyLimitsString = strings.TrimLeft(hourlyLimitsString, "[")
		hourlyLimitsString = strings.TrimRight(hourlyLimitsString, "]")
		limitParts := strings.Split(hourlyLimitsString, " ")
		if len(limitParts) == 24 {
			for i, val := range limitParts {
				result[i] = HourlyEnergyLimit{HourNumber: i, Limit: tryParseIntFromString(val)}
			}

			return result
		} else {
			log.Println(limitParts)
		}
	}

	return result
}

func parseSmcAddressPayload(line string) *SmcAddressParams {
	result := SmcAddressParams{}

	// todo
	/*
			type SmcAddressParams struct {
			SmcUID          string
			PhysicalAddress string
			LogicalAddress  string
			ShortAddress    int
			LastJoiningDate time.Time
		}
	*/

	if result.SmcUID != "" {
		return &result
	}

	return nil
}

func parseSmcConfigPayload(line string) *SmcConfigPayload {
	result := SmcConfigPayload{}

	// todo
	/*
			type SmcConfigPayload struct {
			CustomerSerialNumber           string
			PhysicalAddress                string
			SmcStatus                      string
			CurrentApp1Fw                  string
			CurrentApp2Fw                  string
			CurrentPlcFw                   string
			LastSuccessfulDlmsResponseDate time.Time
			NextHop                        int
		}
	*/

	if result.CustomerSerialNumber != "" {
		return &result
	}

	return nil
}

func parsePodConfigPayload(line string) *PodConfigPayload {
	result := PodConfigPayload{}

	// todo
	/*
			type PodConfigPayload struct {
			SerialNumber            int
			Phase                   int
			PositionInSmc           int
			SoftwareFirmwareVersion string
		}
	*/

	if result.SerialNumber != 0 {
		return &result
	}

	return nil
}
