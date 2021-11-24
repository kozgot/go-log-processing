package contentparser

import (
	"log"
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type MessageEntryParser struct {
	line models.EntryWithLevelAndTimestamp
}

func (messageEntryParser *MessageEntryParser) Parse() *models.DCMessageParams {
	dcMessageParams := models.DCMessageParams{}

	source := messageEntryParser.parseSource()
	dest := messageEntryParser.parseDestination()

	if source == "" && dest == "" {
		// it is not a message entry
		return nil
	}

	if source != "" {
		dcMessageParams.IsInComing = true
		dcMessageParams.SourceOrDestName = source
		messageTypeString := common.ParseFieldInBracketsAsString(
			messageEntryParser.line.Rest, formats.IncomingMessageTypeRegex,
		)
		dcMessageParams.MessageType = models.ParseDCmessageTypeFromString(messageTypeString)
	} else if dest != "" {
		dcMessageParams.IsInComing = false
		dcMessageParams.SourceOrDestName = dest
		dcMessageTypeString := common.ParseFieldInBracketsAsString(
			messageEntryParser.line.Rest, formats.OutGoingMessageTypeRegex,
		)
		dcMessageParams.MessageType = models.ParseDCmessageTypeFromString(dcMessageTypeString)
	}

	dcMessageParams.Payload = messageEntryParser.parseDCMessagePayload(dcMessageParams.MessageType, dest)

	return &dcMessageParams
}

func (messageEntryParser *MessageEntryParser) parseSource() string {
	inComingMessageSource := common.ParseFieldInParenthesesAsString(
		messageEntryParser.line.Rest, formats.IncomingMessageSourceRegex)
	return inComingMessageSource
}

func (messageEntryParser *MessageEntryParser) parseDestination() string {
	outGoingMessageSource := common.ParseFieldInParenthesesAsString(
		messageEntryParser.line.Rest, formats.OutGoingMessageDestRegex)
	return outGoingMessageSource
}

func (messageEntryParser *MessageEntryParser) parsePayloadTime() time.Time {
	// Parse the time[] field of the message.
	// It can be a formatted date or in a date represented by a timestamp in seconds.
	dateTime := common.ParseDateTimeField(messageEntryParser.line.Rest, formats.DateTimeFieldRegex)
	if common.IsValidDate(dateTime) {
		return dateTime
	}
	datefromSeconds := common.ParseTimeFieldFromSeconds(messageEntryParser.line.Rest, formats.TimeTicksRegex)
	if common.IsValidDate(datefromSeconds) {
		return datefromSeconds
	}

	return time.Time{}
}

func (messageEntryParser *MessageEntryParser) parseDCMessagePayload(
	messageType models.DCMessageType,
	destination string,
) *models.DcMessagePayload {
	payload := models.DcMessagePayload{}

	payload.SmcUID = common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.SmcUIDRegex)
	payload.PodUID = common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.PodUIDRegex)
	payload.ServiceLevelID = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.ServiceLevelIDRegex),
	)
	payload.Value = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.ValueRegex),
	)
	payload.Time = messageEntryParser.parsePayloadTime()
	payload.TimeRange = common.ParseTimeRange(messageEntryParser.line.Rest)

	switch messageType {
	case models.NewSmc:
		payload.SmcUID = messageEntryParser.parseNewSmcUID()

	case models.MessageSentToSVI:
		payload.MessagePayload = messageEntryParser.parseMessagePayload()

	case models.Connect:
		if destination == "PLC" {
			payload.ConnectToPLCPayload = messageEntryParser.parseConnectToPLC()
		} else {
			// destination is SVI or UDS
			payload.ConnectOrDisconnectPayload = messageEntryParser.parseConnectOrDisconnectPayload()
		}

	case models.PodConfig:
		payload.PodConfigPayload = messageEntryParser.parsePodConfigPayload()

	case models.SmcConfig:
		payload.SmcConfigPayload = messageEntryParser.parseSmcConfigPayload()

	case models.SmcAddress:
		payload.SmcAddressPayload = messageEntryParser.parseSmcAddressPayload()

	case models.ServiceLevel:
		payload.ServiceLevelPayload = messageEntryParser.parseServiceLevelPayload()

	case models.Settings:
		payload.SettingsPayload = messageEntryParser.parseSettingsPayload()

	case models.DLMSLogs:
		payload.DLMSLogPayload = messageEntryParser.parseDLMSLogPayload()

	case models.IndexReceived:
		payload.IndexPayload = messageEntryParser.parseIndexPayload()

	case models.Consumption:
		// NOOP: all the aparams are parsed in the root payload property
		break

	case models.IndexLowProfileGeneric:
		payload.GenericIndexProfilePayload = messageEntryParser.parseGenericIndexProfile()

	case models.IndexHighProfileGeneric:
		payload.GenericIndexProfilePayload = messageEntryParser.parseGenericIndexProfile()

	case models.ReadIndexLowProfiles:
		payload.ReadIndexLowProfilesEntryPayload = messageEntryParser.parseReadIndexLowProfilesEntry()

	case models.ReadIndexProfiles:
		payload.ReadIndexProfilesEntryPayload = messageEntryParser.parseReadIndexProfilesEntry()

	case models.Statistics:
		if destination == "SVI" {
			// The destination could be 'DB' as well, but in that case, we have no more params to parse.
			payload.StatisticsEntryPayload = messageEntryParser.parseStatisticsEntry()
		}

	case models.UnknownDCMessage:
		// NOOP
		break
	}

	return &payload
}

