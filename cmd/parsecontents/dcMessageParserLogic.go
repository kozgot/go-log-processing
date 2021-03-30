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

	dcMessageParams.Payload = parseDCMessagePayload(lin, dcMessageParams.MessageType)

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

func parseDCMessagePayload(line string, messageType string) *DcMessagePayload {
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

	switch messageType {
	case "message":
		payload.MessagePayload = parseMessagePayload(line)
		if payload.MessagePayload != nil {
			return &payload
		}

	case "connect":
		payload.ConnectOrDisconnectPayload = parseConnectOrDisconnectPayload(line)
		if payload.ConnectOrDisconnectPayload != nil {
			return &payload
		}

	case "pod configuration":
		payload.PodConfigPayload = parsePodConfigPayload(line)
		if payload.PodConfigPayload != nil {
			return &payload
		}
	case "smc configuration":
		payload.SmcConfigPayload = parseSmcConfigPayload(line)
		if payload.SmcConfigPayload != nil {
			return &payload
		}

	case "smc address":
		payload.SmcAddressPayload = parseSmcAddressPayload(line)
		if payload.SmcAddressPayload != nil {
			return &payload
		}

	case "service_level":
		payload.ServiceLevelPayload = parseServiceLevelPayload(line)
		if payload.ServiceLevelPayload != nil {
			return &payload
		}

	case "settings":
		payload.SettingsPayload = parseSettingsPayload(line)
		if payload.SettingsPayload != nil {
			return &payload
		}

	case "DLMS Logs":
		payload.DLMSLogPayload = parseDLMSLogPayload(line)
		if payload.DLMSLogPayload != nil {
			return &payload
		}

	case "index":
		payload.IndexPayload = parseIndexPayload(line)
		if payload.IndexPayload != nil {
			return &payload
		}
	}

	return &payload
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
	previousValueString := parseFieldInBracketsAsString(line, PreviousValueRegex)
	serialNumberString := parseFieldInBracketsAsString(line, SerailNumberRegex)
	previousTimeFromSeconds := parseTimeFieldFromSeconds(line, PreviousTimeRegex)
	if previousTimeFromSeconds.Year() < 1000 && serialNumberString == "" && previousValueString == "" {
		return nil
	}

	previousValue := tryParseIntFromString(previousValueString)
	serialNumber := tryParseIntFromString(serialNumberString)

	result := IndexPayload{}
	result.PreviousTime = previousTimeFromSeconds
	result.PreviousValue = previousValue
	result.SerialNumber = serialNumber

	return &result
}

func parseMessagePayload(line string) *MessagePayload {
	currentString := parseFieldInBracketsAsString(line, CurrentRegex)
	totalString := parseFieldInBracketsAsString(line, TotalRegex)
	if currentString == "" && totalString == "" {
		return nil
	}

	result := MessagePayload{}

	result.Current = tryParseFloat64FromString(currentString)
	result.Total = tryParseFloat64FromString(parseFieldInBracketsAsString(line, TotalRegex))
	result.URL = parseFieldInBracketsAsString(line, URLRegex)
	result.Topic = parseFieldInBracketsAsString(line, TopicRegex)

	return &result
}

