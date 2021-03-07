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
