package models

// ProcessedEntryData contains relevant data extracted from a processed log entry.
type ProcessedEntryData struct {
	SmcData         *SmcData
	SmcEvent        *SmcEvent
	ConsumtionValue *ConsumtionValue
	IndexValue      *IndexValue
}
