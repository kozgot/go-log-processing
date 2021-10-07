package models

import "time"

// InfoParams contains the parsed info parameters.
type InfoParams struct {
	MessageType    string               // one of 'ROUTING', 'JOIN', 'STATUS', or 'DC'
	RoutingMessage RoutingTableParams   // no smc UID for this kind of entries
	JoinMessage    SmcJoinMessageParams // has an SMC UID
	StatusMessage  StatusMessageParams  // no smc UID for this kind of entries
	DCMessage      DCMessageParams      // has an SMC UID
}

// RoutingTableParams contains the parsed routing table message parameters.
type RoutingTableParams struct {
	Address        string
	NextHopAddress string
	RouteCost      int
	HopCount       int
	WeakLink       int
	ValidTimeMins  int
}

// SmcJoinMessageParams contains the parsed SMC join message parameters.
type SmcJoinMessageParams struct {
	Ok         bool
	Response   string
	JoinType   string
	SmcAddress SmcAddressParams
}

// SmcAddressParams contains the parsed SMC address parameters of an SMC join log entry.
type SmcAddressParams struct {
	SmcUID          string
	PhysicalAddress string
	LogicalAddress  string
	ShortAddress    int
	LastJoiningDate time.Time
}

// StatusMessageParams contains the parsed message lines from plc_manager.log.
type StatusMessageParams struct {
	Message    string
	StatusByte string
}

// DCMessageParams contains the parsed info level messages that have been sent or received by the dc.
type DCMessageParams struct {
	IsInComing       bool
	SourceOrDestName string
	MessageType      string // todo: prepare enums for message types
	Payload          *DcMessagePayload
}

// DcMessagePayload contains the parsed payload of info level messages that have been sent or received by the dc.
type DcMessagePayload struct {
	SmcUID         string
	PodUID         string
	ServiceLevelID int
	Value          int
	Time           time.Time

	TimeRange *TimeRange

	ConnectOrDisconnectPayload *ConnectOrDisconnectPayload
	DLMSLogPayload             *DLMSLogPayload
	IndexPayload               *IndexPayload
	MessagePayload             *MessagePayload
	SettingsPayload            *SettingsPayload
	ServiceLevelPayload        *ServiceLevelPayload
	SmcAddressPayload          *SmcAddressParams
	SmcConfigPayload           *SmcConfigPayload
	PodConfigPayload           *PodConfigPayload
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
	StartHourDailyCycle            string // eg. 20h, todo: better type??
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