func parseSettingsPayload(line string) *SettingsPayload {
	dcUID := parseFieldInBracketsAsString(line, DcUidRegex)

	indexCollectionString := parseFieldInBracketsAsString(line, IndexCollectionRegex)
	dataPublishString := parseFieldInBracketsAsString(line, DataPublishRegex)
	frequencyBandChangedString := parseFieldInBracketsAsString(line, FrequencyBandChangedRegex)
	frequencyBandRollBackDonestirng := parseFieldInBracketsAsString(line, FrequencyBandRollbackDoneRegex)

	lastServerCommTimeFromSeconds := parseTimeFieldFromSeconds(line, LastServerCommunicationTimeRegex)
	lastDcStartTimeFromSeconds := parseTimeFieldFromSeconds(line, LastDcStartTimeRegex)

	if dcUID == "" && lastServerCommTimeFromSeconds.Year() < 1000 && lastDcStartTimeFromSeconds.Year() < 1000 && dataPublishString == "" && indexCollectionString == "" && frequencyBandChangedString == "" && frequencyBandRollBackDonestirng == "" {
		return nil
	}

	locality := parseFieldInBracketsAsString(line, LocalityRegex)
	region := parseFieldInBracketsAsString(line, RegionRegex)
	timezone := parseFieldInBracketsAsString(line, TimezoneRegex)
	globalFtpAddress := parseFieldInBracketsAsString(line, GlobalFtpAddressRegex)
	targetFirmwareVersion := parseFieldInBracketsAsString(line, TargetFirmwareVersionRegex)
	dcDistroTargetFirmwareVersion := parseFieldInBracketsAsString(line, DcDistroTargetFirmwareVersionRegex)

	indexCollection := tryParseIntFromString(indexCollectionString)
	dataPublish := tryParseIntFromString(dataPublishString)
	frequencyBandChanged := tryParseIntFromString(frequencyBandChangedString) == 1
	frequencyBandRollBackDone := tryParseIntFromString(frequencyBandRollBackDonestirng) == 1

	result := SettingsPayload{DcUID: dcUID, Locality: locality, Region: region, Timezone: timezone, GlobalFtpAddress: globalFtpAddress}

	result.TargetFirmwareVersion = targetFirmwareVersion
	result.DcDistroTargetFirmwareVersion = dcDistroTargetFirmwareVersion
	result.IndexCollection = indexCollection
	result.DataPublish = dataPublish
	result.FrequencyBandChanged = frequencyBandChanged
	result.FrequencyBandRollBackDone = frequencyBandRollBackDone

	return &result
}

func parseServiceLevelPayload(line string) *ServiceLevelPayload {
	meterModeString := parseFieldInBracketsAsString(line, MeterModeRegex)
	maxActivePowerstring := parseFieldInBracketsAsString(line, MaxActivePowerRegex)
	loadSheddingDailyEnergyBudgetString := parseFieldInBracketsAsString(line, LoadSheddingDailyEnergyBudgetRegex)
	localSheddingDailyEnergyBudgetString := parseFieldInBracketsAsString(line, LocalSheddingDailyEnergyBudgetRegex)
	inServiceString := parseFieldInBracketsAsString(line, InServiceRegex)
	if meterModeString == "" && maxActivePowerstring == "" && loadSheddingDailyEnergyBudgetString == "" && localSheddingDailyEnergyBudgetString == "" && inServiceString == "" {
		// if we could not parse any of these fields, than most likely we will not be able to parse the remaining
		return nil
	}

	meterMode := tryParseIntFromString(meterModeString)
	maxActivePower := tryParseIntFromString(maxActivePowerstring)
	loadSheddingDailyEnergyBudget := tryParseIntFromString(loadSheddingDailyEnergyBudgetString)
	localSheddingDailyEnergyBudget := tryParseIntFromString(localSheddingDailyEnergyBudgetString)
	inService := tryParseIntFromString(inServiceString) == 1

	startHourDailyCycle := parseFieldInBracketsAsString(line, StartHourDailyCycleRegex)
	name := parseFieldInBracketsAsString(line, NameRegex)

	hourlyEnergyLimits := parseHourlyEnergyLimits(line, HourlyEnergyLimitsRegex)
	localHourlyEnergyLimits := parseHourlyEnergyLimits(line, LocalHourlyEnergyLimitsRegex)

	result := ServiceLevelPayload{MeterMode: meterMode, MaxActivePower: maxActivePower, LoadSheddingDailyEnergyBudget: loadSheddingDailyEnergyBudget, LocalSheddingDailyEnergyBudget: localSheddingDailyEnergyBudget, InService: inService, StartHourDailyCycle: startHourDailyCycle, Name: name}
	result.HourlyEnergyLimits = hourlyEnergyLimits
	result.LocalHourlyEnergyLimits = localHourlyEnergyLimits
	return &result
}

