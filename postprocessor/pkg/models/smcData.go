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

// Equals check equality.
func (s *SmcData) Equals(other SmcData) bool {
	if s.SmcUID != other.SmcUID ||
		s.CustomerSerialNumber != other.CustomerSerialNumber ||
		s.LastSuccesfulDlmsResponse != other.LastSuccesfulDlmsResponse ||
		s.LastJoiningDate != other.LastJoiningDate {
		return false
	}

	return s.Address.Equals(other.Address) && s.EqualPodLists(other)
}

// ContainsPod checks if the smc data contains the given pod.
func (s *SmcData) ContainsPod(pod Pod) bool {
	for _, p := range s.Pods {
		if p.UID == pod.UID {
			return true
		}
	}
	return false
}

// EqualPodLists checks if the pod lists are equal in two smc data obbjects.
func (s *SmcData) EqualPodLists(other SmcData) bool {
	if len(s.Pods) != len(other.Pods) {
		return false
	}

	for _, p := range s.Pods {
		if !other.ContainsPod(p) {
			return false
		}
	}

	return true
}
