package processing

import (
	"fmt"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

type EntryProcessor struct {
	eventsBySmcUID    map[string][]models.SmcEvent
	smcDataBySmcUID   map[string]models.SmcData
	smcUIDsByURL      map[string]string
	podUIDToSmcUID    map[string]string
	consumptionValues []models.ConsumtionValue
	indexValues       []models.IndexValue

	esUploader *rabbitmq.EsUploadSender
}

func InitEntryProcessor(uploader *rabbitmq.EsUploadSender) *EntryProcessor {
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
		esUploader:        uploader,
	}

	return &result
}

// Process processes the log entry received as a parameter.
func (processor *EntryProcessor) Process(
	logEntry parsermodels.ParsedLogEntry,
	esIndexName string) {
	var data *models.SmcData
	var event *models.SmcEvent
	var consumption *models.ConsumtionValue
	var index *models.IndexValue
	switch logEntry.Level {
	case "INFO":
		data, event, consumption, index = ProcessInfoEntry(logEntry, processor.podUIDToSmcUID)

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

		if index != nil {
			processor.indexValues = append(processor.indexValues, *index)
		}
		if consumption != nil {
			processor.consumptionValues = append(processor.consumptionValues, *consumption)
		}

	case "WARN":
		data, event = ProcessWarn(logEntry)

	case "WARNING":
		data, event = ProcessWarning(logEntry)

	case "ERROR":
		data, event = ProcessError(logEntry)

	default:
		fmt.Printf("Unknown log level %s", logEntry.Level)
	}

	processor.registerEvent(event, data, esIndexName)
	processor.updateSmcData(data)
}

func initArrayIfNeeded(eventsBySmcUID map[string][]models.SmcEvent, uid string) {
	_, ok := eventsBySmcUID[uid]
	if !ok {
		eventsBySmcUID[uid] = []models.SmcEvent{}
	}
}

func (processor *EntryProcessor) registerEvent(
	event *models.SmcEvent,
	data *models.SmcData,
	esIndexName string) {
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
	processor.esUploader.SendEventToElasticUploader(*event, esIndexName)
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
			if !containsPod(pod, result.Pods) {
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

func containsPod(pod models.Pod, list []models.Pod) bool {
	for _, p := range list {
		if p.UID == pod.UID {
			return true
		}
	}
	return false
}
