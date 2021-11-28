package models

import "time"

// WarningParams contains the parsed warning parameters.
type WarningParams struct {
	WarningType WarningEntryType

	TaskFailedWarningParams *TaskFailedWarningParams
	JoinMessageParams       *SmcJoinMessageParams
	TimeoutParams           *TimeOutParams
	LostConnectionParams    *LostConnectionParams
}

// TaskFailedWarningParams contains a parsed task failed type, warn level log entry.
type TaskFailedWarningParams struct {
	Name          string // The name of the failed task.
	SmcUID        string
	UID           int
	Priority      int
	Retry         int
	FileName      string
	Creation      time.Time
	MinLaunchTime time.Time

	Details *ErrorParams
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

// TimeOutParams contains the parsed timeout WARN level log entries.
type TimeOutParams struct {
	Protocol string
	URL      string
}
