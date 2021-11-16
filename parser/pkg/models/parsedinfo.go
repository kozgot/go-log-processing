package models

import (
	"time"
)

// InfoParams contains the parsed info parameters.
type InfoParams struct {
	EntryType               EntryType
	RoutingMessage          *RoutingTableParams      // no smc UID for this kind of entries
	JoinMessage             *SmcJoinMessageParams    // has an SMC UID
	StatusMessage           *StatusMessageParams     // no smc UID for this kind of entries
	DCMessage               *DCMessageParams         // has an SMC UID
	ConnectionAttempt       *ConnectionAttemptParams // has an SMC UID and url
	SmcConfigUpdate         *SmcConfigUpdateParams
	ConnectionReleased      *ConnectionReleasedParams
	InitConnection          *InitConnectionParams
	InternalDiagnosticsData *InternalDiagnosticsData
}

// InternalDiagnosticsData contains a parsed internal diagnostics log entry.
type InternalDiagnosticsData struct {
	SmcUID                         string
	LastSuccessfulDlmsResponseDate time.Time
}

// InitConnectionParams contains a parsed initialize dlms connection log entry.
type InitConnectionParams struct {
	URL string
}

// ConnectionAttemptParams contains a parsed connection attempt log entry.
type ConnectionAttemptParams struct {
	URL    string
	SmcUID string
	At     string // eg. (@ 000A)
}

// ConnectionReleasedParams contains a parsed connection attempt log entry.
type ConnectionReleasedParams struct {
	URL string
}

// SmcConfigUpdateParams contains a parsed SMC config update log entry.
// Update SMC configuration in DB smc_uid[dc18-smc32] physical_address[EEBEDDFFFE6210AD]
//    logical_address[FE80::4021:FF:FE00:000a:61616] short_address[10]
//    last_joining_date[Wed Jun 10 09:20:14 2020]! (distribution_controller_plc_interface.cc::68).
type SmcConfigUpdateParams struct {
	PhysicalAddress string
	LogicalAddress  string
	ShortAddress    int
	LastJoiningDate time.Time
	SmcUID          string
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
	MessageType      DCMessageType
	Payload          *DcMessagePayload
}
