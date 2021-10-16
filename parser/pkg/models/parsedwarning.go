package models

import "time"

// WarningParams contains the parsed warning parameters.
type WarningParams struct {
	Name              string
	SmcUID            string
	UID               int
	Priority          int
	Retry             int
	Creation          time.Time
	MinLaunchTime     time.Time
	Details           ErrorParams
	FileName          string
	JoinMessageParams SmcJoinMessageParams
	TimeoutParams     TimelineOutParams
}

// TimelineOutParams contains the parsed timeout WARN level log entries.
type TimelineOutParams struct {
	Protocol string
	URL      string
}
