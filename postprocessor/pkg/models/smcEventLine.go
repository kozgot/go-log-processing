package models

import "time"

// SmcEventLine shows the events happening to an smc over time.
type SmcEventLine struct {
	SmcData SmcData
	Events  []SmcEvent
}

// SmcEvent is an event happening to a specific smc at a specific time.
type SmcEvent struct {
	Time      time.Time
	EventType EventType
	Label     string
	SmcUID    string
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
	ConnectionFailed
	IndexCollectionStarted
	IndexRead
	IndexLowProfileGenericReceived
	IndexHighProfileGenericReceived
	DLMSLogsSent
	ConfigurationChanged
	InternalDiagnostics
	StatisticsSent
)