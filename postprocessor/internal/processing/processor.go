package processing

import (
	"fmt"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

func Process(logEntry parsermodels.ParsedLine) models.ProcessedEntries {
	entriesBySmcUID := make(map[string][]models.SmcEntry)
	routingEntries := []models.RoutingEntry{}
	statusEntries := []models.StatusEntry{}

	result := models.ProcessedEntries{}

	switch logEntry.Level {
	case "INFO":
		smcEntry, routingEntry, statusEntry := ProcessInfo(logEntry)
		if smcEntry != nil {
			uid := smcEntry.UID
			_, ok := entriesBySmcUID[uid]
			if !ok {
				entriesBySmcUID[uid] = []models.SmcEntry{}
			}

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
		}

		if routingEntry != nil {
			routingEntries = append(routingEntries, *routingEntry)
		}
		if statusEntry != nil {
			statusEntries = append(statusEntries, *statusEntry)
		}
	case "WARN":
		smcEntry := ProcessWarn(logEntry)
		if smcEntry != nil {
			uid := smcEntry.UID
			_, ok := entriesBySmcUID[uid]
			if !ok {
				entriesBySmcUID[uid] = []models.SmcEntry{}
			}

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
		}
	case "WARNING":
		smcEntry := ProcessWarning(logEntry)
		if smcEntry != nil {
			uid := smcEntry.UID
			_, ok := entriesBySmcUID[uid]
			if !ok {
				entriesBySmcUID[uid] = []models.SmcEntry{}
			}

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
		}
	case "ERROR":
		smcEntry := ProcessError(logEntry)

		if smcEntry != nil {
			uid := smcEntry.UID
			_, ok := entriesBySmcUID[uid]
			if !ok {
				entriesBySmcUID[uid] = []models.SmcEntry{}
			}

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
		}
	default:
		fmt.Printf("Unknown log level %s", logEntry.Level)
	}

	result.RoutingEntries = routingEntries
	result.StatusEntries = statusEntries
	result.SmcEntries = entriesBySmcUID

	return result
}
