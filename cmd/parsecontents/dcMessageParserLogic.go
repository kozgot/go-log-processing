package parsecontents

import (
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

	// Parse the time[] field of the message. It can be a formatted date or in a date represented by a timestamp in seconds.
	dateTime := parseDateTimeField(line, DateTimeFieldRegex)
	if dateTime.Year() > 1000 {
		payload.Time = dateTime // for some reason, it can parse the long date format to int, so that needs to be handled as well (hence the if-else)
	} else {
		datefromSeconds := parseTimeFieldFromTimeStamp(line, TimeTicksRegex)
		if datefromSeconds.Year() > 1000 {
			payload.Time = datefromSeconds
		}
	}

	payload.TimeRange = parseTimeRange(line)

	payload.ConnectOrDisconnectPayload = parseConnectOrDisconnectPayload(line)
	payload.DLMSLogPayload = parseDLMSLogPayload(line)
	payload.IndexPayload = parseIndexPayload(line)
	payload.MessagePayload = parseMessagePayload(line)
	payload.PodConfigPayload = parsePodConfigPayload(line)
	payload.SmcConfigPayload = parseSmcConfigPayload(line)
	payload.SmcAddressPayload = parseSmcAddressPayload(line)
	payload.SettingsPayload = parseSettingsPayload(line)
	payload.ServiceLevelPayload = parseServiceLevelPayload(line)

	return payload
}

func parseTimeRange(line string) TimeRange {
	timeRange := TimeRange{}

	from := parseDateTimeField(line, TimeRangeFromRegex)
	if from.Year() > 1000 {
		timeRange.From = from
	} else {
		start := parseTimeFieldFromTimeStamp(line, TimeRangeStartTicksRegex)
		if start.Year() > 1000 {
			timeRange.From = start
		}
	}

	to := parseDateTimeField(line, TimeRangeToRegex)
	if to.Year() > 1000 {
		timeRange.To = to
	} else {
		end := parseTimeFieldFromTimeStamp(line, TimeRangeEndTicksRegex)
		if end.Year() > 1000 {
			timeRange.To = end
		}
	}

	return timeRange
}

func parseTimeFieldFromTimeStamp(line string, timeStampRegex string) time.Time {
	seconds := tryParseInt64FromString(parseFieldInBracketsAsString(line, timeStampRegex))
	if seconds != 0 {
		dateTimeFromsSecs := time.Unix(seconds, 0)
		return dateTimeFromsSecs
	}

	return time.Time{}
}

func parseConnectOrDisconnectPayload(line string) ConnectOrDisconnectPayload {
	result := ConnectOrDisconnectPayload{}

	// todo
	/*
			type ConnectOrDisconnectPayload struct {
			Type      int
			ClientId  string
			URL       string
			Topic     string
			Timeout   int
			Connected bool
		}
	*/

	return result
}

func parseDLMSLogPayload(line string) DLMSLogPayload {
	result := DLMSLogPayload{}

	// todo
	/*
			type DLMSLogPayload struct {
			DLMSRequestTime  time.Time
			DLMSResponseTime time.Time
			DLMSError        string
		}
	*/

	return result
}

func parseIndexPayload(line string) IndexPayload {
	result := IndexPayload{}

	// todo
	/*
			type IndexPayload struct {
			PreviousTime  time.Time // it might be ticks or something (eg. 1591776000)
			PreviousValue int
			SerialNumber  int
		}
	*/

	return result
}

func parseMessagePayload(line string) MessagePayload {
	result := MessagePayload{}

	// todo
	/*
			type MessagePayload struct {
			Current float32
			Total   float32
			URL     string
			Topic   string
		}
	*/

	return result
}

func parseSettingsPayload(line string) SettingsPayload {
	result := SettingsPayload{}

	// todo
	/*
			type SettingsPayload struct {
			DcUID                         string
			Locality                      string
			Region                        string
			Timezone                      string
			GlobalFtpAddress              string
			TargetFirmwareVersion         string
			IndexCollection               int
			DataPublish                   int
			LastServerCommunicationTime   time.Time
			DcDistroTargetFirmwareVersion string
			LastDcStartTime               time.Time // it might be ticks or something (eg. 1591780709)
			FrequencyBandChanged          bool
			FrequencyBandRollBackDone     bool
		}
	*/

	return result
}

func parseServiceLevelPayload(line string) ServiceLevelPayload {
	result := ServiceLevelPayload{}

	// todo
	/*
			type ServiceLevelPayload struct {
			MeterMode                      int
			StartHourDailyCycle            string // eg. 20h, todo: better type??
			LoadSheddingDailyEnergyBudget  int
			LocalSheddingDailyEnergyBudget int
			MaxActivePower                 int
			InService                      int
			Name                           string
			HourlyEnergyLimits             []HourlyEnergyLimit
			LocalHourlyEnergyLimits        []HourlyEnergyLimit
		}
	*/

	return result
}

func parseSmcAddressPayload(line string) SmcAddressParams {
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

	return result
}

func parseSmcConfigPayload(line string) SmcConfigPayload {
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

	return result
}

func parsePodConfigPayload(line string) PodConfigPayload {
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

	return result
}
