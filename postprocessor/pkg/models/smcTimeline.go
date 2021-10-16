package models

import "time"

type SmcTimeline struct {
	From     time.Time
	To       time.Time
	SmcUID   string
	Address  AddressDetails
	Sections []TimelineSection
}

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
	IndexCollected
)

type AddressDetails struct {
	ShortAddress    string
	PhysicalAddress string
	LogicalAddress  string
	URL             string
}
