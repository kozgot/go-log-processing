package processorunittests

import (
	"log"
	"testing"
	"time"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/internal/processing"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

type infoProcessorTest struct {
	inputEntry          parsermodels.ParsedLogEntry
	expectedSmcData     *models.SmcData
	expectedSmcEvent    *models.SmcEvent
	expectedConsumption *models.ConsumtionValue
	expectedIndex       *models.IndexValue
}

func TestProcessInfoEntry(t *testing.T) {
	podUIDToSmcUID := make(map[string]string)
	infoProcessorTests := []infoProcessorTest{
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp:  time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
				Level:      "INFO",
				InfoParams: nil,
			},
			expectedSmcData:     nil,
			expectedSmcEvent:    nil,
			expectedConsumption: nil,
			expectedIndex:       nil,
		},
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp: time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
				Level:     "INFO",
				InfoParams: &parsermodels.InfoParams{
					EntryType: parsermodels.UnknownEntryType,
				},
			},
			expectedSmcData:     nil,
			expectedSmcEvent:    nil,
			expectedConsumption: nil,
			expectedIndex:       nil,
		},
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp: time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
				Level:     "INFO",
				InfoParams: &parsermodels.InfoParams{
					EntryType: parsermodels.Routing,
					RoutingMessage: &parsermodels.RoutingTableParams{
						Address:        "0x0014",
						NextHopAddress: "0x0014",
						RouteCost:      11,
						HopCount:       0,
						WeakLink:       1,
						ValidTimeMins:  240,
					},
				},
			},
			expectedSmcData:     nil,
			expectedSmcEvent:    nil,
			expectedConsumption: nil,
			expectedIndex:       nil,
		},
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp: time.Date(2020, time.June, 10, 10, 2, 37, 0, time.UTC),
				Level:     "INFO",
				InfoParams: &parsermodels.InfoParams{
					EntryType: parsermodels.DCMessage,
					DCMessage: &parsermodels.DCMessageParams{
						IsInComing:       true,
						SourceOrDestName: "SMC",
						MessageType:      parsermodels.Consumption,
						Payload: &parsermodels.DcMessagePayload{
							TimeRange: &parsermodels.TimeRange{
								From: time.Date(2020, time.June, 10, 8, 0, 0, 0, time.UTC),
								To:   time.Date(2020, time.June, 10, 8, 10, 0, 0, time.UTC),
							},
							Value:          5,
							ServiceLevelID: 9,
						},
					},
				},
			},
			expectedSmcData:  nil,
			expectedSmcEvent: nil,
			expectedConsumption: &models.ConsumtionValue{
				ReceiveTime:  time.Date(2020, time.June, 10, 10, 2, 37, 0, time.UTC),
				StartTime:    time.Date(2020, time.June, 10, 8, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2020, time.June, 10, 8, 10, 0, 0, time.UTC),
				Value:        5,
				ServiceLevel: 9,
			},
			expectedIndex: nil,
		},
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp: time.Date(2020, time.June, 10, 9, 24, 18, 0, time.UTC),
				Level:     "INFO",
				InfoParams: &parsermodels.InfoParams{
					EntryType: parsermodels.SMCJoin,
					JoinMessage: &parsermodels.SmcJoinMessageParams{
						Ok:       true,
						Response: "Confirmed",
						JoinType: "LBD",
						SmcAddress: parsermodels.SmcAddressParams{
							SmcUID:          "dc18-smc24",
							PhysicalAddress: "EEBEDDFFFE62106D",
							LogicalAddress:  "FE80::4021:FF:FE00:001f:61616",
							ShortAddress:    31,
							LastJoiningDate: time.Date(2020, time.June, 10, 9, 24, 16, 0, time.UTC),
						},
					},
				},
			},
			expectedSmcData: &models.SmcData{
				Address: models.AddressDetails{
					PhysicalAddress: "EEBEDDFFFE62106D",
					LogicalAddress:  "FE80::4021:FF:FE00:001f:61616",
					ShortAddress:    31,
					URL:             "",
				},
				SmcUID:               "dc18-smc24",
				CustomerSerialNumber: "",
				LastJoiningDate:      time.Date(2020, time.June, 10, 9, 24, 16, 0, time.UTC),
			},
			expectedSmcEvent: &models.SmcEvent{
				Time:            time.Date(2020, time.June, 10, 9, 24, 18, 0, time.UTC),
				EventType:       models.SmcJoined,
				EventTypeString: models.EventTypeToString(models.SmcJoined),
				Label:           "Smc " + "dc18-smc24" + " has joined",
				SmcUID:          "dc18-smc24",
				DataPayload: models.SmcData{
					Address: models.AddressDetails{
						PhysicalAddress: "EEBEDDFFFE62106D",
						LogicalAddress:  "FE80::4021:FF:FE00:001f:61616",
						ShortAddress:    31,
						URL:             "",
					},
					SmcUID:               "dc18-smc24",
					CustomerSerialNumber: "",
					LastJoiningDate:      time.Date(2020, time.June, 10, 9, 24, 16, 0, time.UTC),
				},
			},
			expectedConsumption: nil,
			expectedIndex:       nil,
		},
	}

	for i, test := range infoProcessorTests {
		infoProcessor := processing.InfoEntryProcessor{
			PodUIDToSmcUID: podUIDToSmcUID,
		}

		smcData, event, consumption, index := infoProcessor.ProcessInfoEntry(test.inputEntry)

		assertEqualSmcData(smcData, test.expectedSmcData, t, i)
		assertEqualSmcEvent(event, test.expectedSmcEvent, t, i)
		assertEqualConsumption(consumption, test.expectedConsumption, t, i)
		assertEqualIndex(index, test.expectedIndex, t, i)
	}

	log.Printf("Successfully run %d test cases.", len(infoProcessorTests))
}

func assertEqualSmcData(actual *models.SmcData, expected *models.SmcData, t *testing.T, testIndex int) {
	if actual == nil && expected != nil ||
		actual != nil && expected == nil ||
		actual != nil && !actual.Equals(*expected) {
		t.Fatalf("Expected SMC Data does not match actual SMC Data in test no. %d.", testIndex+1)
	}
}

func assertEqualSmcEvent(actual *models.SmcEvent, expected *models.SmcEvent, t *testing.T, testIndex int) {
	if actual == nil && expected != nil ||
		actual != nil && expected == nil ||
		actual != nil && !actual.Equals(*expected) {
		t.Fatalf("Expected SMC Event does not match actual SMC Event in test no. %d.", testIndex+1)
	}
}

func assertEqualConsumption(
	actual *models.ConsumtionValue,
	expected *models.ConsumtionValue,
	t *testing.T,
	testIndex int) {
	if actual == nil && expected != nil ||
		actual != nil && expected == nil ||
		actual != nil && !actual.Equals(*expected) {
		t.Fatalf("Expected Consumtion Value does not match actual Consumtion Value in test no. %d.", testIndex+1)
	}
}

func assertEqualIndex(actual *models.IndexValue, expected *models.IndexValue, t *testing.T, testIndex int) {
	if actual == nil && expected != nil ||
		actual != nil && expected == nil ||
		actual != nil && !actual.Equals(*expected) {
		t.Fatalf("Expected Index Value does not match actual Index Value in test no. %d.", testIndex+1)
	}
}
