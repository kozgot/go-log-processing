package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// SmcEvent is an event happening to a specific smc at a specific time.
type SmcEvent struct {
	Time            time.Time
	EventType       EventType
	EventTypeString string
	Label           string
	SmcUID          string
	DataPayload     SmcData
}

// Serialize serializes an smc event and returns a byte array.
func (e SmcEvent) Serialize() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		fmt.Println("Can't serialize event ", e)
	}

	return bytes
}

// EventType represents the type of the event happening in this moment.
type EventType int64

const (
	// UnknownEventType is the default value of EventType.
	UnknownEventType EventType = iota
	NewSmc
	PodConfiguration
	SmcJoined
	ConnectionAttempt
	StartToConnect
	SmcAddressUpdated
	SmcAddressInvalidated
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
	ConfigurationReadFromDB
	ConfigurationUpdated
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

	case ConfigurationReadFromDB:
		return "ConfigurationReadFromDB"

	case ConfigurationUpdated:
		return "ConfigurationUpdated"

	case InternalDiagnostics:
		return "InternalDiagnostics"

	case StatisticsSent:
		return "StatisticsSent"

	case SmcAddressInvalidated:
		return "SmcAddressInvalidated"

	default:
		return "None"
	}
}
