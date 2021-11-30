package processorunittests

import (
	"log"
	"testing"
	"time"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/internal/processing"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/tests/testutils"
)

type warningProcessorTest struct {
	inputEntry       parsermodels.ParsedLogEntry
	expectedSmcData  *models.SmcData
	expectedSmcEvent *models.SmcEvent
}

func TestProcessWarnEntry(t *testing.T) {
	warnTests := []warningProcessorTest{
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp:     time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
				Level:         "WARN",
				WarningParams: nil,
			},
			expectedSmcData:  nil,
			expectedSmcEvent: nil,
		},
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp: time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
				Level:     "WARN",
				WarningParams: &parsermodels.WarningParams{
					TimeoutParams: &parsermodels.TimeOutParams{
						Protocol: "plc-udp",
						URL:      "fe80::4021:ff:fe00:a:61616",
					},
				},
			},
			expectedSmcData: &models.SmcData{
				Address: models.AddressDetails{
					URL: "fe80::4021:ff:fe00:a:61616",
				},
			},
			expectedSmcEvent: &models.SmcEvent{
				Time:            time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
				EventType:       models.TimeoutWarning,
				EventTypeString: models.EventTypeToString(models.TimeoutWarning),
				Label:           "Timeout for URL " + "fe80::4021:ff:fe00:a:61616",
				SMC: models.SmcData{
					Address: models.AddressDetails{
						URL: "fe80::4021:ff:fe00:a:61616",
					},
				},
			},
		},
	}

	for i, test := range warnTests {
		warningProcessor := processing.WarningProcessor{}
		data, event := warningProcessor.ProcessWarn(test.inputEntry)

		testutils.AssertEqualSmcData(data, test.expectedSmcData, t, i)
		testutils.AssertEqualSmcEvent(event, test.expectedSmcEvent, t, i)
	}

	log.Printf("Successfully run %d tests.", len(warnTests))
}

func TestProcessWarniongEntry(t *testing.T) {
	warningTests := []warningProcessorTest{
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp:     time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
				Level:         "WARNING",
				WarningParams: nil,
			},
			expectedSmcData:  nil,
			expectedSmcEvent: nil,
		},
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp: time.Date(2020, time.June, 10, 9, 26, 53, 0, time.UTC),
				Level:     "WARNING",
				WarningParams: &parsermodels.WarningParams{
					JoinMessageParams: &parsermodels.SmcJoinMessageParams{
						Ok:       false,
						Response: "Rejected",
						JoinType: "LBD",
						SmcAddress: parsermodels.SmcAddressParams{
							SmcUID:          "dc18-smc36",
							PhysicalAddress: "EEBEDDFFFE621128",
							LogicalAddress:  "FE80::4021:FF:FE00:0008:61616",
							ShortAddress:    8,
							LastJoiningDate: time.Date(2020, time.June, 10, 9, 26, 43, 0, time.UTC),
						},
					},
				},
			},
			expectedSmcData: &models.SmcData{
				Address: models.AddressDetails{
					PhysicalAddress: "EEBEDDFFFE621128",
					LogicalAddress:  "FE80::4021:FF:FE00:0008:61616",
					ShortAddress:    8,
				},
				SmcUID: "dc18-smc36",
			},
			expectedSmcEvent: &models.SmcEvent{
				Time:            time.Date(2020, time.June, 10, 9, 26, 53, 0, time.UTC),
				EventType:       models.JoinRejectedWarning,
				EventTypeString: models.EventTypeToString(models.JoinRejectedWarning),
				Label:           "SMC join rejected for " + "dc18-smc36",
				SmcUID:          "dc18-smc36",
				SMC: models.SmcData{
					Address: models.AddressDetails{
						PhysicalAddress: "EEBEDDFFFE621128",
						LogicalAddress:  "FE80::4021:FF:FE00:0008:61616",
						ShortAddress:    8,
					},
					SmcUID: "dc18-smc36",
				},
			},
		},
	}

	for i, test := range warningTests {
		warningProcessor := processing.WarningProcessor{}
		data, event := warningProcessor.ProcessWarning(test.inputEntry)

		testutils.AssertEqualSmcData(data, test.expectedSmcData, t, i)
		testutils.AssertEqualSmcEvent(event, test.expectedSmcEvent, t, i)
	}

	log.Printf("Successfully run %d tests.", len(warningTests))
}
