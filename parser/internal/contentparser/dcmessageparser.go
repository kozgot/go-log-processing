package contentparser

import (
	"log"
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func parseDCMessage(lin string) *models.DCMessageParams {
	dcMessageParams := models.DCMessageParams{}

	source := parseDCMessageSource(lin)
	dest := parseDCMessageDest(lin)

	if source == "" && dest == "" {
		// it is not a message entry
		return nil
	}

	if source != "" {
		dcMessageParams.IsInComing = true
		dcMessageParams.SourceOrDestName = source
		messageTypeString := common.ParseFieldInBracketsAsString(lin, formats.IncomingMessageTypeRegex)
		dcMessageParams.MessageType = models.ParseDCmessageTypeFromString(messageTypeString)
	} else if dest != "" {
		dcMessageParams.IsInComing = false
		dcMessageParams.SourceOrDestName = dest
		dcMessageTypeString := common.ParseFieldInBracketsAsString(lin, formats.OutGoingMessageTypeRegex)
		dcMessageParams.MessageType = models.ParseDCmessageTypeFromString(dcMessageTypeString)
	}

	dcMessageParams.Payload = parseDCMessagePayload(lin, dcMessageParams.MessageType, dest)

	return &dcMessageParams
}

func parseDCMessageSource(line string) string {
	inComingMessageSource := common.ParseFieldInParenthesesAsString(line, formats.IncomingMessageSourceRegex)
	return inComingMessageSource
}

func parseDCMessageDest(line string) string {
	outGoingMessageSource := common.ParseFieldInParenthesesAsString(line, formats.OutGoingMessageDestRegex)
	return outGoingMessageSource
}

func parsePayloadTime(line string) time.Time {
	// Parse the time[] field of the message.
	// It can be a formatted date or in a date represented by a timestamp in seconds.
	dateTime := common.ParseDateTimeField(line, formats.DateTimeFieldRegex)
	if common.IsValidDate(dateTime) {
		return dateTime
	}
	datefromSeconds := common.ParseTimeFieldFromSeconds(line, formats.TimeTicksRegex)
	if common.IsValidDate(datefromSeconds) {
		return datefromSeconds
	}

	return time.Time{}
}

func parseDCMessagePayload(line string, messageType models.DCMessageType, destination string) *models.DcMessagePayload {
	payload := models.DcMessagePayload{}

	payload.SmcUID = common.ParseFieldInBracketsAsString(line, formats.SmcUIDRegex)
	payload.PodUID = common.ParseFieldInBracketsAsString(line, formats.PodUIDRegex)
	payload.ServiceLevelID = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(line, formats.ServiceLevelIDRegex))
	payload.Value = common.TryParseIntFromString(common.ParseFieldInBracketsAsString(line, formats.ValueRegex))
	payload.Time = parsePayloadTime(line)
	payload.TimeRange = common.ParseTimeRange(line)

	switch messageType {
	case models.NewSmc:
		payload.SmcUID = parseNewSmcUID(line)

	case models.MessageSentToSVI:
		payload.MessagePayload = parseMessagePayload(line)

	case models.Connect:
		if destination == "PLC" {
			payload.ConnectToPLCPayload = parseConnectToPLC(line)
		} else {
			payload.ConnectOrDisconnectPayload = parseConnectOrDisconnectPayload(line) // destination is SVI or UDS
		}

	case models.PodConfig:
		payload.PodConfigPayload = parsePodConfigPayload(line)

	case models.SmcConfig:
		payload.SmcConfigPayload = parseSmcConfigPayload(line)

	case models.SmcAddress:
		payload.SmcAddressPayload = parseSmcAddressPayload(line)

	case models.ServiceLevel:
		payload.ServiceLevelPayload = parseServiceLevelPayload(line)

	case models.Settings:
		payload.SettingsPayload = parseSettingsPayload(line)

	case models.DLMSLogs:
		payload.DLMSLogPayload = parseDLMSLogPayload(line)

	case models.IndexReceived:
		payload.IndexPayload = parseIndexPayload(line)

	case models.Consumption:
		// NOOP: all the aparams are parsed in the root payload property
		break

	case models.IndexLowProfileGeneric:
		payload.GenericIndexProfilePayload = parseGenericIndexProfile(line)

	case models.IndexHighProfileGeneric:
		payload.GenericIndexProfilePayload = parseGenericIndexProfile(line)

	case models.ReadIndexLowProfiles:
		payload.ReadIndexLowProfilesEntryPayload = parseReadIndexLowProfilesEntry(line)

	case models.ReadIndexProfiles:
		payload.ReadIndexProfilesEntryPayload = parseReadIndexProfilesEntry(line)

	case models.Statistics:
		if destination == "SVI" {
			// The destination could be 'DB' as well, but in that case, we have no more params to parse.
			payload.StatisticsEntryPayload = parseStatisticsEntry(line)
		}

	case models.UnknownDCMessage:
		// NOOP
		break
	}

	return &payload
}