func (
	messageEntryParser *MessageEntryParser,
) parseReadIndexLowProfilesEntry() *models.ReadIndexLowProfilesEntryPayload {
	timeRange := common.ParseTimeRange(messageEntryParser.line.Rest)
	result := models.ReadIndexLowProfilesEntryPayload{}
	result.To = timeRange.To
	result.From = timeRange.From
	result.SmcUID = common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.SMCUIDRegex)

	return &result
}

func (messageEntryParser *MessageEntryParser) parseReadIndexProfilesEntry() *models.ReadIndexProfilesEntryPayload {
	result := models.ReadIndexProfilesEntryPayload{}
	result.SmcUID = common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.SMCUIDRegex)

	// Get the count part between the parentheses.
	// <--[read index profiles]--(SMC) smc_uid[dc18-smc9] (6) (smart_meter_cabinet.cc::190)
	correctCount := 4
	parts := strings.Split(messageEntryParser.line.Rest, "(")
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

func (messageEntryParser *MessageEntryParser) parseStatisticsEntry() *models.StatisticsEntryPayload {
	result := models.StatisticsEntryPayload{}
	result.Type = common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.StatisticsTypeRegex)
	result.SourceID = common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.StatisticsSourceIDRegex)
	result.Time = common.ParseTimeFieldFromSeconds(messageEntryParser.line.Rest, formats.TimeTicksRegex)
	valueString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.StatisticsValueRegex)

	value := common.TryParseFloat64FromString(valueString)
	result.Value = value

	return &result
}

func (messageEntryParser *MessageEntryParser) parseGenericIndexProfile() *models.GenericIndexProfilePayload {
	result := models.GenericIndexProfilePayload{}
	capturePeriodString := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.IndexProfileCapturePeriodRegex,
	)
	capturePeriod := common.TryParseIntFromString(capturePeriodString)

	captureObjectsString := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.IndexProfileCaptureObjectsRegex,
	)
	captureObjects := common.TryParseIntFromString(captureObjectsString)

	result.CaptureObjects = captureObjects
	result.CapturePeriod = capturePeriod

	return &result
}

func (messageEntryParser *MessageEntryParser) parseConnectToPLC() *models.ConnectToPLCPayload {
	result := models.ConnectToPLCPayload{}
	result.Interface = common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.ConnectToPLCIfaceRegex)
	result.DestinationAddress = common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.ConnectToPLCDestAddressRegex,
	)

	return &result
}

