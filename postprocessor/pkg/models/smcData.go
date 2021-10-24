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

// AddressDetails contains data related to an SMC's address.
type AddressDetails struct {
	ShortAddress    int
	PhysicalAddress string
	LogicalAddress  string
	URL             string
}
