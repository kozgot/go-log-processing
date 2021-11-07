package models

import "time"

// WarningParams contains the parsed warning parameters.
type WarningParams struct {
	Name          string
	SmcUID        string
	UID           int
	Priority      int
	Retry         int
	FileName      string
	Creation      time.Time
	MinLaunchTime time.Time

	Details              *ErrorParams
	JoinMessageParams    *SmcJoinMessageParams
	TimeoutParams        *TimelineOutParams
	LostConnectionParams *LostConnectionParams
}

// LostConnectionParams contains a parsed lost connaction warning entry.
type LostConnectionParams struct {
	Type      int
	Reason    string
	ClientID  string
	URL       string
	Topic     string
	Timeout   int
	Connected bool
}

// TimelineOutParams contains the parsed timeout WARN level log entries.
type TimelineOutParams struct {
	Protocol string
	URL      string
}

// Equals checks equality.
func (w *WarningParams) Equals(other WarningParams) bool {
	if w.Creation != other.Creation ||
		w.Details != other.Details ||
		w.Name != other.Name ||
		w.SmcUID != other.SmcUID ||
		w.MinLaunchTime != other.MinLaunchTime ||
		w.Priority != other.Priority ||
		w.UID != other.UID ||
		w.Retry != other.Retry {
		return false
	}

	// Details
	if w.Details == nil && other.Details != nil || w.Details != nil && other.Details == nil {
		return false
	}
	if w.Details != nil && other.Details != nil && !w.Details.Equals(*other.Details) {
		return false
	}

	// LostConnectionParams
	if w.LostConnectionParams == nil && other.LostConnectionParams != nil ||
		w.LostConnectionParams != nil && other.LostConnectionParams == nil {
		return false
	}
	if w.LostConnectionParams != nil && other.LostConnectionParams != nil &&
		!w.LostConnectionParams.Equals(*other.LostConnectionParams) {
		return false
	}

	// TimeoutParams
	if w.TimeoutParams == nil && other.TimeoutParams != nil ||
		w.TimeoutParams != nil && other.TimeoutParams == nil {
		return false
	}
	if w.TimeoutParams != nil && other.TimeoutParams != nil &&
		!w.TimeoutParams.Equals(*other.TimeoutParams) {
		return false
	}

	// JoinMessageParams
	if w.JoinMessageParams == nil && other.JoinMessageParams != nil ||
		w.JoinMessageParams != nil && other.JoinMessageParams == nil {
		return false
	}
	if w.JoinMessageParams != nil && other.JoinMessageParams != nil &&
		!w.JoinMessageParams.Equals(*other.JoinMessageParams) {
		return false
	}

	return true
}

// Equals checks equality.
func (l *LostConnectionParams) Equals(other LostConnectionParams) bool {
	if l.Type != other.Type ||
		l.Reason != other.Reason ||
		l.ClientID != other.ClientID ||
		l.URL != other.URL ||
		l.Topic != other.Topic ||
		l.Timeout != other.Timeout ||
		l.Connected != other.Connected {
		return false
	}

	return true
}

// Equals checks equality.
func (t *TimelineOutParams) Equals(other TimelineOutParams) bool {
	if t.Protocol != other.Protocol || t.URL != other.URL {
		return false
	}

	return true
}