func (messageEntryParser *MessageEntryParser) parseNewSmcUID() string {
	// parse the smc uid from this:
	//  <--[new smc]--(SVI) dc18-smc32 (distribution_controller_initializer.cc::280)

	minLengthIfContainsSeparator := 2

	// get the part after (SVI)
	firstPart := strings.Split(messageEntryParser.line.Rest, ")")
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

func (messageEntryParser *MessageEntryParser) parseConnectOrDisconnectPayload() *models.ConnectOrDisconnectPayload {
	resultType := common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.ConnectOrDisconnectTypeRegex))
	clientID := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.ClientIDRegex)
	URL := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.URLRegex)
	topic := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.TopicRegex)
	timeout := common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.TimeoutRegex))
	connected := common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.ConnectedRegex)) == 1

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

func (messageEntryParser *MessageEntryParser) parseDLMSLogPayload() *models.DLMSLogPayload {
	requestTimeFromSeconds := common.ParseTimeFieldFromMilliSeconds(
		messageEntryParser.line.Rest, formats.DLMSRequestTimeRegex,
	)
	responseTimeFromSeconds := common.ParseTimeFieldFromMilliSeconds(
		messageEntryParser.line.Rest, formats.DLMSResponseTimeRegex,
	)
	DLMSError := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.DLMSErrorRegex)

	if requestTimeFromSeconds.Year() > 1500 || responseTimeFromSeconds.Year() > 1500 || DLMSError != "" {
		result := models.DLMSLogPayload{
			DLMSRequestTime:  requestTimeFromSeconds,
			DLMSResponseTime: responseTimeFromSeconds,
			DLMSError:        DLMSError}
		return &result
	}

	return nil
}

func (messageEntryParser *MessageEntryParser) parseIndexPayload() *models.IndexPayload {
	previousValueString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.PreviousValueRegex)
	serialNumberString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.SerialNumberRegex)
	previousTimeFromSeconds := common.ParseTimeFieldFromSeconds(messageEntryParser.line.Rest, formats.PreviousTimeRegex)
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

func (messageEntryParser *MessageEntryParser) parseMessagePayload() *models.MessagePayload {
	currentString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.CurrentRegex)
	totalString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.TotalRegex)
	if currentString == "" && totalString == "" {
		return nil
	}

	result := models.MessagePayload{}

	result.Current = common.TryParseFloat64FromString(currentString)
	result.Total = common.TryParseFloat64FromString(
		common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.TotalRegex))
	result.URL = common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.URLRegex)
	result.Topic = common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.TopicRegex)

	return &result
}

func (messageEntryParser *MessageEntryParser) parseSettingsPayload() *models.SettingsPayload {
	dcUID := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.DcUIDRegex)

	indexCollectionString := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.IndexCollectionRegex,
	)
	dataPublishString := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.DataPublishRegex,
	)
	frequencyBandChangedString := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.FrequencyBandChangedRegex,
	)
	frequencyBandRollBackDonestirng := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.FrequencyBandRollbackDoneRegex,
	)

	lastServerCommTimeFromSeconds := common.ParseTimeFieldFromSeconds(
		messageEntryParser.line.Rest,
		formats.LastServerCommunicationTimeRegex,
	)
	lastDcStartTimeFromSeconds := common.ParseTimeFieldFromSeconds(
		messageEntryParser.line.Rest,
		formats.LastDcStartTimeRegex,
	)

	if dcUID == "" &&
		lastServerCommTimeFromSeconds.Year() < 1000 &&
		lastDcStartTimeFromSeconds.Year() < 1000 &&
		dataPublishString == "" &&
		indexCollectionString == "" &&
		frequencyBandChangedString == "" &&
		frequencyBandRollBackDonestirng == "" {
		return nil
	}

	locality := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.LocalityRegex)
	region := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.RegionRegex)
	timezone := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.TimezoneRegex)
	globalFtpAddress := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.GlobalFtpAddressRegex)
	targetFirmwareVersion := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.TargetFirmwareVersionRegex,
	)
	dcDistroTargetFirmwareVersion := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.DcDistroTargetFirmwareVersionRegex,
	)

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

