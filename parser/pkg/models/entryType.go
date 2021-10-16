package models

// EntryType represents the type of the log entry.
type EntryType int64

const (
	// UnknownEntryType is the default value of EntryType.
	UnknownEntryType EntryType = iota
	Routing                    // Routing Table: Addr[0x0008]...
	NetworkStatus              // <--[Network status]--(PLC)
	SMCJoin                    // SMC Join OK [Confirmed] <--     or     SMC Join NOT OK [Rejected] <--

	ConnectionAttempt   // Attempt to connect to SMC
	SmcConfigUpdate     // Update SMC configuration in DB
	DCMessage           // entries that are described by the DCMessageType enum
	ConnectionReleased  // Successfully Released DLMS connection
	InitDLMSConnection  // Initialize DLMS connection
	InternalDiagnostics // SMC internal diagnostics
)
