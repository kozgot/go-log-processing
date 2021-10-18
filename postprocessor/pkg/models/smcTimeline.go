package models

import "time"

// SmcTimeline represents the state changes of an smc over time.
type SmcTimeline struct {
	From     time.Time
	To       time.Time
	SmcUID   string
	Address  AddressDetails
	Sections []TimelineSection
}

// TimelineSection represents the state of in a period of time.
type TimelineSection struct {
	From  time.Time
	To    time.Time
	State SmcState
	Label string
}

// SmcState represents the state the smc is in.
type SmcState int64

const (
	// UnknownSmcState is the default value of SmcState.
	UnknownSmcState SmcState = iota
	New
	Joined
	Connected
	Failing
	CollectingIndex
)

type AddressDetails struct {
	ShortAddress    string
	PhysicalAddress string
	LogicalAddress  string
	URL             string
}
