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

type errorProcessorTest struct {
	inputEntry       parsermodels.ParsedLogEntry
	expectedSmcData  *models.SmcData
	expectedSmcEvent *models.SmcEvent
}

func TestProcessErrorEntry(t *testing.T) {
	errorTests := []errorProcessorTest{
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp:   time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
				Level:       "ERROR",
				ErrorParams: nil,
			},
			expectedSmcData:  nil,
			expectedSmcEvent: nil,
		},
		{
			inputEntry: parsermodels.ParsedLogEntry{
				Timestamp: time.Date(2020, time.June, 10, 10, 26, 37, 0, time.UTC),
				Level:     "ERROR",
				ErrorParams: &parsermodels.ErrorParams{
					Source:      "dc18-smc32",
					Message:     "DLMS error",
					Severity:    3,
					Description: "n/a",
					ErrorCode:   241,
				},
			},
			expectedSmcData: &models.SmcData{
				SmcUID: "dc18-smc32",
			},
			expectedSmcEvent: &models.SmcEvent{
				Time:            time.Date(2020, time.June, 10, 10, 26, 37, 0, time.UTC),
				EventType:       models.DLMSError,
				EventTypeString: models.EventTypeToString(models.DLMSError),
				Label:           "Error type DLMS error" + ", severity: 3",
				SmcUID:          "dc18-smc32",
				SMC: models.SmcData{
					SmcUID: "dc18-smc32",
				},
			},
		},
	}

	for i, test := range errorTests {
		errorProcessor := processing.ErrorProcessor{}
		data, event := errorProcessor.ProcessError(test.inputEntry)

		testutils.AssertEqualSmcData(data, test.expectedSmcData, t, i)
		testutils.AssertEqualSmcEvent(event, test.expectedSmcEvent, t, i)
	}

	log.Printf("Successfully run %d tests", len(errorTests))
}
