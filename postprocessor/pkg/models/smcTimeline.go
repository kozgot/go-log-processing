package models

import "time"

// SmcData stores data related to an smc.
type SmcData struct {
	SmcUID                    string
	Address                   AddressDetails
	CustomerSerialNumber      string
	Pods                      []Pod
	LastSuccesfulDlmsResponse time.Time
	LastJoiningDate           time.Time
}

// Pod stores data related to a pod.
type Pod struct {
	UID            string
	SmcUID         string
	SerialNumber   int
	Phase          int
	ServiceLevelID int
	PositionInSmc  int
}

// SmcTimeline represents the state changes of an smc over time.
type SmcTimeline struct {
	From     time.Time
	To       time.Time
	SmcData  SmcData
	Sections []TimelineSection
}

// TimelineSection represents the state of in a period of time.
type TimelineSection struct {
	From  time.Time
	To    time.Time
	State SmcState
	Label string
}

// AddressDetails contains data related to an SMC's address.
type AddressDetails struct {
	ShortAddress    int
	PhysicalAddress string
	LogicalAddress  string
	URL             string
}
