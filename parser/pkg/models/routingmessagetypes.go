package models

// todo: define enum for these
const (
	RountingMessageType = "ROUTING"
	JoinMessageType     = "JOIN"
	StatusMessageType   = "STATUS"
	DCMessageType       = "DC"
	ConnectionAttempt   = "CONNECTION_ATTEMPT"
	SmcConfigUpdate     = "SMC_CONFIG_UPDATE"
	ConnectionReleased  = "CONNECTION_RELEASED"
	InitDlmsConnection  = "INITIALIZE_DLMS_CONNECTION"
	InternalDiagnostics = "SMC_INTERNAL_DIAGNOSTICS"
)
