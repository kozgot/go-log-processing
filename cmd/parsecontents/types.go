package parsecontents

import "time"

// ErrorParams contains the parsed error parameters
type ErrorParams struct {
	ErrorCode   int
	Message     string
	Severity    int
	Description string
	Source      string
}

// InfoMessageParams contains the parsed info level messages that have been sent or recieved by the dc
type InfoMessageParams struct {
	SentByDC bool
	Sender   string
	Receiver string
	Message  string
	Payload  InfoPayload
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
	ErrorParams   ErrorParams
	FileName      string
}

// InfoParams contains the parsed info parameters
type InfoParams struct {
	RoutingParams       RoutingTableParams
	JoinMessageParams   SmcJoinMessageParams
	StatusMessageParams StatusMessageParams
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
	Ok        bool
	Response  string
	JoinType  string
	SmcConfig SmcConfigParams
}

type SmcConfigParams struct {
	SmcUID          string
	PhysicalAddress string
	LogicalAddress  string
	ShortAddress    string
	LastJoiningDate time.Time
}

// StatusMessageParams contains the parsed message lines from plc_manager.log
type StatusMessageParams struct {
	Message    string
	StatusByte string
}

// ParsedLine contains a parsed line from the log file
type ParsedLine struct {
	Timestamp     time.Time
	Level         string
	ErrorParams   ErrorParams
	WarningParams WarningParams
	InfoParams    InfoParams
}

/* Info Message Payload types*/
// InfoMessagePayload contains the parsed payload of info level messages that have been sent or recieved by the dc
type InfoPayload struct {
	CommonProps CommonMessageProps
	TimeRange   TimeRange

	ConnectOrDisconnectPayload ConnectOrDisconnectPayload
	DLMSLogPayload             DLMSLogPayload
	IndexPayload               IndexPayload
	MessagePayload             MessagePayload
	TarifSettingsPayload       TarifSettingsPayload
	SettingsPayload            SettingsPayload
	ServiceLevelPayload        ServiceLevelPayload
	SmcAddressPayload          SmcAddressPayload
	SmcConfigPayload           SmcConfigPayload
	PodConfigPayload           PodConfigPayload
}

type CommonMessageProps struct {
	// <--[new smc]--(SVI) dc18-smc18 esetén ez egy smc uid
	// <--[last index]--(DB)
	// <--[read index profiles]--(SMC) smc_uid[dc18-smc10] (1)
	// --[statistics]-->(DB) itt nincs payload
	SmcUID         string
	PodUID         string
	ServiceLevelId int
	Value          int
	Time           time.Time // it is only present in 3 types of messages
}

type SettingsPayload struct {
	// <--[settings]--(DB) and --[settings]-->(DB)     // exactly the same properties
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
	//  <--[service_level]--(DB)
	MeterMode                      int
	StartHourDailyCycle            string // pl.: 20h better type??
	LoadSheddingDailyEnergyBudget  int
	LocalSheddingDailyEnergyBudget int
	MaxActivePower                 int
	InService                      int
	Name                           string
	HourlyEnergyLimits             []HourlyEnergyLimit
	LocalHourlyEnergyLimits        []HourlyEnergyLimit
}

type SmcAddressPayload struct {
	// <--[smc address]--(DB)
	LogicalAddress  string
	ShortAddress    string
	PhysicalAddress string
	LastJoiningDate time.Time
}

type SmcConfigPayload struct {
	// <--[smc configuration]--(DB)
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
	// --[message]-->(SVI)
	Current float32
	Total   float32
	URL     string
	Topic   string
}

type PodConfigPayload struct {
	// <--[pod configuration]--(DB)
	SerialNumber            string // mitgh be int
	Phase                   int
	PositionInSmc           int
	SoftwareFirmwareVersion string
}

type TarifSettingsPayload struct {
	//  <--[tariff_settings]--(DB) ebbol csak 1 van
	ActiveCalendar            Calendar
	PassiveCalendar           Calendar
	ActivePassiveCalendarTime int
}

type TimeRange struct {
	// --[read index low profiles]-->(SMC)
	// <--[consumption]--(SMC)

	From time.Time
	To   time.Time
}

type DLMSLogPayload struct {
	// --[DLMS Logs]-->(SVI)
	Time3            time.Time
	DLMSRequestTime  time.Time
	DLMSResponseTime time.Time
	DLMSError        string
}

type IndexPayload struct {
	// <--[index]--(SMC) & --[index]-->(SVI)
	PreviousTime  time.Time // it might be ticks or something (eg. 1591776000)
	PreviousValue int
	SerialNumber  string // might be int
}

type ConnectOrDisconnectPayload struct {
	// --[connect]-->(SVI) & --[disconnect]-->(SVI)
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
	SeasonProfile string // for now, the exact type is unkown
	WeekProfile   string // for now, the exact type is unkown
	DayProfile    string // for now, the exact type is unkown
}

type HourlyEnergyLimit struct {
	HourNumber int
	Limit      int
}
