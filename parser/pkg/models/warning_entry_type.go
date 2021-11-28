package models

// WarningEntryType represents the type of a warning log entry.
type WarningEntryType int64

const (
	// UnknownWarningTypeType is the default value of WarningEntryType.
	UnknownWarningTypeType WarningEntryType = iota
	TaskFailedWarning
	JoinRejectedWarning
	TimeoutWarning
	ConnectionLostWarning
)