func (messageEntryParser *MessageEntryParser) parseServiceLevelPayload() *models.ServiceLevelPayload {
	meterModeString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.MeterModeRegex)
	maxActivePowerstring := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.MaxActivePowerRegex)
	loadSheddingDailyEnergyBudgetString := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.LoadSheddingDailyEnergyBudgetRegex)
	localSheddingDailyEnergyBudgetString := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.LocalSheddingDailyEnergyBudgetRegex)
	inServiceString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.InServiceRegex)
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

	startHourDailyCycle := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.StartHourDailyCycleRegex,
	)
	name := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.NameRegex)

	hourlyEnergyLimits := messageEntryParser.parseHourlyEnergyLimits(formats.HourlyEnergyLimitsRegex)
	localHourlyEnergyLimits := messageEntryParser.parseHourlyEnergyLimits(formats.LocalHourlyEnergyLimitsRegex)

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

func (messageEntryParser *MessageEntryParser) parseHourlyEnergyLimits(
	energyLimitRegex string,
) [24]models.HourlyEnergyLimit {
	var result [24]models.HourlyEnergyLimit
	hoursInADay := 24
	hourlyLimitsString := common.ParseFieldInDoubleBracketsAsString(messageEntryParser.line.Rest, energyLimitRegex)
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

func (messageEntryParser *MessageEntryParser) parseSmcAddressPayload() *models.SmcAddressParams {
	smcUID := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.SMCUIDRegex)
	physicalAddress := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.PhysicalAddressRegex)
	logicalAddress := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.LogicalAddressRegex)
	shortAddressString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.ShortAddressRegex)
	if smcUID == "" && physicalAddress == "" && logicalAddress == "" && shortAddressString == "" {
		return nil
	}

	lastJoiningDate := common.ParseDateTimeField(messageEntryParser.line.Rest, formats.LastJoiningDateRegex)
	shortAddress := common.TryParseIntFromString(shortAddressString)

	result := models.SmcAddressParams{
		SmcUID:          smcUID,
		PhysicalAddress: physicalAddress,
		LogicalAddress:  logicalAddress,
		ShortAddress:    shortAddress,
		LastJoiningDate: lastJoiningDate}
	return &result
}

func (messageEntryParser *MessageEntryParser) parseSmcConfigPayload() *models.SmcConfigPayload {
	customerSerialNumber := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.CustomerSerialNumberRegex,
	)
	physicalAddress := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.PhysicalAddressRegex)
	smcStatus := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.SmcStatusRegex)
	nextHopString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.NextHopRegex)

	if customerSerialNumber == "" && physicalAddress == "" && smcStatus == "" && nextHopString == "" {
		return nil
	}

	currentApp1Fw := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.CurrentApp1FwRegex)
	currentApp2Fw := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.CurrentApp2FwRegex)
	currentPlcFw := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.CurrentPlcFwRegex)

	lastSuccessfulDlmsResponseDate := common.ParseDateTimeField(
		messageEntryParser.line.Rest,
		formats.LastSuccessfulRespDateRegex,
	)
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

func (messageEntryParser *MessageEntryParser) parsePodConfigPayload() *models.PodConfigPayload {
	serialNumberString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.SerialNumberRegex)
	phaseString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.PhaseRegex)
	positionInSmcString := common.ParseFieldInBracketsAsString(messageEntryParser.line.Rest, formats.PositionInSmcRegex)
	softwareFirmwareVersion := common.ParseFieldInBracketsAsString(
		messageEntryParser.line.Rest,
		formats.SoftwareFirmwareVersionRegex,
	)
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
