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

// todo: define enum for these
const (
	RountingMessageType = "ROUTING"
	JoinMessageType     = "JOIN"
	StatusMessageType   = "STATUS"
	DCMessageType       = "DC"
)

// InfoParams contains the parsed info parameters
type InfoParams struct {
	MessageType    string // one of 'ROUTING', 'JOIN', 'STATUS', or 'DC'
	RoutingMessage RoutingTableParams
	JoinMessage    SmcJoinMessageParams
	StatusMessage  StatusMessageParams
	DCMessage      DCMessageParams
}

// RoutingTableParams contains the parsed routing table message parameters
type RoutingTableParams struct { // Routing Table: Addr[0x0018] NextHopAddr[0x001F] RouteCost[15] HopCount[0] WeakLink[2] ValidTime[240] (min)
	Address        string
	NextHopAddress string
	RouteCost      int
	HopCount       int
	WeakLink       int
	ValidTimeMins  int
}

// SmcJoinMessageParams contains the parsed SMC join message parameters
type SmcJoinMessageParams struct { // SMC Join OK [Confirmed] <-- [join_type[LBA] smc_uid[..] physical_address[..] logical_address[..] short_address[14] last_joining_date[..]]--(PLC)
	Ok         bool
	Response   string
	JoinType   string
	SmcAddress SmcAddressParams
}

type SmcAddressParams struct {
	SmcUID          string
	PhysicalAddress string
	LogicalAddress  string
	ShortAddress    string
	LastJoiningDate time.Time
}

// StatusMessageParams contains the parsed message lines from plc_manager.log
type StatusMessageParams struct { // LOADNG_SEQ_NUM_REPORTED status_byte[0xA5]<--[Network status]--(PLC)
	Message    string
	StatusByte string
}

// DCMessageParams contains the parsed info level messages that have been sent or recieved by the dc
type DCMessageParams struct {
	Sender      string
	Receiver    string
	MessageType string // todo: prepare enums for message types
	Payload     InfoPayload
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
	SmcUID         string
	PodUID         string
	ServiceLevelId int
	Value          int
	Time           time.Time

	TimeRange TimeRange

	ConnectOrDisconnectPayload ConnectOrDisconnectPayload
	DLMSLogPayload             DLMSLogPayload
	IndexPayload               IndexPayload
	MessagePayload             MessagePayload
	TarifSettingsPayload       TarifSettingsPayload
	SettingsPayload            SettingsPayload
	ServiceLevelPayload        ServiceLevelPayload
	SmcAddressPayload          SmcAddressParams
	SmcConfigPayload           SmcConfigPayload
	PodConfigPayload           PodConfigPayload
}

type SettingsPayload struct {
	// <--[settings]--(DB) and --[settings]-->(DB)
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
	StartHourDailyCycle            string // eg. 20h, todo: better type??
	LoadSheddingDailyEnergyBudget  int
	LocalSheddingDailyEnergyBudget int
	MaxActivePower                 int
	InService                      int
	Name                           string
	HourlyEnergyLimits             []HourlyEnergyLimit
	LocalHourlyEnergyLimits        []HourlyEnergyLimit
}

type HourlyEnergyLimit struct {
	HourNumber int
	Limit      int
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
	SerialNumber            int
	Phase                   int
	PositionInSmc           int
	SoftwareFirmwareVersion string
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
	SerialNumber  int
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

// parsing of these can wait, does not seem so important as it only has a single occurance
type TarifSettingsPayload struct {
	//  <--[tariff_settings]--(DB) ebbol csak 1 van
	ActiveCalendar            Calendar
	PassiveCalendar           Calendar
	ActivePassiveCalendarTime int
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