func parseReadIndexLowProfilesEntry(line string) *models.ReadIndexLowProfilesEntryPayload {
	timeRange := common.ParseTimeRange(line)
	result := models.ReadIndexLowProfilesEntryPayload{}
	result.To = timeRange.To
	result.From = timeRange.From
	result.SmcUID = common.ParseFieldInBracketsAsString(line, formats.SMCUIDRegex)

	return &result
}

func parseReadIndexProfilesEntry(line string) *models.ReadIndexProfilesEntryPayload {
	result := models.ReadIndexProfilesEntryPayload{}
	result.SmcUID = common.ParseFieldInBracketsAsString(line, formats.SMCUIDRegex)

	// Get the count part between the parentheses.
	// <--[read index profiles]--(SMC) smc_uid[dc18-smc9] (6) (smart_meter_cabinet.cc::190)
	correctCount := 4
	parts := strings.Split(line, "(")
	if len(parts) < correctCount {
		// Default if we could not parse...
		result.Count = 0
	}

	// eg.: '6) '
	countString := parts[2]

	// trim off the ) and space from the end
	countString = strings.Replace(countString, ") ", "", 1)

	// convert to int
	count := common.TryParseIntFromString(countString)
	result.Count = count

	return &result
}

func parseStatisticsEntry(line string) *models.StatisticsEntryPayload {
	result := models.StatisticsEntryPayload{}
	result.Type = common.ParseFieldInBracketsAsString(line, formats.StatisticsTypeRegex)
	result.SourceID = common.ParseFieldInBracketsAsString(line, formats.StatisticsSourceIDRegex)
	result.Time = common.ParseTimeFieldFromSeconds(line, formats.TimeTicksRegex)
	valueString := common.ParseFieldInBracketsAsString(line, formats.StatisticsValueRegex)

	value := common.TryParseFloat64FromString(valueString)
	result.Value = value

	return &result
}

func parseGenericIndexProfile(line string) *models.GenericIndexProfilePayload {
	result := models.GenericIndexProfilePayload{}
	capturePeriodString := common.ParseFieldInBracketsAsString(line, formats.IndexProfileCapturePeriodRegex)
	capturePeriod := common.TryParseIntFromString(capturePeriodString)

	captureObjectsString := common.ParseFieldInBracketsAsString(line, formats.IndexProfileCaptureObjectsRegex)
	captureObjects := common.TryParseIntFromString(captureObjectsString)

	result.CaptureObjects = captureObjects
	result.CapturePeriod = capturePeriod

	return &result
}

func parseConnectToPLC(line string) *models.ConnectToPLCPayload {
	result := models.ConnectToPLCPayload{}
	result.Interface = common.ParseFieldInBracketsAsString(line, formats.ConnectToPLCIfaceRegex)
	result.DestinationAddress = common.ParseFieldInBracketsAsString(line, formats.ConnectToPLCDestAddressRegex)

	return &result
}

func parseNewSmcUID(line string) string {
	// parse the smc uid from this:
	//  <--[new smc]--(SVI) dc18-smc32 (distribution_controller_initializer.cc::280)

	minLengthIfContainsSeparator := 2

	// get the part after (SVI)
	firstPart := strings.Split(line, ")")
	if len(firstPart) < minLengthIfContainsSeparator {
		return ""
	}

	// get the part before (distribution_controller_initializer.cc::280)
	firstPart = strings.Split(firstPart[1], "(")
	if len(firstPart) < minLengthIfContainsSeparator {
		return ""
	}

	// trim the spces from ' dc18-smc32 '
	result := strings.Trim(firstPart[0], " ")
	return result
}

func parseConnectOrDisconnectPayload(line string) *models.ConnectOrDisconnectPayload {
	resultType := common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(line, formats.ConnectOrDisconnectTypeRegex))
	clientID := common.ParseFieldInBracketsAsString(line, formats.ClientIDRegex)
	URL := common.ParseFieldInBracketsAsString(line, formats.URLRegex)
	topic := common.ParseFieldInBracketsAsString(line, formats.TopicRegex)
	timeout := common.TryParseIntFromString(common.ParseFieldInBracketsAsString(line, formats.TimeoutRegex))
	connected := common.TryParseIntFromString(common.ParseFieldInBracketsAsString(line, formats.ConnectedRegex)) == 1

	if clientID != "" || resultType != 0 || URL != "" || topic != "" || timeout != 0 {
		result := models.ConnectOrDisconnectPayload{
			ClientID:  clientID,
			Type:      resultType,
			URL:       URL,
			Topic:     topic,
			Timeout:   timeout,
			Connected: connected}
		return &result
	}

	return nil
}

