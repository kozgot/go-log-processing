package models

import "time"

// WarningParams contains the parsed warning parameters.
type WarningParams struct {
	Name                 string
	SmcUID               string
	UID                  int
	Priority             int
	Retry                int
	Creation             time.Time
	MinLaunchTime        time.Time
	Details              ErrorParams
	FileName             string
	JoinMessageParams    SmcJoinMessageParams
	TimeoutParams        TimelineOutParams
	LostConnectionParams LostConnectionParams
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
