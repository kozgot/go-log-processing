package processing

import (
	"fmt"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/streadway/amqp"
)

// Process processes the log entry received as a parameter.
func Process(logEntry parsermodels.ParsedLogEntry, channel *amqp.Channel) {
	/*
		entriesBySmcUID := make(map[string][]models.SmcEntry)
		routingEntries := []models.RoutingEntry{}
		statusEntries := []models.StatusEntry{}
	*/

	eventsBySmcUID := make(map[string][]models.SmcEvent)
	smcDataBySmcUID := make(map[string]models.SmcData)
	smcUIDsByURL := make(map[string]string)

	switch logEntry.Level {
	case "INFO":
		data, event := ProcessInfoEntry(logEntry)
		// todo
		if event != nil && event.EventType == models.ConnectionAttempt {
			// This is the only entry where the URL and SMC UID parameters are given at the same time.
			URL := data.Address.URL
			UID := data.SmcUID

			// Put it in the dictionary
			_, ok := smcUIDsByURL[URL]
			if !ok {
				smcUIDsByURL[URL] = UID
			}
		}

		registerEvent(eventsBySmcUID, smcUIDsByURL, event, data)
		updateSmc(smcDataBySmcUID, smcUIDsByURL, data)

	case "WARN":
		data, event := ProcessWarn(logEntry)
		// todo
		registerEvent(eventsBySmcUID, smcUIDsByURL, event, data)
		updateSmc(smcDataBySmcUID, smcUIDsByURL, data)

	case "WARNING":
		data, event := ProcessWarning(logEntry)
		// todo
		registerEvent(eventsBySmcUID, smcUIDsByURL, event, data)
		updateSmc(smcDataBySmcUID, smcUIDsByURL, data)

	case "ERROR":
		data, event := ProcessError(logEntry)
		// todo
		registerEvent(eventsBySmcUID, smcUIDsByURL, event, data)
		updateSmc(smcDataBySmcUID, smcUIDsByURL, data)

	default:
		fmt.Printf("Unknown log level %s", logEntry.Level)
	}
}

/*
func saveToDb(entry models.SmcEntry, channel *amqp.Channel) {
	rabbitmq.SendEntryToElasticUploader(entry, channel, "smc")
}
*/

/*
func initArrayIfNeeded(eventsBySmcUID map[string][]models.SmcEvent, uid string) {
	_, ok := eventsBySmcUID[uid]
	if !ok {
		eventsBySmcUID[uid] = []models.SmcEvent{}
	}
}
*/

func registerEvent(eventsBySmcUID map[string][]models.SmcEvent, smcUIDsByURL map[string]string, event *models.SmcEvent, data *models.SmcData) {
	// todo
}

func updateSmc(smcDataBySmcUID map[string]models.SmcData, smcUIDsByURL map[string]string, data *models.SmcData) {
	// todo
	if data != nil {
		return
	}

	if data.SmcUID == "" && data.Address.URL != "" {
		smcUID := smcUIDsByURL[data.Address.URL]
		// todo update data...
		smcData := smcDataBySmcUID[smcUID]
		fmt.Println(smcData)
	}

	if data.SmcUID != "" {
		// todo: check what needs to be updated
		smcData := smcDataBySmcUID[data.SmcUID]
		fmt.Println(smcData)
	}
}
