package models

// InfoEntryType represents the type of an INFO level log entry.
type InfoEntryType int64

const (
	// UnknownInfoType is the default value of EntryType.
	UnknownInfoType     InfoEntryType = iota
	Routing                           // Routing Table: Addr[0x0008]...
	NetworkStatus                     // <--[Network status]--(PLC)
	SMCJoin                           // SMC Join OK [Confirmed] <--     or     SMC Join NOT OK [Rejected] <--
	ConnectionAttempt                 // Attempt to connect to SMC
	SmcConfigUpdate                   // Update SMC configuration in DB
	DCMessage                         // entries that are described by the DCMessageType enum
	ConnectionReleased                // Successfully Released DLMS connection
	InitDLMSConnection                // Initialize DLMS connection
	InternalDiagnostics               // SMC internal diagnostics
)