func parseHourlyEnergyLimits(line string, energyLimitRegex string) [24]HourlyEnergyLimit {
	var result [24]HourlyEnergyLimit
	hourlyLimitsString := parseFieldInDoubleBracketsAsString(line, energyLimitRegex)
	if hourlyLimitsString != "" {
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
	smcUID := parseFieldInBracketsAsString(line, SMCUIDRegex)
	physicalAddress := parseFieldInBracketsAsString(line, PhysicalAddressRegex)
	logicalAddress := parseFieldInBracketsAsString(line, LogicalAddressRegex)
	shortAddressString := parseFieldInBracketsAsString(line, ShortAddressRegex)
	if smcUID == "" && physicalAddress == "" && logicalAddress == "" && shortAddressString == "" {
		return nil
	}

	lastJoiningDate := parseDateTimeField(line, LastJoiningDateRegex)
	shortAddress := tryParseIntFromString(shortAddressString)

	result := SmcAddressParams{SmcUID: smcUID, PhysicalAddress: physicalAddress, LogicalAddress: logicalAddress, ShortAddress: shortAddress, LastJoiningDate: lastJoiningDate}
	return &result
}

func parseSmcConfigPayload(line string) *SmcConfigPayload {
	customerSerialNumber := parseFieldInBracketsAsString(line, CustomerSerialNumberRegex)
	physicalAddress := parseFieldInBracketsAsString(line, PhysicalAddressRegex)
	smcStatus := parseFieldInBracketsAsString(line, SmcStatusRegex)
	nextHopString := parseFieldInBracketsAsString(line, NextHopRegex)

	if customerSerialNumber == "" && physicalAddress == "" && smcStatus == "" && nextHopString == "" {
		return nil
	}

	currentApp1Fw := parseFieldInBracketsAsString(line, CurrentApp1FwRegex)
	currentApp2Fw := parseFieldInBracketsAsString(line, CurrentApp2FwRegex)
	currentPlcFw := parseFieldInBracketsAsString(line, CurrentPlcFwRegex)

	lastSuccessfulDlmsResponseDate := parseDateTimeField(line, LastSuccessfulDlmsResponseDateRegex)
	nextHop := tryParseIntFromString(nextHopString)

	result := SmcConfigPayload{}
	result.CurrentApp1Fw = currentApp1Fw
	result.CurrentApp2Fw = currentApp2Fw
	result.CurrentPlcFw = currentPlcFw
	result.CustomerSerialNumber = customerSerialNumber
	result.PhysicalAddress = physicalAddress
	result.SmcStatus = smcStatus
	result.NextHop = nextHop
	result.LastSuccessfulDlmsResponseDate = lastSuccessfulDlmsResponseDate

	return &result
}

func parsePodConfigPayload(line string) *PodConfigPayload {
	serialNumberString := parseFieldInBracketsAsString(line, SerialNumberRegex)
	phaseString := parseFieldInBracketsAsString(line, PhaseRegex)
	positionInSmcString := parseFieldInBracketsAsString(line, PositionInSmcRegex)
	softwareFirmwareVersion := parseFieldInBracketsAsString(line, SoftwareFirmwareVersionRegex)
	if serialNumberString == "" && phaseString == "" && positionInSmcString == "" && softwareFirmwareVersion == "" {
		return nil
	}

	serialNumber := tryParseIntFromString(serialNumberString)
	phase := tryParseIntFromString(phaseString)
	positionInSmc := tryParseIntFromString(positionInSmcString)

	result := PodConfigPayload{Phase: phase, SerialNumber: serialNumber, PositionInSmc: positionInSmc, SoftwareFirmwareVersion: softwareFirmwareVersion}
	return &result
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
