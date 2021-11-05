package service

import (
	"encoding/json"
	"fmt"
	"strings"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/internal/processing"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"
	"github.com/streadway/amqp"
)

// MessageConsumer encapsulates messages needed to consume messages.
type MessageConsumer interface {
	Consume() (<-chan amqp.Delivery, error)
}

// HandleEntries consumes entries from the provided MessageConsumer,
// and uploads them to ES using the provided ESUploader.
func HandleEntries(messageConsumer MessageConsumer, esUploader processing.ESUploader) {
	processor := processing.InitEntryProcessor(esUploader)

	msgs, err := messageConsumer.Consume()
	utils.FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			if strings.Contains(string(d.Body), "END") {
				fmt.Println("End of entries...")

				// Further processing to get consumption and index info.
				processor.ProcessConsumptionAndIndexValues()

				// Acknowledge the message after it has been processed.
				err := d.Ack(false)
				utils.FailOnError(err, "Could not acknowledge END message")
				continue
			}

			entry := deserializeMessage(d.Body)
			processor.ProcessEntry(entry)

			// Acknowledge the message after it has been processed.
			err := d.Ack(false)
			utils.FailOnError(err,
				"Could not acknowledge message with timestamp: "+entry.Timestamp.Format("2 Jan 2006 15:04:05"))
		}
	}()
}

func deserializeMessage(message []byte) parsermodels.ParsedLogEntry {
	var data parsermodels.ParsedLogEntry
	err := json.Unmarshal(message, &data)
	utils.FailOnError(err, "Failed to unmarshal log entry")
	return data
}
