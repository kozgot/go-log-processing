package models

// ProcessedEntries contains all the processed data from all the log entries.
type ProcessedEntries struct {
	SmcEntries     map[string][]SmcEntry
	RoutingEntries []RoutingEntry
	StatusEntries  []StatusEntry
}
