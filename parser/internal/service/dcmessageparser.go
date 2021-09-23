package service

import (
	"log"
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/parser/pkg/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func parseDCMessage(lin string) *models.DCMessageParams {
	dcMessageParams := models.DCMessageParams{}

	source := parseDCMessageSource(lin)
	dest := parseDCMessageDest(lin)
	if source != "" {
		dcMessageParams.IsInComing = true
		dcMessageParams.SourceOrDestName = source
		dcMessageParams.MessageType = parseFieldInBracketsAsString(lin, formats.IncomingMessageTypeRegex)
	} else if dest != "" {
		dcMessageParams.IsInComing = false
		dcMessageParams.SourceOrDestName = dest
		dcMessageParams.MessageType = parseFieldInBracketsAsString(lin, formats.OutGoingMessageTypeRegex)
	} else {
		return nil
	}

	dcMessageParams.Payload = parseDCMessagePayload(lin, dcMessageParams.MessageType)

	return &dcMessageParams
}

func parseDCMessageSource(line string) string {
	inComingMessageSource := parseFieldInParenthesesAsString(line, formats.IncomingMessageSourceRegex)
	return inComingMessageSource
}

func parseDCMessageDest(line string) string {
	outGoingMessageSource := parseFieldInParenthesesAsString(line, formats.OutGoingMessageDestRegex)
	return outGoingMessageSource
}

func parseDCMessagePayload(line string, messageType string) *models.DcMessagePayload {
	payload := models.DcMessagePayload{}
	payload.SmcUID = parseFieldInBracketsAsString(line, formats.SmcUidRegex)
	payload.PodUID = parseFieldInBracketsAsString(line, formats.PodUidRegex)
	payload.ServiceLevelId = tryParseIntFromString(parseFieldInBracketsAsString(line, formats.ServiceLevelIdRegex))
	payload.Value = tryParseIntFromString(parseFieldInBracketsAsString(line, formats.ValueRegex))

	// Parse the time[] field of the message.
	// It can be a formatted date or in a date represented by a timestamp in seconds.
	dateTime := parseDateTimeField(line, formats.DateTimeFieldRegex)
	if dateTime.Year() > 1000 {
		// for some reason, it can parse the long date format to int, so that needs to be handled as well (hence the if-else)
		payload.Time = dateTime
	} else {
		datefromSeconds := parseTimeFieldFromSeconds(line, formats.TimeTicksRegex)
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

func parseConnectOrDisconnectPayload(line string) *models.ConnectOrDisconnectPayload {
	resultType := tryParseIntFromString(parseFieldInBracketsAsString(line, formats.ConnectOrDisconnectTypeRegex))
	clientID := parseFieldInBracketsAsString(line, formats.ClientIdRegex)
	URL := parseFieldInBracketsAsString(line, formats.URLRegex)
	topic := parseFieldInBracketsAsString(line, formats.TopicRegex)
	timeout := tryParseIntFromString(parseFieldInBracketsAsString(line, formats.TimeoutRegex))
	connected := tryParseIntFromString(parseFieldInBracketsAsString(line, formats.ConnectedRegex)) == 1

	if clientID != "" || resultType != 0 || URL != "" || topic != "" || timeout != 0 {
		result := models.ConnectOrDisconnectPayload{
			ClientId:  clientID,
			Type:      resultType,
			URL:       URL,
			Topic:     topic,
			Timeout:   timeout,
			Connected: connected}
		return &result
	}

	// todo: ellenőrzés
	return nil
}

func parseDLMSLogPayload(line string) *models.DLMSLogPayload {
	requestTimeFromSeconds := parseTimeFieldFromMilliSeconds(line, formats.DLMSRequestTimeRegex)
	responseTimeFromSeconds := parseTimeFieldFromMilliSeconds(line, formats.DLMSResponseTimeRegex)
	DLMSError := parseFieldInBracketsAsString(line, formats.DLMSErrorRegex)

	if requestTimeFromSeconds.Year() > 1500 || responseTimeFromSeconds.Year() > 1500 || DLMSError != "" {
		result := models.DLMSLogPayload{
			DLMSRequestTime:  requestTimeFromSeconds,
			DLMSResponseTime: responseTimeFromSeconds,
			DLMSError:        DLMSError}
		return &result
	}

	return nil
}

func parseIndexPayload(line string) *models.IndexPayload {
	previousValueString := parseFieldInBracketsAsString(line, formats.PreviousValueRegex)
	serialNumberString := parseFieldInBracketsAsString(line, formats.SerailNumberRegex)
	previousTimeFromSeconds := parseTimeFieldFromSeconds(line, formats.PreviousTimeRegex)
	if previousTimeFromSeconds.Year() < 1000 && serialNumberString == "" && previousValueString == "" {
		return nil
	}

	previousValue := tryParseIntFromString(previousValueString)
	serialNumber := tryParseIntFromString(serialNumberString)

	result := models.IndexPayload{}
	result.PreviousTime = previousTimeFromSeconds
	result.PreviousValue = previousValue
	result.SerialNumber = serialNumber

	return &result
}

func parseMessagePayload(line string) *models.MessagePayload {
	currentString := parseFieldInBracketsAsString(line, formats.CurrentRegex)
	totalString := parseFieldInBracketsAsString(line, formats.TotalRegex)
	if currentString == "" && totalString == "" {
		return nil
	}

	result := models.MessagePayload{}

	result.Current = tryParseFloat64FromString(currentString)
	result.Total = tryParseFloat64FromString(parseFieldInBracketsAsString(line, formats.TotalRegex))
	result.URL = parseFieldInBracketsAsString(line, formats.URLRegex)
	result.Topic = parseFieldInBracketsAsString(line, formats.TopicRegex)

	return &result
}

func parseSettingsPayload(line string) *models.SettingsPayload {
	dcUID := parseFieldInBracketsAsString(line, formats.DcUidRegex)

	indexCollectionString := parseFieldInBracketsAsString(line, formats.IndexCollectionRegex)
	dataPublishString := parseFieldInBracketsAsString(line, formats.DataPublishRegex)
	frequencyBandChangedString := parseFieldInBracketsAsString(line, formats.FrequencyBandChangedRegex)
	frequencyBandRollBackDonestirng := parseFieldInBracketsAsString(line, formats.FrequencyBandRollbackDoneRegex)

	lastServerCommTimeFromSeconds := parseTimeFieldFromSeconds(line, formats.LastServerCommunicationTimeRegex)
	lastDcStartTimeFromSeconds := parseTimeFieldFromSeconds(line, formats.LastDcStartTimeRegex)

	if dcUID == "" &&
		lastServerCommTimeFromSeconds.Year() < 1000 &&
		lastDcStartTimeFromSeconds.Year() < 1000 &&
		dataPublishString == "" &&
		indexCollectionString == "" &&
		frequencyBandChangedString == "" &&
		frequencyBandRollBackDonestirng == "" {
		return nil
	}

	locality := parseFieldInBracketsAsString(line, formats.LocalityRegex)
	region := parseFieldInBracketsAsString(line, formats.RegionRegex)
	timezone := parseFieldInBracketsAsString(line, formats.TimezoneRegex)
	globalFtpAddress := parseFieldInBracketsAsString(line, formats.GlobalFtpAddressRegex)
	targetFirmwareVersion := parseFieldInBracketsAsString(line, formats.TargetFirmwareVersionRegex)
	dcDistroTargetFirmwareVersion := parseFieldInBracketsAsString(line, formats.DcDistroTargetFirmwareVersionRegex)

	indexCollection := tryParseIntFromString(indexCollectionString)
	dataPublish := tryParseIntFromString(dataPublishString)
	frequencyBandChanged := tryParseIntFromString(frequencyBandChangedString) == 1
	frequencyBandRollBackDone := tryParseIntFromString(frequencyBandRollBackDonestirng) == 1

	result := models.SettingsPayload{
		DcUID:            dcUID,
		Locality:         locality,
		Region:           region,
		Timezone:         timezone,
		GlobalFtpAddress: globalFtpAddress}

	result.TargetFirmwareVersion = targetFirmwareVersion
	result.DcDistroTargetFirmwareVersion = dcDistroTargetFirmwareVersion
	result.IndexCollection = indexCollection
	result.DataPublish = dataPublish
	result.FrequencyBandChanged = frequencyBandChanged
	result.FrequencyBandRollBackDone = frequencyBandRollBackDone

	return &result
}

func parseServiceLevelPayload(line string) *models.ServiceLevelPayload {
	meterModeString := parseFieldInBracketsAsString(line, formats.MeterModeRegex)
	maxActivePowerstring := parseFieldInBracketsAsString(line, formats.MaxActivePowerRegex)
	loadSheddingDailyEnergyBudgetString := parseFieldInBracketsAsString(line, formats.LoadSheddingDailyEnergyBudgetRegex)
	localSheddingDailyEnergyBudgetString := parseFieldInBracketsAsString(line, formats.LocalSheddingDailyEnergyBudgetRegex)
	inServiceString := parseFieldInBracketsAsString(line, formats.InServiceRegex)
	if meterModeString == "" &&
		maxActivePowerstring == "" &&
		loadSheddingDailyEnergyBudgetString == "" &&
		localSheddingDailyEnergyBudgetString == "" &&
		inServiceString == "" {
		// if we could not parse any of these fields, than most likely we will not be able to parse the remaining
		return nil
	}

	meterMode := tryParseIntFromString(meterModeString)
	maxActivePower := tryParseIntFromString(maxActivePowerstring)
	loadSheddingDailyEnergyBudget := tryParseIntFromString(loadSheddingDailyEnergyBudgetString)
	localSheddingDailyEnergyBudget := tryParseIntFromString(localSheddingDailyEnergyBudgetString)
	inService := tryParseIntFromString(inServiceString) == 1

	startHourDailyCycle := parseFieldInBracketsAsString(line, formats.StartHourDailyCycleRegex)
	name := parseFieldInBracketsAsString(line, formats.NameRegex)

	hourlyEnergyLimits := parseHourlyEnergyLimits(line, formats.HourlyEnergyLimitsRegex)
	localHourlyEnergyLimits := parseHourlyEnergyLimits(line, formats.LocalHourlyEnergyLimitsRegex)

	result := models.ServiceLevelPayload{
		MeterMode:                      meterMode,
		MaxActivePower:                 maxActivePower,
		LoadSheddingDailyEnergyBudget:  loadSheddingDailyEnergyBudget,
		LocalSheddingDailyEnergyBudget: localSheddingDailyEnergyBudget,
		InService:                      inService,
		StartHourDailyCycle:            startHourDailyCycle,
		Name:                           name}
	result.HourlyEnergyLimits = hourlyEnergyLimits
	result.LocalHourlyEnergyLimits = localHourlyEnergyLimits
	return &result
}

func parseHourlyEnergyLimits(line string, energyLimitRegex string) [24]models.HourlyEnergyLimit {
	var result [24]models.HourlyEnergyLimit
	hourlyLimitsString := parseFieldInDoubleBracketsAsString(line, energyLimitRegex)
	if hourlyLimitsString != "" {
		limitParts := strings.Split(hourlyLimitsString, " ")
		if len(limitParts) == 24 {
			for i, val := range limitParts {
				result[i] = models.HourlyEnergyLimit{HourNumber: i, Limit: tryParseIntFromString(val)}
			}

			return result
		}

		// This is an unexpected format, we always expect the limits string to be 24-long
		log.Println(limitParts)
	}

	return result
}

func parseSmcAddressPayload(line string) *models.SmcAddressParams {
	smcUID := parseFieldInBracketsAsString(line, formats.SMCUIDRegex)
	physicalAddress := parseFieldInBracketsAsString(line, formats.PhysicalAddressRegex)
	logicalAddress := parseFieldInBracketsAsString(line, formats.LogicalAddressRegex)
	shortAddressString := parseFieldInBracketsAsString(line, formats.ShortAddressRegex)
	if smcUID == "" && physicalAddress == "" && logicalAddress == "" && shortAddressString == "" {
		return nil
	}

	lastJoiningDate := parseDateTimeField(line, formats.LastJoiningDateRegex)
	shortAddress := tryParseIntFromString(shortAddressString)

	result := models.SmcAddressParams{
		SmcUID:          smcUID,
		PhysicalAddress: physicalAddress,
		LogicalAddress:  logicalAddress,
		ShortAddress:    shortAddress,
		LastJoiningDate: lastJoiningDate}
	return &result
}

func parseSmcConfigPayload(line string) *models.SmcConfigPayload {
	customerSerialNumber := parseFieldInBracketsAsString(line, formats.CustomerSerialNumberRegex)
	physicalAddress := parseFieldInBracketsAsString(line, formats.PhysicalAddressRegex)
	smcStatus := parseFieldInBracketsAsString(line, formats.SmcStatusRegex)
	nextHopString := parseFieldInBracketsAsString(line, formats.NextHopRegex)

	if customerSerialNumber == "" && physicalAddress == "" && smcStatus == "" && nextHopString == "" {
		return nil
	}

	currentApp1Fw := parseFieldInBracketsAsString(line, formats.CurrentApp1FwRegex)
	currentApp2Fw := parseFieldInBracketsAsString(line, formats.CurrentApp2FwRegex)
	currentPlcFw := parseFieldInBracketsAsString(line, formats.CurrentPlcFwRegex)

	lastSuccessfulDlmsResponseDate := parseDateTimeField(line, formats.LastSuccessfulDlmsResponseDateRegex)
	nextHop := tryParseIntFromString(nextHopString)

	result := models.SmcConfigPayload{}
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

func parsePodConfigPayload(line string) *models.PodConfigPayload {
	serialNumberString := parseFieldInBracketsAsString(line, formats.SerialNumberRegex)
	phaseString := parseFieldInBracketsAsString(line, formats.PhaseRegex)
	positionInSmcString := parseFieldInBracketsAsString(line, formats.PositionInSmcRegex)
	softwareFirmwareVersion := parseFieldInBracketsAsString(line, formats.SoftwareFirmwareVersionRegex)
	if serialNumberString == "" && phaseString == "" && positionInSmcString == "" && softwareFirmwareVersion == "" {
		return nil
	}

	serialNumber := tryParseIntFromString(serialNumberString)
	phase := tryParseIntFromString(phaseString)
	positionInSmc := tryParseIntFromString(positionInSmcString)

	result := models.PodConfigPayload{
		Phase:                   phase,
		SerialNumber:            serialNumber,
		PositionInSmc:           positionInSmc,
		SoftwareFirmwareVersion: softwareFirmwareVersion}
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

func parseTimeRange(line string) *models.TimeRange {
	from := parseDateTimeField(line, formats.TimeRangeFromRegex)

	if from.Year() < 1500 {
		from = parseTimeFieldFromSeconds(line, formats.TimeRangeStartTicksRegex)
	}

	to := parseDateTimeField(line, formats.TimeRangeToRegex)
	if to.Year() < 1500 {
		to = parseTimeFieldFromSeconds(line, formats.TimeRangeEndTicksRegex)
	}

	if from.Year() > 1500 && to.Year() > 1500 {
		result := models.TimeRange{From: from, To: to}
		return &result
	}

	return nil
}
