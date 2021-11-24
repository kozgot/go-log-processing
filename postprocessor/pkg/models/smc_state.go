package models

// SmcState represents the state the smc is in.
type SmcState int64

const (
	// UnknownSmcState is the default value of SmcState.
	UnknownSmcState SmcState = iota
	New
	Joined
	Connecting
	Error
	Disconnected
	CollectingIndex
)

// SmcStateToString returns the string representation of an SMC state.
func SmcStateToString(state SmcState) string {
	switch state {
	case New:
		return "New"
	case Joined:
		return "Joined"
	case Connecting:
		return "Connecting"
	case Error:
		return "Error"
	case CollectingIndex:
		return "CollectingIndex"
	case Disconnected:
		return "Disconnected"
	case UnknownSmcState:
		return "UnknownSmcState"
	default:
		return "None"
	}
}
