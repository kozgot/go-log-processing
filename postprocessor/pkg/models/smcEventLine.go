package models

import "time"

// SmcEventLine shows the events happening to an smc over time.
type SmcEventLine struct {
	SmcData SmcData
	Events  []SmcEvent
}

// SmcEvent is an event happening to a specific smc at a specific time.
type SmcEvent struct {
	Time            time.Time
	EventType       EventType
	EventTypeString string
	Label           string
	SmcUID          string
	DataPayload     SmcData
}

// EventType represents the type of the event happening in this moment.
type EventType int64

const (
	// UnknownEventType is the default value of EventType.
	UnknownEventType EventType = iota
	NewSmc
	PodConfiguration
	SmcJoined
	SmcConnected
	ConnectionAttempt
	StartToConnect
	SmcAddressUpdated
	TimeoutWarning
	JoinRejectedWarning
	DLMSError
	InitConnection
	ConnectionReleased
	IndexCollectionStarted
	IndexRead
	IndexLowProfileGenericReceived
	IndexHighProfileGenericReceived
	DLMSLogsSent
	ConfigurationChanged
	InternalDiagnostics
	StatisticsSent
)

func EventTypeToString(eventType EventType) string {
	switch eventType {
	case UnknownEventType:
		return "UnknownEventType"

	case NewSmc:
		return "NewSmc"

	case PodConfiguration:
		return "PodConfiguration"

	case SmcJoined:
		return "SmcJoined"

	case SmcConnected:
		return "SmcConnected"

	case ConnectionAttempt:
		return "ConnectionAttempt"

	case StartToConnect:
		return "StartToConnect"

	case SmcAddressUpdated:
		return "SmcAddressUpdated"

	case TimeoutWarning:
		return "TimeoutWarning"

	case JoinRejectedWarning:
		return "JoinRejectedWarning"

	case DLMSError:
		return "DLMSError"

	case InitConnection:
		return "InitConnection"

	case ConnectionReleased:
		return "ConnectionReleased"

	case IndexCollectionStarted:
		return "IndexCollectionStarted"

	case IndexRead:
		return "IndexRead"

	case IndexLowProfileGenericReceived:
		return "IndexLowProfileGenericReceived"

	case IndexHighProfileGenericReceived:
		return "IndexHighProfileGenericReceived"

	case DLMSLogsSent:
		return "DLMSLogsSent"

	case ConfigurationChanged:
		return "ConfigurationChanged"

	case InternalDiagnostics:
		return "InternalDiagnostics"

	case StatisticsSent:
		return "StatisticsSent"

	default:
		return "None"
	}
}
