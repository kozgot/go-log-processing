package processing

import (
	"fmt"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/streadway/amqp"
)

// Process processes the log entry received as a parameter.
func Process(logEntry parsermodels.ParsedLogEntry, channel *amqp.Channel) {
	/*
		entriesBySmcUID := make(map[string][]models.SmcEntry)
		routingEntries := []models.RoutingEntry{}
		statusEntries := []models.StatusEntry{}
	*/

	switch logEntry.Level {
	case "INFO":
		data, event := ProcessInfoEntry(logEntry)
		// todo
		if data != nil {
			fmt.Println(data)
		}

		if event != nil {
			fmt.Println(event)
		}
	case "WARN":
		data, event := ProcessWarn(logEntry)
		// todo
		if data != nil {
			fmt.Println(data)
		}

		if event != nil {
			fmt.Println(event)
		}

	case "WARNING":
		data, event := ProcessWarning(logEntry)
		// todo
		if data != nil {
			fmt.Println(data)
		}

		if event != nil {
			fmt.Println(event)
		}

	case "ERROR":
		data, event := ProcessError(logEntry)
		// todo
		if data != nil {
			fmt.Println(data)
		}

		if event != nil {
			fmt.Println(event)
		}

	default:
		fmt.Printf("Unknown log level %s", logEntry.Level)
	}
}

/*
func saveToDb(entry models.SmcEntry, channel *amqp.Channel) {
	rabbitmq.SendEntryToElasticUploader(entry, channel, "smc")
}

func initArrayIfNeeded(entriesBySmcUID map[string][]models.SmcEntry, uid string) {
	_, ok := entriesBySmcUID[uid]
	if !ok {
		entriesBySmcUID[uid] = []models.SmcEntry{}
	}
}
*/
