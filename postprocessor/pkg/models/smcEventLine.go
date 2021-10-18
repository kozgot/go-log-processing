package models

import "time"

// SmcEventLine shows the events happening to an smc over time.
type SmcEventLine struct {
	SmcUID  string
	Address AddressDetails
	Events  []SmcEvent
}

// SmcEvent is an event happening to a specific smc at a specific time.
type SmcEvent struct {
	Time      time.Time
	EventType EventType
	Label     string
}

// EventType represents the type of the event happening in this moment.
type EventType int64

const (
	// UnknownEventType is the default value of EventType.
	UnknownEventType EventType = iota
	NewSmc
	SmcJoined
	SmcConnected
	ConnectionFailed
	IndexCollected
)
