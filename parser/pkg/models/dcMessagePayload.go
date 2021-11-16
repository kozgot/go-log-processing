package models

import "time"

// DcMessagePayload contains the parsed payload of info level messages that have been sent or received by the dc.
type DcMessagePayload struct {
	SmcUID         string
	PodUID         string
	ServiceLevelID int
	Value          int
	Time           time.Time

	TimeRange *TimeRange

	ConnectOrDisconnectPayload       *ConnectOrDisconnectPayload
	DLMSLogPayload                   *DLMSLogPayload
	IndexPayload                     *IndexPayload
	GenericIndexProfilePayload       *GenericIndexProfilePayload
	MessagePayload                   *MessagePayload
	SettingsPayload                  *SettingsPayload
	ServiceLevelPayload              *ServiceLevelPayload
	SmcAddressPayload                *SmcAddressParams
	SmcConfigPayload                 *SmcConfigPayload
	PodConfigPayload                 *PodConfigPayload
	ConnectToPLCPayload              *ConnectToPLCPayload
	StatisticsEntryPayload           *StatisticsEntryPayload
	ReadIndexLowProfilesEntryPayload *ReadIndexLowProfilesEntryPayload
	ReadIndexProfilesEntryPayload    *ReadIndexProfilesEntryPayload
}

// ReadIndexLowProfilesEntryPayload contains the parsed --[read index low profiles]-->(SMC) entries.
type ReadIndexLowProfilesEntryPayload struct {
	SmcUID string
	From   time.Time
	To     time.Time
}

// ReadIndexProfilesEntryPayload contains the parsed <--[read index profiles]--(SMC) entries.
type ReadIndexProfilesEntryPayload struct {
	SmcUID string
	Count  int
}

// StatisticsEntryPayload contains the parsed statistics log entry sent to the SVI.
type StatisticsEntryPayload struct {
	Type     string
	Value    float64
	Time     time.Time
	SourceID string
}

// GenericIndexProfilePayload contains the parsed index high/low profile generic payload.
type GenericIndexProfilePayload struct {
	CapturePeriod  int
	CaptureObjects int
}

// ConnectToPLCPayload contains the parsed connect to PLC payload.
type ConnectToPLCPayload struct {
	Interface          string
	DestinationAddress string
}

// SettingsPayload contains the parsed settings payload
// of info level messages that have been sent or received by the dc.
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

// ServiceLevelPayload contains the parsed service level related data
// in info level log entries that describe messages sent or received by the dc.
type ServiceLevelPayload struct {
	MeterMode                      int
	StartHourDailyCycle            string // eg. 20h
	LoadSheddingDailyEnergyBudget  int
	LocalSheddingDailyEnergyBudget int
	MaxActivePower                 int
	InService                      bool
	Name                           string
	HourlyEnergyLimits             [24]HourlyEnergyLimit
	LocalHourlyEnergyLimits        [24]HourlyEnergyLimit
}

// HourlyEnergyLimit contains the value and the hour number of an hourly energy limit.
type HourlyEnergyLimit struct {
	HourNumber int
	Limit      int
}

// SmcConfigPayload contains data related to the configuration of an SMC.
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

// MessagePayload contains the parameters of a log entry related to a DC message.
type MessagePayload struct {
	Current float64
	Total   float64
	URL     string
	Topic   string
}

// PodConfigPayload contains the parameters of a log entry related to a pod configuration event.
type PodConfigPayload struct {
	SerialNumber            int
	Phase                   int
	PositionInSmc           int
	SoftwareFirmwareVersion string
}

// TimeRange represents a time range.
type TimeRange struct {
	From time.Time
	To   time.Time
}

// DLMSLogPayload contains data of a log entry related to DLMS Log contents.
type DLMSLogPayload struct {
	DLMSRequestTime  time.Time
	DLMSResponseTime time.Time
	DLMSError        string
}

// IndexPayload contains data of an index log entry.
type IndexPayload struct {
	PreviousTime  time.Time // it might be ticks or something (eg. 1591776000)
	PreviousValue int
	SerialNumber  int
}

// ConnectOrDisconnectPayload contains information related to a connect or disconnect event.
type ConnectOrDisconnectPayload struct {
	Type      int
	ClientID  string
	URL       string
	Topic     string
	Timeout   int
	Connected bool
}

// Calendar contains data related to the calendar of a tariff settings log entry.
type Calendar struct {
	CalendarName CalendarName
}

// CalendarName contains data related to the calendar of a tariff settings log entry.
type CalendarName struct {
	IsActive      bool
	SeasonProfile string // for now, the exact type is unknown
	WeekProfile   string // for now, the exact type is unknown
	DayProfile    string // for now, the exact type is unknown
}
