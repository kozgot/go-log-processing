package models

import "time"

// RoutingEntry contains data from a log entry related to the current state of the routing table.
type RoutingEntry struct {
	TimeStamp      time.Time // the exact time the event happened
	Address        string
	NextHopAddress string
	RouteCost      int
	HopCount       int
	WeakLink       int
	ValidTimeMins  int
}
