package processing

import (
	"fmt"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/streadway/amqp"
)

// Process processes the log entry received as a parameter.
func Process(logEntry parsermodels.ParsedLine, channel *amqp.Channel) models.ProcessedEntries {
	entriesBySmcUID := make(map[string][]models.SmcEntry)
	routingEntries := []models.RoutingEntry{}
	statusEntries := []models.StatusEntry{}

	result := models.ProcessedEntries{}

	switch logEntry.Level {
	case "INFO":
		smcEntry, routingEntry, statusEntry := ProcessInfo(logEntry)
		if smcEntry != nil && smcEntry.UID != "" {
			uid := smcEntry.UID
			initArrayIfNeeded(entriesBySmcUID, uid)

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
			saveToDb(*smcEntry, channel)
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
			initArrayIfNeeded(entriesBySmcUID, uid)

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
			saveToDb(*smcEntry, channel)
		}
	case "WARNING":
		smcEntry := ProcessWarning(logEntry)
		if smcEntry != nil {
			uid := smcEntry.UID
			initArrayIfNeeded(entriesBySmcUID, uid)

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
			saveToDb(*smcEntry, channel)
		}
	case "ERROR":
		smcEntry := ProcessError(logEntry)

		if smcEntry != nil {
			uid := smcEntry.UID
			initArrayIfNeeded(entriesBySmcUID, uid)

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
			saveToDb(*smcEntry, channel)
		}
	default:
		fmt.Printf("Unknown log level %s", logEntry.Level)
	}

	result.RoutingEntries = routingEntries
	result.StatusEntries = statusEntries
	result.SmcEntries = entriesBySmcUID

	return result
}

func saveToDb(entry models.SmcEntry, channel *amqp.Channel) {
	rabbitmq.SendEntryToElasticUploader(entry, channel, "smc")
}

func initArrayIfNeeded(entriesBySmcUID map[string][]models.SmcEntry, uid string) {
	_, ok := entriesBySmcUID[uid]
	if !ok {
		entriesBySmcUID[uid] = []models.SmcEntry{}
	}
}