func parseDLMSLogPayload(line string) *models.DLMSLogPayload {
	requestTimeFromSeconds := common.ParseTimeFieldFromMilliSeconds(line, formats.DLMSRequestTimeRegex)
	responseTimeFromSeconds := common.ParseTimeFieldFromMilliSeconds(line, formats.DLMSResponseTimeRegex)
	DLMSError := common.ParseFieldInBracketsAsString(line, formats.DLMSErrorRegex)

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
	previousValueString := common.ParseFieldInBracketsAsString(line, formats.PreviousValueRegex)
	serialNumberString := common.ParseFieldInBracketsAsString(line, formats.SerialNumberRegex)
	previousTimeFromSeconds := common.ParseTimeFieldFromSeconds(line, formats.PreviousTimeRegex)
	if previousTimeFromSeconds.Year() < 1000 && serialNumberString == "" && previousValueString == "" {
		return nil
	}

	previousValue := common.TryParseIntFromString(previousValueString)
	serialNumber := common.TryParseIntFromString(serialNumberString)

	result := models.IndexPayload{}
	result.PreviousTime = previousTimeFromSeconds
	result.PreviousValue = previousValue
	result.SerialNumber = serialNumber

	return &result
}

func parseMessagePayload(line string) *models.MessagePayload {
	currentString := common.ParseFieldInBracketsAsString(line, formats.CurrentRegex)
	totalString := common.ParseFieldInBracketsAsString(line, formats.TotalRegex)
	if currentString == "" && totalString == "" {
		return nil
	}

	result := models.MessagePayload{}

	result.Current = common.TryParseFloat64FromString(currentString)
	result.Total = common.TryParseFloat64FromString(common.ParseFieldInBracketsAsString(line, formats.TotalRegex))
	result.URL = common.ParseFieldInBracketsAsString(line, formats.URLRegex)
	result.Topic = common.ParseFieldInBracketsAsString(line, formats.TopicRegex)

	return &result
}

func parseSettingsPayload(line string) *models.SettingsPayload {
	dcUID := common.ParseFieldInBracketsAsString(line, formats.DcUIDRegex)

	indexCollectionString := common.ParseFieldInBracketsAsString(line, formats.IndexCollectionRegex)
	dataPublishString := common.ParseFieldInBracketsAsString(line, formats.DataPublishRegex)
	frequencyBandChangedString := common.ParseFieldInBracketsAsString(line, formats.FrequencyBandChangedRegex)
	frequencyBandRollBackDonestirng := common.ParseFieldInBracketsAsString(line, formats.FrequencyBandRollbackDoneRegex)

	lastServerCommTimeFromSeconds := common.ParseTimeFieldFromSeconds(line, formats.LastServerCommunicationTimeRegex)
	lastDcStartTimeFromSeconds := common.ParseTimeFieldFromSeconds(line, formats.LastDcStartTimeRegex)

	if dcUID == "" &&
		lastServerCommTimeFromSeconds.Year() < 1000 &&
		lastDcStartTimeFromSeconds.Year() < 1000 &&
		dataPublishString == "" &&
		indexCollectionString == "" &&
		frequencyBandChangedString == "" &&
		frequencyBandRollBackDonestirng == "" {
		return nil
	}

	locality := common.ParseFieldInBracketsAsString(line, formats.LocalityRegex)
	region := common.ParseFieldInBracketsAsString(line, formats.RegionRegex)
	timezone := common.ParseFieldInBracketsAsString(line, formats.TimezoneRegex)
	globalFtpAddress := common.ParseFieldInBracketsAsString(line, formats.GlobalFtpAddressRegex)
	targetFirmwareVersion := common.ParseFieldInBracketsAsString(line, formats.TargetFirmwareVersionRegex)
	dcDistroTargetFirmwareVersion := common.ParseFieldInBracketsAsString(line, formats.DcDistroTargetFirmwareVersionRegex)

	indexCollection := common.TryParseIntFromString(indexCollectionString)
	dataPublish := common.TryParseIntFromString(dataPublishString)
	frequencyBandChanged := common.TryParseIntFromString(frequencyBandChangedString) == 1
	frequencyBandRollBackDone := common.TryParseIntFromString(frequencyBandRollBackDonestirng) == 1

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
	meterModeString := common.ParseFieldInBracketsAsString(line, formats.MeterModeRegex)
	maxActivePowerstring := common.ParseFieldInBracketsAsString(line, formats.MaxActivePowerRegex)
	loadSheddingDailyEnergyBudgetString := common.ParseFieldInBracketsAsString(
		line,
		formats.LoadSheddingDailyEnergyBudgetRegex)
	localSheddingDailyEnergyBudgetString := common.ParseFieldInBracketsAsString(
		line,
		formats.LocalSheddingDailyEnergyBudgetRegex)
	inServiceString := common.ParseFieldInBracketsAsString(line, formats.InServiceRegex)
	if meterModeString == "" &&
		maxActivePowerstring == "" &&
		loadSheddingDailyEnergyBudgetString == "" &&
		localSheddingDailyEnergyBudgetString == "" &&
		inServiceString == "" {
		// if we could not parse any of these fields, than most likely we will not be able to parse the remaining
		return nil
	}

	meterMode := common.TryParseIntFromString(meterModeString)
	maxActivePower := common.TryParseIntFromString(maxActivePowerstring)
	loadSheddingDailyEnergyBudget := common.TryParseIntFromString(loadSheddingDailyEnergyBudgetString)
	localSheddingDailyEnergyBudget := common.TryParseIntFromString(localSheddingDailyEnergyBudgetString)
	inService := common.TryParseIntFromString(inServiceString) == 1

	startHourDailyCycle := common.ParseFieldInBracketsAsString(line, formats.StartHourDailyCycleRegex)
	name := common.ParseFieldInBracketsAsString(line, formats.NameRegex)

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
	hoursInADay := 24
	hourlyLimitsString := common.ParseFieldInDoubleBracketsAsString(line, energyLimitRegex)
	if hourlyLimitsString != "" {
		limitParts := strings.Split(hourlyLimitsString, " ")
		if len(limitParts) == hoursInADay {
			for i, val := range limitParts {
				result[i] = models.HourlyEnergyLimit{HourNumber: i, Limit: common.TryParseIntFromString(val)}
			}

			return result
		}

		// This is an unexpected format, we always expect the limits string to be 24-long
		log.Println(limitParts)
	}

	return result
}

