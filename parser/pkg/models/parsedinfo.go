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

// Equals checks equality.
func (i *InfoParams) Equals(other InfoParams) bool {
	if i.EntryType != other.EntryType {
		return false
	}

	// RoutingMessage
	if i.RoutingMessage != nil && other.RoutingMessage == nil ||
		i.RoutingMessage == nil && other.RoutingMessage != nil {
		return false
	}
	if i.RoutingMessage != nil && !i.RoutingMessage.Equals(*other.RoutingMessage) {
		return false
	}

	// JoinMessage
	if i.JoinMessage != nil && other.JoinMessage == nil ||
		i.JoinMessage == nil && other.JoinMessage != nil {
		return false
	}
	if i.JoinMessage != nil && !i.JoinMessage.Equals(*other.JoinMessage) {
		return false
	}

	// StatusMessage
	if i.StatusMessage != nil && other.StatusMessage == nil ||
		i.StatusMessage == nil && other.StatusMessage != nil {
		return false
	}
	if i.StatusMessage != nil && !i.StatusMessage.Equals(*other.StatusMessage) {
		return false
	}

	// DCMessage
	if i.DCMessage != nil && other.DCMessage == nil ||
		i.DCMessage == nil && other.DCMessage != nil {
		return false
	}
	if i.DCMessage != nil && !i.DCMessage.Equals(*other.DCMessage) {
		return false
	}

	// ConnectionAttempt
	if i.ConnectionAttempt != nil && other.ConnectionAttempt == nil ||
		i.ConnectionAttempt == nil && other.ConnectionAttempt != nil {
		return false
	}
	if i.ConnectionAttempt != nil && !i.ConnectionAttempt.Equals(*other.ConnectionAttempt) {
		return false
	}

	// SmcConfigUpdate
	if i.SmcConfigUpdate != nil && other.SmcConfigUpdate == nil ||
		i.SmcConfigUpdate == nil && other.SmcConfigUpdate != nil {
		return false
	}
	if i.SmcConfigUpdate != nil && !i.SmcConfigUpdate.Equals(*other.SmcConfigUpdate) {
		return false
	}

	// ConnectionReleased
	if i.ConnectionReleased != nil && other.ConnectionReleased == nil ||
		i.ConnectionReleased == nil && other.ConnectionReleased != nil {
		return false
	}
	if i.ConnectionReleased != nil && !i.ConnectionReleased.Equals(*other.ConnectionReleased) {
		return false
	}

	// InitConnection
	if i.InitConnection != nil && other.InitConnection == nil ||
		i.InitConnection == nil && other.InitConnection != nil {
		return false
	}
	if i.InitConnection != nil && !i.InitConnection.Equals(*other.InitConnection) {
		return false
	}

	// InternalDiagnosticsData
	if i.InternalDiagnosticsData != nil && other.InternalDiagnosticsData == nil ||
		i.InternalDiagnosticsData == nil && other.InternalDiagnosticsData != nil {
		return false
	}
	if i.InternalDiagnosticsData != nil && !i.InternalDiagnosticsData.Equals(*other.InternalDiagnosticsData) {
		return false
	}

	return true
}

// InternalDiagnosticsData contains a parsed internal diagnostics log entry.
type InternalDiagnosticsData struct {
	SmcUID                         string
	LastSuccessfulDlmsResponseDate time.Time
}

// Equals checks equality.
func (i *InternalDiagnosticsData) Equals(other InternalDiagnosticsData) bool {
	return i.LastSuccessfulDlmsResponseDate == other.LastSuccessfulDlmsResponseDate &&
		i.SmcUID == other.SmcUID
}

// InitConnectionParams contains a parsed initialize dlms connection log entry.
type InitConnectionParams struct {
	URL string
}

// Equals checks equality.
func (i *InitConnectionParams) Equals(other InitConnectionParams) bool {
	return i.URL == other.URL
}

// ConnectionAttemptParams contains a parsed connection attempt log entry.
type ConnectionAttemptParams struct {
	URL    string
	SmcUID string
	At     string // eg. (@ 000A)
}

// Equals checks equality.
func (c *ConnectionAttemptParams) Equals(other ConnectionAttemptParams) bool {
	return c.URL == other.URL &&
		c.SmcUID == other.SmcUID &&
		c.At == other.At
}

// ConnectionReleasedParams contains a parsed connection attempt log entry.
type ConnectionReleasedParams struct {
	URL string
}

// Equals checks equality.
func (c *ConnectionReleasedParams) Equals(other ConnectionReleasedParams) bool {
	return c.URL == other.URL
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

// Equals checks equality.
func (s *SmcConfigUpdateParams) Equals(other SmcConfigUpdateParams) bool {
	return s.PhysicalAddress == other.PhysicalAddress &&
		s.LogicalAddress == other.LogicalAddress &&
		s.ShortAddress == other.ShortAddress &&
		s.LastJoiningDate == other.LastJoiningDate &&
		s.SmcUID == other.SmcUID
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

// Equals checks equality.
func (r *RoutingTableParams) Equals(other RoutingTableParams) bool {
	return r.Address == other.Address &&
		r.NextHopAddress == other.NextHopAddress &&
		r.RouteCost == other.RouteCost &&
		r.HopCount == other.HopCount &&
		r.WeakLink == other.WeakLink &&
		r.ValidTimeMins == other.ValidTimeMins
}

// SmcJoinMessageParams contains the parsed SMC join message parameters.
type SmcJoinMessageParams struct {
	Ok         bool
	Response   string
	JoinType   string
	SmcAddress SmcAddressParams
}

// Equals checks equality.
func (j *SmcJoinMessageParams) Equals(other SmcJoinMessageParams) bool {
	if j.JoinType != other.JoinType ||
		j.Ok != other.Ok ||
		j.Response != other.Response {
		return false
	}

	if !j.SmcAddress.Equals(other.SmcAddress) {
		return false
	}

	return true
}

// SmcAddressParams contains the parsed SMC address parameters of an SMC join log entry.
type SmcAddressParams struct {
	SmcUID          string
	PhysicalAddress string
	LogicalAddress  string
	ShortAddress    int
	LastJoiningDate time.Time
}

// Equals checks equality.
func (a *SmcAddressParams) Equals(other SmcAddressParams) bool {
	if a.LastJoiningDate != other.LastJoiningDate ||
		a.LogicalAddress != other.LogicalAddress ||
		a.PhysicalAddress != other.PhysicalAddress ||
		a.ShortAddress != other.ShortAddress ||
		a.SmcUID != other.SmcUID {
		return false
	}

	return true
}

// StatusMessageParams contains the parsed message lines from plc_manager.log.
type StatusMessageParams struct {
	Message    string
	StatusByte string
}

// Equals checks equality.
func (s *StatusMessageParams) Equals(other StatusMessageParams) bool {
	return s.Message == other.Message && s.StatusByte == other.StatusByte
}

// DCMessageParams contains the parsed info level messages that have been sent or received by the dc.
type DCMessageParams struct {
	IsInComing       bool
	SourceOrDestName string
	MessageType      DCMessageType
	Payload          *DcMessagePayload
}

// Equals checks equality.
func (d *DCMessageParams) Equals(other DCMessageParams) bool {
	if d.IsInComing != other.IsInComing ||
		d.SourceOrDestName != other.SourceOrDestName ||
		d.MessageType != other.MessageType {
		return false
	}

	// Check Payload
	if d.Payload != nil && other.Payload == nil ||
		d.Payload == nil && other.Payload != nil {
		return false
	}
	if d.Payload != nil && !d.Payload.Equals(*other.Payload) {
		return false
	}

	return true
}
