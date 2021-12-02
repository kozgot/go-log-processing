package processing

import (
	"encoding/json"
	"log"
	"strings"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/postprocessor/internal/utils"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

type EntryProcessor struct {
	eventsBySmcUID    map[string][]models.SmcEvent
	smcDataBySmcUID   map[string]models.SmcData
	smcUIDsByURL      map[string]string
	podUIDToSmcUID    map[string]string
	consumptionValues []models.ConsumtionValue
	indexValues       []models.IndexValue

	messageProducer rabbitmq.MessageProducer
	messageConsumer rabbitmq.MessageConsumer
}

func NewEntryProcessor(
	uploader rabbitmq.MessageProducer,
	messageConsumer rabbitmq.MessageConsumer,
) *EntryProcessor {
	eventsBySmcUID := make(map[string][]models.SmcEvent)
	smcDataBySmcUID := make(map[string]models.SmcData)
	smcUIDsByURL := make(map[string]string)
	podUIDToSmcUID := make(map[string]string)
	consumptionValues := []models.ConsumtionValue{}
	indexValues := []models.IndexValue{}

	result := EntryProcessor{
		eventsBySmcUID:    eventsBySmcUID,
		smcDataBySmcUID:   smcDataBySmcUID,
		smcUIDsByURL:      smcUIDsByURL,
		podUIDToSmcUID:    podUIDToSmcUID,
		consumptionValues: consumptionValues,
		indexValues:       indexValues,
		messageProducer:   uploader,
		messageConsumer:   messageConsumer,
	}

	return &result
}

// HandleEntries consumes entries from the provided MessageConsumer,
// and publishes them to a rabbitmq queue using the provided MessageProducer.
func (processor *EntryProcessor) HandleEntries() {
	msgs := processor.messageConsumer.ConsumeMessages()

	go func() {
		for d := range msgs {
			if strings.Contains(string(d.Body), "END") {
				log.Println(" [PROCESSOR] End of entries...")

				// Further processing to get consumption and index info.
				consumptionProcessor := NewConsumptionProcessor(
					processor.consumptionValues,
					processor.indexValues,
					processor.messageProducer,
				)
				consumptionProcessor.ProcessConsumptionAndIndexValues()

				log.Println(" [PROCESSOR] Done processing consumption data")

				// Acknowledge the message after it has been processed.
				err := d.Ack(false)
				utils.FailOnError(err, " [PROCESSOR] Could not acknowledge END message")

				// Clear previous processed data.
				processor.reset()
				continue
			}

			entry := deserializeParsedLogEntry(d.Body)
			processor.ProcessEntry(entry)

			// Acknowledge the message after it has been processed.
			err := d.Ack(false)
			utils.FailOnError(err,
				" [PROCESSOR] Could not acknowledge message with timestamp: "+entry.Timestamp.Format("2 Jan 2006 15:04:05"))
		}
	}()
}

// ProcessEntry processes the log entry received as a parameter.
func (processor *EntryProcessor) ProcessEntry(logEntry parsermodels.ParsedLogEntry) {
	var data *models.SmcData
	var event *models.SmcEvent
	var consumption *models.ConsumtionValue
	var indexvalue *models.IndexValue
	switch logEntry.Level {
	case "INFO":
		infoProcessor := InfoProcessor{
			PodUIDToSmcUID: processor.podUIDToSmcUID,
		}
		data, event, consumption, indexvalue = infoProcessor.ProcessInfoEntry(logEntry)

		if event != nil && event.EventType == models.ConnectionAttempt {
			// This is the only entry where the URL and SMC UID parameters are given at the same time.
			URL := data.Address.URL
			UID := data.SmcUID

			// Save it in the dictionary so we can use it later in the processing logic.
			_, ok := processor.smcUIDsByURL[URL]
			if !ok {
				processor.smcUIDsByURL[URL] = UID
			}
		}

		if indexvalue != nil {
			processor.indexValues = append(processor.indexValues, *indexvalue)
		}
		if consumption != nil {
			processor.consumptionValues = append(processor.consumptionValues, *consumption)
		}

	case "WARN":
		warningProcessor := WarningProcessor{}
		data, event = warningProcessor.ProcessWarn(logEntry)

	case "WARNING":
		warningProcessor := WarningProcessor{}
		data, event = warningProcessor.ProcessWarning(logEntry)

	case "ERROR":
		errorProcessor := ErrorProcessor{}
		data, event = errorProcessor.ProcessError(logEntry)

	default:
		log.Printf(" [PROCESSOR] Unknown log level %s", logEntry.Level)
	}

	processor.registerEvent(event, data)
	processor.updateSmcData(data)
}

func initArrayIfNeeded(eventsBySmcUID map[string][]models.SmcEvent, uid string) {
	_, ok := eventsBySmcUID[uid]
	if !ok {
		eventsBySmcUID[uid] = []models.SmcEvent{}
	}
}

func (processor *EntryProcessor) registerEvent(event *models.SmcEvent, data *models.SmcData) {
	if data == nil {
		return
	}

	if event == nil {
		return
	}

	smcUID := data.SmcUID

	// If only a URL is provided, use that to get the SMC UID.
	if smcUID == "" && data.Address.URL != "" {
		smcUID = processor.smcUIDsByURL[data.Address.URL]
	}

	if event.SmcUID == "" {
		event.SmcUID = smcUID
	}

	// Append the event to the corresponding array.
	initArrayIfNeeded(processor.eventsBySmcUID, smcUID)
	processor.eventsBySmcUID[smcUID] = append(processor.eventsBySmcUID[smcUID], *event)

	// send to ES
	processor.messageProducer.PublishEvent(*event)
}

func (processor *EntryProcessor) updateSmcData(data *models.SmcData) {
	if data == nil {
		return
	}

	// We have to find the smc by URL, UID is not provided.
	if data.SmcUID == "" && data.Address.URL != "" {
		smcUID := processor.smcUIDsByURL[data.Address.URL]
		smcData, ok := processor.smcDataBySmcUID[smcUID]
		if !ok {
			// add new value
			processor.smcDataBySmcUID[smcUID] = *data
		} else {
			newSmcData := updateChangedProperties(smcData, *data)

			// replace existing value
			processor.smcDataBySmcUID[smcUID] = newSmcData
		}
	}

	// UID is provided.
	if data.SmcUID != "" {
		smcData, ok := processor.smcDataBySmcUID[data.SmcUID]
		if !ok {
			// add new value
			processor.smcDataBySmcUID[data.SmcUID] = *data
		} else {
			newSmcData := updateChangedProperties(smcData, *data)

			// replace existing value
			processor.smcDataBySmcUID[data.SmcUID] = newSmcData
		}
	}
}

func updateChangedProperties(existingSmcData models.SmcData, newSmcData models.SmcData) models.SmcData {
	result := existingSmcData

	// Update address if there are valid changes
	result.Address = updateAddresIfNeeded(result.Address, newSmcData.Address)

	// Update pod list if there is a difference
	if len(result.Pods) < len(newSmcData.Pods) {
		for i := 0; i < len(newSmcData.Pods); i++ {
			pod := newSmcData.Pods[i]
			if !result.ContainsPod(pod) {
				result.Pods = append(result.Pods, pod)
			}
		}
	}

	// Update other properties if there are valid changes
	if result.CustomerSerialNumber != newSmcData.CustomerSerialNumber && newSmcData.CustomerSerialNumber != "" {
		result.CustomerSerialNumber = newSmcData.CustomerSerialNumber
	}

	if result.LastJoiningDate != newSmcData.LastJoiningDate && newSmcData.LastJoiningDate.Year() > 1500 {
		result.LastJoiningDate = newSmcData.LastJoiningDate
	}

	if result.LastSuccesfulDlmsResponse != newSmcData.LastSuccesfulDlmsResponse &&
		newSmcData.LastSuccesfulDlmsResponse.Year() > 1500 { // Also check date validity
		result.LastSuccesfulDlmsResponse = newSmcData.LastSuccesfulDlmsResponse
	}

	if result.SmcUID != newSmcData.SmcUID && newSmcData.SmcUID != "" {
		result.SmcUID = newSmcData.SmcUID
	}

	return result
}

func updateAddresIfNeeded(oldAddress models.AddressDetails, newAddress models.AddressDetails) models.AddressDetails {
	result := oldAddress

	// Update address if there are valid changes
	if oldAddress != newAddress {
		// logical address can be empty when invalidation logical address,
		// but in that case, the other properties (except for URL) cannot be empty (at least check the physical address)
		if result.LogicalAddress != newAddress.LogicalAddress &&
			(newAddress.LogicalAddress != "" || newAddress.PhysicalAddress != "") {
			result.LogicalAddress = newAddress.LogicalAddress
		}

		if result.PhysicalAddress != newAddress.PhysicalAddress && newAddress.PhysicalAddress != "" {
			result.PhysicalAddress = newAddress.PhysicalAddress
		}

		if result.ShortAddress != newAddress.ShortAddress && newAddress.ShortAddress != 0 {
			result.ShortAddress = newAddress.ShortAddress
		}

		if result.URL != newAddress.URL && newAddress.URL != "" {
			result.URL = newAddress.URL
		}
	}

	return result
}

func (processor *EntryProcessor) reset() {
	for k := range processor.eventsBySmcUID {
		delete(processor.eventsBySmcUID, k)
	}

	for k := range processor.smcDataBySmcUID {
		delete(processor.smcDataBySmcUID, k)
	}

	for k := range processor.smcUIDsByURL {
		delete(processor.smcUIDsByURL, k)
	}

	for k := range processor.podUIDToSmcUID {
		delete(processor.podUIDToSmcUID, k)
	}

	processor.consumptionValues = []models.ConsumtionValue{}
	processor.indexValues = []models.IndexValue{}
}

func deserializeParsedLogEntry(bytes []byte) parsermodels.ParsedLogEntry {
	var parsedEntry parsermodels.ParsedLogEntry
	err := json.Unmarshal(bytes, &parsedEntry)
	utils.FailOnError(err, "Failed to unmarshal log entry")
	return parsedEntry
}