func parseSmcAddressPayload(line string) *models.SmcAddressParams {
	smcUID := common.ParseFieldInBracketsAsString(line, formats.SMCUIDRegex)
	physicalAddress := common.ParseFieldInBracketsAsString(line, formats.PhysicalAddressRegex)
	logicalAddress := common.ParseFieldInBracketsAsString(line, formats.LogicalAddressRegex)
	shortAddressString := common.ParseFieldInBracketsAsString(line, formats.ShortAddressRegex)
	if smcUID == "" && physicalAddress == "" && logicalAddress == "" && shortAddressString == "" {
		return nil
	}

	lastJoiningDate := common.ParseDateTimeField(line, formats.LastJoiningDateRegex)
	shortAddress := common.TryParseIntFromString(shortAddressString)

	result := models.SmcAddressParams{
		SmcUID:          smcUID,
		PhysicalAddress: physicalAddress,
		LogicalAddress:  logicalAddress,
		ShortAddress:    shortAddress,
		LastJoiningDate: lastJoiningDate}
	return &result
}

func parseSmcConfigPayload(line string) *models.SmcConfigPayload {
	customerSerialNumber := common.ParseFieldInBracketsAsString(line, formats.CustomerSerialNumberRegex)
	physicalAddress := common.ParseFieldInBracketsAsString(line, formats.PhysicalAddressRegex)
	smcStatus := common.ParseFieldInBracketsAsString(line, formats.SmcStatusRegex)
	nextHopString := common.ParseFieldInBracketsAsString(line, formats.NextHopRegex)

	if customerSerialNumber == "" && physicalAddress == "" && smcStatus == "" && nextHopString == "" {
		return nil
	}

	currentApp1Fw := common.ParseFieldInBracketsAsString(line, formats.CurrentApp1FwRegex)
	currentApp2Fw := common.ParseFieldInBracketsAsString(line, formats.CurrentApp2FwRegex)
	currentPlcFw := common.ParseFieldInBracketsAsString(line, formats.CurrentPlcFwRegex)

	lastSuccessfulDlmsResponseDate := common.ParseDateTimeField(line, formats.LastSuccessfulRespDateRegex)
	nextHop := common.TryParseIntFromString(nextHopString)

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
	serialNumberString := common.ParseFieldInBracketsAsString(line, formats.SerialNumberRegex)
	phaseString := common.ParseFieldInBracketsAsString(line, formats.PhaseRegex)
	positionInSmcString := common.ParseFieldInBracketsAsString(line, formats.PositionInSmcRegex)
	softwareFirmwareVersion := common.ParseFieldInBracketsAsString(line, formats.SoftwareFirmwareVersionRegex)
	if serialNumberString == "" && phaseString == "" && positionInSmcString == "" && softwareFirmwareVersion == "" {
		return nil
	}

	serialNumber := common.TryParseIntFromString(serialNumberString)
	phase := common.TryParseIntFromString(phaseString)
	positionInSmc := common.TryParseIntFromString(positionInSmcString)

	result := models.PodConfigPayload{
		Phase:                   phase,
		SerialNumber:            serialNumber,
		PositionInSmc:           positionInSmc,
		SoftwareFirmwareVersion: softwareFirmwareVersion}
	return &result
}
