package models

// Pod stores data related to a pod.
type Pod struct {
	UID            string
	SmcUID         string
	SerialNumber   int
	Phase          int
	ServiceLevelID int
	PositionInSmc  int
}
