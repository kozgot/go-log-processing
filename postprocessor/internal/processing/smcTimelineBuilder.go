package processing

import (
	"fmt"
	"sort"
	"time"

	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/streadway/amqp"
)

func CreateSMCTimelines(
	eventsBySmcUID map[string][]models.SmcEvent,
	dataByUID map[string]models.SmcData,
	channel *amqp.Channel) {
	// Sort events by time.
	for k := range eventsBySmcUID {
		eventSlice := eventsBySmcUID[k]
		sort.Slice(eventSlice, func(i, j int) bool {
			return eventSlice[i].Time.Before(eventSlice[j].Time)
		})

		eventsBySmcUID[k] = eventSlice
	}

	// Build timelines from the events.
	smcTimelines := make(map[string]models.SmcTimeline)
	for k := range eventsBySmcUID {
		eventSlice := eventsBySmcUID[k]
		data := dataByUID[k]
		timeline := createTimelineForSMC(eventSlice)
		timeline.SmcData = data
		smcTimelines[k] = timeline
		saveSmcTimeline(timeline, channel)
	}
}

func saveSmcTimeline(timeline models.SmcTimeline, channel *amqp.Channel) {
	rabbitmq.SendTimelineToElasticUploader(timeline, channel, "timelines")
	fmt.Println("Sent timeline " + timeline.SmcData.SmcUID)
}

func createTimelineForSMC(smcEvents []models.SmcEvent) models.SmcTimeline {
	result := models.SmcTimeline{}
	sections := []models.TimelineSection{}
	currentSection := models.TimelineSection{}

	for _, event := range smcEvents {
		switch event.EventType {
		case models.NewSmc:
			currentSection = handleNewSmc(currentSection, event)

		case models.ConfigurationUpdated:
			currentSection, sections = handleConfigUpdated(currentSection, event, sections)

		case models.SmcJoined:
			currentSection, sections = handleJoined(currentSection, event, sections)

		case models.ConnectionAttempt:
			currentSection, sections = handleConnectionAttempt(currentSection, event, sections)

		case models.TimeoutWarning:
			currentSection, sections = handleTimeoutWarning(currentSection, event, sections)

		case models.DLMSError:
			currentSection, sections = handleError(currentSection, event, sections)

		case models.ConnectionReleased:
			currentSection, sections = handleConnectionReleased(currentSection, event, sections)

		case models.IndexCollectionStarted:
			currentSection, sections = handleIndexCollectionStarted(currentSection, event, sections)

		case models.IndexLowProfileGenericReceived:
			currentSection, sections = handleIndexProfileReceived(currentSection, event, sections)

		case models.IndexHighProfileGenericReceived:
			currentSection, sections = handleIndexProfileReceived(currentSection, event, sections)

		case models.SmcAddressInvalidated:
			break
		case models.IndexRead:
			break
		case models.JoinRejectedWarning:
			break
		case models.InitConnection:
			break
		case models.StartToConnect:
			break
		case models.SmcAddressUpdated:
			break
		case models.InternalDiagnostics:
			break
		case models.ConfigurationReadFromDB:
			break
		case models.PodConfiguration:
			break
		case models.DLMSLogsSent:
			break
		case models.StatisticsSent:
			break
		case models.UnknownEventType:
			break
		default:
			break
		}
	}

	result.Sections = sections
	if len(sections) > 0 {
		result.From = sections[0].From
		result.To = sections[len(sections)-1].To
	}

	return result
}

func handleNewSmc(currentSection models.TimelineSection, event models.SmcEvent) models.TimelineSection {
	result := currentSection
	if result.State == models.UnknownSmcState {
		result.From = event.Time
		result.State = models.New
	}
	return result
}

func handleConfigUpdated(
	currentSection models.TimelineSection,
	event models.SmcEvent,
	sections []models.TimelineSection) (models.TimelineSection, []models.TimelineSection) {
	result := currentSection
	// If the last joining date attribute has a valid value, that marks the date when the smc has joined.
	// check date validity
	if event.DataPayload.LastJoiningDate.After(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		if result.State != models.UnknownSmcState && result.State != models.Joined {
			result.To = event.DataPayload.LastJoiningDate
			sections = append(sections, result)
		}

		result = models.TimelineSection{From: event.DataPayload.LastJoiningDate, State: models.Joined}
	}

	return result, sections
}

func handleJoined(
	currentSection models.TimelineSection,
	event models.SmcEvent,
	sections []models.TimelineSection) (models.TimelineSection, []models.TimelineSection) {
	result := currentSection

	if result.State != models.UnknownSmcState && result.State != models.Joined {
		result.To = event.Time
		sections = append(sections, result)

		// reset current section
		result = models.TimelineSection{From: event.Time, State: models.Joined}
	}

	return result, sections
}

func handleConnectionAttempt(
	currentSection models.TimelineSection,
	event models.SmcEvent,
	sections []models.TimelineSection) (models.TimelineSection, []models.TimelineSection) {
	result := currentSection

	if result.State != models.UnknownSmcState && result.State != models.Connecting {
		result.To = event.Time
		sections = append(sections, result)

		// reset current section
		result = models.TimelineSection{From: event.Time, State: models.Connecting}
	}

	return result, sections
}

func handleTimeoutWarning(
	currentSection models.TimelineSection,
	event models.SmcEvent,
	sections []models.TimelineSection) (models.TimelineSection, []models.TimelineSection) {
	result := currentSection

	if result.State != models.UnknownSmcState && result.State != models.Error {
		// Finish and add previous section.
		result.To = event.Time
		sections = append(sections, result)

		// reset current section
		result = models.TimelineSection{From: event.Time, State: models.Error}
	}

	return result, sections
}

func handleError(
	currentSection models.TimelineSection,
	event models.SmcEvent,
	sections []models.TimelineSection) (models.TimelineSection, []models.TimelineSection) {
	result := currentSection

	if result.State != models.UnknownSmcState && result.State != models.Error {
		// Finish and add previous section.
		result.To = event.Time
		sections = append(sections, result)

		// reset current section
		result = models.TimelineSection{From: event.Time, State: models.Error}
	}

	return result, sections
}

func handleConnectionReleased(
	currentSection models.TimelineSection,
	event models.SmcEvent,
	sections []models.TimelineSection) (models.TimelineSection, []models.TimelineSection) {
	result := currentSection

	if result.State != models.UnknownSmcState {
		// Finish and add previous section.
		result.To = event.Time
		sections = append(sections, result)

		// reset current section
		result = models.TimelineSection{From: event.Time, State: models.Disconnected}
	}

	return result, sections
}

func handleIndexCollectionStarted(
	currentSection models.TimelineSection,
	event models.SmcEvent,
	sections []models.TimelineSection) (models.TimelineSection, []models.TimelineSection) {
	result := currentSection

	if result.State != models.UnknownSmcState {
		// Finish and add previous section.
		result.To = event.Time
		sections = append(sections, result)
	}

	// Start current section.
	result = models.TimelineSection{From: event.Time, State: models.CollectingIndex}

	return result, sections
}

func handleIndexProfileReceived(
	currentSection models.TimelineSection,
	event models.SmcEvent,
	sections []models.TimelineSection) (models.TimelineSection, []models.TimelineSection) {
	result := currentSection

	if result.State == models.Connecting {
		result.To = event.Time
		sections = append(sections, result)

		result = models.TimelineSection{From: event.Time, State: models.CollectingIndex}
	}

	return result, sections
}
