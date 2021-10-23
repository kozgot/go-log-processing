package processing

import (
	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

func processDCMessageEntry(logEntry parsermodels.ParsedLogEntry) models.ProcessedEntryData {
	messageType := logEntry.InfoParams.DCMessage.MessageType
	switch messageType {
	case parsermodels.Connect:
		data, event := processConnect(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.DLMSLogs:
		data, event := processDLMSLogsEntry(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.IndexHighProfileGeneric:
		data, event := processIndexHighProfileGeneric(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.IndexLowProfileGeneric:
		data, event := processIndexLowProfileGeneric(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.ReadIndexLowProfiles:
		data, event := processReadIndexLowProfiles(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.ReadIndexProfiles:
		data, event := processReadIndexProfiles(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.IndexReceived:
		indexValue := processIndexReceived(logEntry)
		// todo
		result := models.ProcessedEntryData{
			SmcData:         nil,
			SmcEvent:        nil,
			ConsumtionValue: nil,
			IndexValue:      indexValue,
		}
		return result

	case parsermodels.Consumption:
		consumptionValue := processConsumption(logEntry)
		// todo
		result := models.ProcessedEntryData{
			SmcData:         nil,
			SmcEvent:        nil,
			ConsumtionValue: consumptionValue,
			IndexValue:      nil,
		}
		return result

	case parsermodels.NewSmc:
		data, event := processNewSmc(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.MessageSentToSVI:
		data, event := processSVIMessage(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.PodConfig:
		data, event := processPodConfig(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.SmcConfig:
		data, event := processSmcConfig(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.SmcAddress:
		data, event := processSmcAddress(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.ServiceLevel:
		data, event := processServicelevelEntry(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.Settings:
		data, event := processSettings(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.Statistics:
		data, event := processStatistics(logEntry)
		result := models.ProcessedEntryData{
			SmcData:         data,
			SmcEvent:        event,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result

	case parsermodels.UnknownDCMessage:
		result := models.ProcessedEntryData{
			SmcData:         nil,
			SmcEvent:        nil,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result
	default:
		result := models.ProcessedEntryData{
			SmcData:         nil,
			SmcEvent:        nil,
			ConsumtionValue: nil,
			IndexValue:      nil,
		}
		return result
	}
}

func processIndexLowProfileGeneric(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// there are two more params (capture period and capture objects), but they are not really interesting for us here.
	smcUID := logEntry.InfoParams.DCMessage.Payload.SmcUID
	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.IndexLowProfileGenericReceived,
		Label:     "Index low profile generic from SMC",
		SmcUID:    smcUID,
	}

	data := models.SmcData{
		SmcUID: smcUID,
	}

	return &data, &event
}

func processIndexHighProfileGeneric(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// there are two more params (capture period and capture objects), but they are not really interesting for us here.
	smcUID := logEntry.InfoParams.DCMessage.Payload.SmcUID
	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.IndexHighProfileGenericReceived,
		Label:     "Index high profile generic from SMC",
		SmcUID:    smcUID,
	}

	data := models.SmcData{
		SmcUID: smcUID,
	}

	return &data, &event
}

func processReadIndexProfiles(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	smcUID := logEntry.InfoParams.DCMessage.Payload.ReadIndexProfilesEntryPayload.SmcUID
	data := models.SmcData{
		SmcUID: smcUID,
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.IndexRead,
		Label:     "Index read",
		SmcUID:    smcUID,
	}

	return &data, &event
}

func processReadIndexLowProfiles(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	smcUID := logEntry.InfoParams.DCMessage.Payload.ReadIndexLowProfilesEntryPayload.SmcUID
	to := logEntry.InfoParams.DCMessage.Payload.ReadIndexLowProfilesEntryPayload.To
	from := logEntry.InfoParams.DCMessage.Payload.ReadIndexLowProfilesEntryPayload.From

	data := models.SmcData{
		SmcUID: smcUID,
	}

	toDateString := to.Format("2 Jan 2006 15:04:05")
	fromDateString := from.Format("2 Jan 2006 15:04:05")

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.IndexCollectionStarted,
		Label:     "Index collection started from " + fromDateString + " to " + toDateString,
		SmcUID:    smcUID,
	}

	return &data, &event
}

// These entries have the same timestamp as the corresponding <--[read index profiles] (IndexRead) entries.
// If the changes in dex or consumption values are interesting, we can get them from these messages.
// The pod UID and serial number fields can be used to pair these entries to SMC-s.
func processIndexReceived(logEntry parsermodels.ParsedLogEntry) *models.IndexValue {
	result := models.IndexValue{
		ReceiveTime:   logEntry.Timestamp,
		PreviousTime:  logEntry.InfoParams.DCMessage.Payload.IndexPayload.PreviousTime,
		Time:          logEntry.InfoParams.DCMessage.Payload.Time,
		PreviousValue: logEntry.InfoParams.DCMessage.Payload.IndexPayload.PreviousValue,
		Value:         logEntry.InfoParams.DCMessage.Payload.Value,
		ServiceLevel:  logEntry.InfoParams.DCMessage.Payload.ServiceLevelID,
		PodUID:        logEntry.InfoParams.DCMessage.Payload.PodUID,
		SerialNumber:  logEntry.InfoParams.DCMessage.Payload.IndexPayload.SerialNumber,
	}

	return &result
}

// These entries have the same timestamp as the corresponding <--[read index profiles]
// and <--[index]--(SMC) (IndexRead, IndexReceived) entries.
// If the changes in dex or consumption values are interesting, we can get them from these messages.
// There are no UID fields so only the timestamp and start/end fields can help us pair them to the smc-s.
func processConsumption(logEntry parsermodels.ParsedLogEntry) *models.ConsumtionValue {
	result := models.ConsumtionValue{
		ReceiveTime:  logEntry.Timestamp,
		StartTime:    logEntry.InfoParams.DCMessage.Payload.TimeRange.From,
		EndTime:      logEntry.InfoParams.DCMessage.Payload.TimeRange.To,
		Value:        logEntry.InfoParams.DCMessage.Payload.Value,
		ServiceLevel: logEntry.InfoParams.DCMessage.Payload.ServiceLevelID,
	}
	return &result
}

func processNewSmc(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	data := models.SmcData{
		SmcUID: logEntry.InfoParams.DCMessage.Payload.SmcUID,
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.NewSmc,
		Label:     "New SMC, UID: " + data.SmcUID,
		SmcUID:    data.SmcUID,
	}

	return &data, &event
}

func processPodConfig(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	smcUID := logEntry.InfoParams.DCMessage.Payload.SmcUID
	poidUID := logEntry.InfoParams.DCMessage.Payload.PodUID
	podData := models.Pod{
		SmcUID:         smcUID,
		UID:            poidUID,
		PositionInSmc:  logEntry.InfoParams.DCMessage.Payload.PodConfigPayload.PositionInSmc,
		Phase:          logEntry.InfoParams.DCMessage.Payload.PodConfigPayload.Phase,
		ServiceLevelID: logEntry.InfoParams.DCMessage.Payload.ServiceLevelID,
		SerialNumber:   logEntry.InfoParams.DCMessage.Payload.PodConfigPayload.SerialNumber,
	}

	pods := []models.Pod{podData}

	data := models.SmcData{
		Pods:   pods,
		SmcUID: smcUID,
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.PodConfiguration,
		Label:     "Pod configuration read for pod " + poidUID,
		SmcUID:    smcUID,
	}
	return &data, &event
}

func processConnect(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	if logEntry.InfoParams.DCMessage.SourceOrDestName != "PLC" {
		// We only care about PLC connects
		return nil, nil
	}

	URL := logEntry.InfoParams.DCMessage.Payload.ConnectToPLCPayload.DestinationAddress
	address := models.AddressDetails{
		URL: URL,
	}
	data := models.SmcData{
		Address: address,
	}
	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.StartToConnect,
		Label:     "Trying to connect to " + URL + " ...",
	}
	return &data, &event
}

func processSmcAddress(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	address := models.AddressDetails{
		PhysicalAddress: logEntry.InfoParams.DCMessage.Payload.SmcAddressPayload.PhysicalAddress,
		LogicalAddress:  logEntry.InfoParams.DCMessage.Payload.SmcAddressPayload.LogicalAddress,
		ShortAddress:    logEntry.InfoParams.DCMessage.Payload.SmcAddressPayload.ShortAddress,
	}
	data := models.SmcData{
		LastJoiningDate: logEntry.InfoParams.DCMessage.Payload.SmcAddressPayload.LastJoiningDate,
		SmcUID:          logEntry.InfoParams.DCMessage.Payload.SmcUID,
		Address:         address,
	}

	label := "SMC address updated"

	if address.LogicalAddress == "" {
		label = "SMC logical address invalidated"
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.SmcAddressUpdated,
		Label:     label,
		SmcUID:    data.SmcUID,
	}

	return &data, &event
}

func processSmcConfig(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	address := models.AddressDetails{
		PhysicalAddress: logEntry.InfoParams.DCMessage.Payload.SmcConfigPayload.PhysicalAddress,
	}

	data := models.SmcData{
		SmcUID:                    logEntry.InfoParams.DCMessage.Payload.SmcUID,
		Address:                   address,
		LastSuccesfulDlmsResponse: logEntry.InfoParams.DCMessage.Payload.SmcConfigPayload.LastSuccessfulDlmsResponseDate,
		CustomerSerialNumber:      logEntry.InfoParams.DCMessage.Payload.SmcConfigPayload.CustomerSerialNumber,
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.ConfigurationChanged,
		Label:     "SMC Config updated",
		SmcUID:    data.SmcUID,
	}

	return &data, &event
}

func processStatistics(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	statisticsPayload := logEntry.InfoParams.DCMessage.Payload.StatisticsEntryPayload

	if statisticsPayload == nil {
		return nil, nil
	}

	smcUID := statisticsPayload.SourceID
	data := models.SmcData{
		SmcUID: smcUID,
	}

	event := models.SmcEvent{
		Time:      logEntry.Timestamp,
		EventType: models.StatisticsSent,
		Label:     "Statistics sent to SVI (" + statisticsPayload.Type + ")",
		SmcUID:    smcUID,
	}

	return &data, &event
}

func processSettings(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo: is this interesting?
	// --[settings]-->(DB) dc_uid[dc18] locality[Tanambao Daoud] region[Madagascar] timezone[Indian/Antananarivo]
	// global_ftp_address[sftp://sagemcom@172.30.31.20:firmwares] target_firmware_version[] index_collection[600]
	// data_publish[2400] last_server_communication_time[1591775824] dc_distro_target_firmware_version[]
	// last_dc_start_time[1591780709] frequency_band_changed[0] frequency_band_rollback_done[0]
	return nil, nil
}

func processDLMSLogsEntry(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo is this interesting??
	return nil, nil
}

func processServicelevelEntry(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo is this interesting??
	// <--[service_level]--(DB) service_level_id[1] meter_mode[2] start_hour_daily_cycle[20h]
	// load_shedding_daily_energy_budget[0] local_shedding_daily_energy_budget[0] max_active_power[0]
	// in_service[1] name[1. Suspension] hourly_energy_limits[[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]]
	// local_hourly_energy_limits[[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]]
	return nil, nil
}

func processSVIMessage(logEntry parsermodels.ParsedLogEntry) (*models.SmcData, *models.SmcEvent) {
	// todo is this interesting??
	// --[message]-->(SVI) current[9.8955] total[9.8955] url[tcp://172.30.31.20:9062] topic[dc.measurementQueue]

	return nil, nil
}
