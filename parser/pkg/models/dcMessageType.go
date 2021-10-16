package models

// DCMessageType represents the type of a dc message log entry.
type DCMessageType int64

const (
	// UnknownDCMessage is the default value of DCMessageType.
	UnknownDCMessage        DCMessageType = iota
	IndexReceived                         // <--[index]--(SMC)
	MessageSentToSVI                      // --[message]-->(SVI)
	PodConfig                             // <--[pod configuration]--(DB)
	SmcAddress                            // <--[smc address]--(DB)
	SmcConfig                             // <--[smc configuration]--(DB)
	ServiceLevel                          // <--[service_level]--(DB)
	Settings                              // --[settings]-->(DB)
	DLMSLogs                              // --[DLMS Logs]-->(SVI)
	NewSmc                                // <--[new smc]--(SVI)
	IndexLowProfileGeneric                // <--[index low profile generic]--(SMC)
	IndexHighProfileGeneric               // <--[index high profile generic]--(SMC)
	Connect                               // --[connect]-->(PLC) or --[connect]-->(SVI) or --[connect]-->(UDS)
	Statistics                            // --[statistics]-->(SVI)
	ReadIndexLowProfiles                  // --[read index low profiles]-->(SMC)
	ReadIndexProfiles                     // <--[read index profiles]--(SMC)
)

// ParseDCmessageTypeFromString parses the dc message type from a string representation.
func ParseDCmessageTypeFromString(messageTypeString string) DCMessageType {
	switch messageTypeString {
	case "index":
		return IndexReceived

	case "message":
		return MessageSentToSVI

	case "pod configuration":
		return PodConfig

	case "smc address":
		return SmcAddress

	case "smc configuration":
		return SmcConfig

	case "service_level":
		return ServiceLevel

	case "DLMS Logs":
		return DLMSLogs

	case "new smc":
		return NewSmc

	case "settings":
		return Settings

	case "index high profile generic":
		return IndexHighProfileGeneric

	case "index low profile generic":
		return IndexLowProfileGeneric

	case "connect":
		return Connect

	case "statistics":
		return Statistics

	case "read index low profiles":
		return ReadIndexLowProfiles

	case "read index profiles":
		return ReadIndexProfiles

	default:
		return UnknownDCMessage
	}
}
