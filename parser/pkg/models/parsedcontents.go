package models

import "time"

// ParsedLine contains a parsed line from the log file.
type ParsedLine struct {
	Timestamp     time.Time
	Level         string
	ErrorParams   ErrorParams
	WarningParams WarningParams
	InfoParams    InfoParams
}

// ErrorParams contains the parsed error parameters
type ErrorParams struct {
	ErrorCode   int
	Message     string
	Severity    int
	Description string
	Source      string
}

// WarningParams contains the parsed warning parameters
type WarningParams struct {
	Name          string
	SmcUID        string
	UID           int
	Priority      int
	Retry         int
	Creation      time.Time
	MinLaunchTime time.Time
	Details       ErrorParams
	FileName      string
}

// InfoParams contains the parsed info parameters
type InfoParams struct {
	MessageType    string // one of 'ROUTING', 'JOIN', 'STATUS', or 'DC'
	RoutingMessage RoutingTableParams
	JoinMessage    SmcJoinMessageParams
	StatusMessage  StatusMessageParams
	DCMessage      DCMessageParams
}

// RoutingTableParams contains the parsed routing table message parameters
type RoutingTableParams struct {
	Address        string
	NextHopAddress string
	RouteCost      int
	HopCount       int
	WeakLink       int
	ValidTimeMins  int
}

// SmcJoinMessageParams contains the parsed SMC join message parameters
type SmcJoinMessageParams struct {
	Ok         bool
	Response   string
	JoinType   string
	SmcAddress SmcAddressParams
}

type SmcAddressParams struct {
	SmcUID          string
	PhysicalAddress string
	LogicalAddress  string
	ShortAddress    int
	LastJoiningDate time.Time
}

// StatusMessageParams contains the parsed message lines from plc_manager.log
type StatusMessageParams struct {
	Message    string
	StatusByte string
}

// DCMessageParams contains the parsed info level messages that have been sent or recieved by the dc
type DCMessageParams struct {
	IsInComing       bool
	SourceOrDestName string
	MessageType      string // todo: prepare enums for message types
	Payload          *DcMessagePayload
}

/* Dc Message Payload types*/
// DcMessagePayload contains the parsed payload of info level messages that have been sent or recieved by the dc
type DcMessagePayload struct {
	SmcUID         string
	PodUID         string
	ServiceLevelId int
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

type HourlyEnergyLimit struct {
	HourNumber int
	Limit      int
}

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

type MessagePayload struct {
	Current float64
	Total   float64
	URL     string
	Topic   string
}

type PodConfigPayload struct {
	SerialNumber            int
	Phase                   int
	PositionInSmc           int
	SoftwareFirmwareVersion string
}

type TimeRange struct {
	From time.Time
	To   time.Time
}

type DLMSLogPayload struct {
	DLMSRequestTime  time.Time
	DLMSResponseTime time.Time
	DLMSError        string
}

type IndexPayload struct {
	PreviousTime  time.Time // it might be ticks or something (eg. 1591776000)
	PreviousValue int
	SerialNumber  int
}

type ConnectOrDisconnectPayload struct {
	Type      int
	ClientId  string
	URL       string
	Topic     string
	Timeout   int
	Connected bool
}

type Calendar struct {
	CalendarName CalendarName
}

type CalendarName struct {
	IsActive      bool
	SeasonProfile string // for now, the exact type is unknown
	WeekProfile   string // for now, the exact type is unknown
	DayProfile    string // for now, the exact type is unknown
}
