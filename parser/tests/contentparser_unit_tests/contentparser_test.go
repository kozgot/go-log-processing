package contentparserunittests

import (
	"reflect"
	"testing"
	"time"

	"github.com/kozgot/go-log-processing/parser/internal/contentparser"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func TestParseContents(t *testing.T) {
	errorEntry := models.EntryWithLevelAndTimestamp{
		Level:     "ERROR",
		Timestamp: time.Date(2020, time.June, 10, 14, 56, 3, 0, time.UTC),
		Rest: "error_code[241] message[DLMS error] severity[3] description[n/a] source[dc18-smc18]" +
			"  (smart_meter_cabinet.cc::129)",
	}

	warningEntry := models.EntryWithLevelAndTimestamp{
		Level:     "WARN",
		Timestamp: time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
		Rest:      "Timeout protocol[plc-udp] url[fe80::4021:ff:fe00:10:61616] (plc_bridge_connector.cc::227)",
	}

	infoConsumptionEntry := models.EntryWithLevelAndTimestamp{
		Level:     "INFO",
		Timestamp: time.Date(2020, time.June, 10, 14, 56, 19, 0, time.UTC),

		Rest: "<--[consumption]--(SMC) start[1591800000] end[1591800600] value[0] service_level_id[9]" +
			" (abstract_smart_meter_cabinet.h::314)",
	}

	expectedErrorEntry := models.ParsedLogEntry{
		Level:         "ERROR",
		Timestamp:     time.Date(2020, time.June, 10, 14, 56, 3, 0, time.UTC),
		InfoParams:    nil,
		WarningParams: nil,
		ErrorParams: &models.ErrorParams{
			ErrorCode:   241,
			Message:     "DLMS error",
			Severity:    3,
			Description: "n/a",
			Source:      "dc18-smc18",
		},
	}

	expectedWarningEntry := models.ParsedLogEntry{
		Level:       "WARN",
		Timestamp:   time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
		InfoParams:  nil,
		ErrorParams: nil,
		WarningParams: &models.WarningParams{
			WarningType: models.TimeoutWarning,
			TimeoutParams: &models.TimeOutParams{
				Protocol: "plc-udp",
				URL:      "fe80::4021:ff:fe00:10:61616",
			},
		},
	}

	expectedInfoConsumptionEntry := models.ParsedLogEntry{
		Level:     "INFO",
		Timestamp: time.Date(2020, time.June, 10, 14, 56, 19, 0, time.UTC),
		InfoParams: &models.InfoParams{
			EntryType: models.DCMessage,
			DCMessage: &models.DCMessageParams{
				IsInComing:       true,
				SourceOrDestName: "SMC",
				MessageType:      models.Consumption,
				Payload: &models.DcMessagePayload{
					ServiceLevelID: 9,
					Value:          0,
					TimeRange: &models.TimeRange{
						From: time.Date(2020, time.June, 10, 14, 40, 0, 0, time.UTC), // 2020. June 10., Wednesday 14:40:00 = 1591800000
						To:   time.Date(2020, time.June, 10, 14, 50, 0, 0, time.UTC), // 2020. June 10., Wednesday 14:50:00 = 1591800600
					},
				},
			},
		},
	}

	parsedErrorEntry := contentparser.ParseEntryContents(errorEntry)
	parsedWarningEntry := contentparser.ParseEntryContents(warningEntry)
	parsedInfoConsumptionEntry := contentparser.ParseEntryContents(infoConsumptionEntry)

	if !reflect.DeepEqual(parsedErrorEntry, &expectedErrorEntry) {
		t.Fatal("Got error entry does not match expected error entry")
	}

	if !reflect.DeepEqual(parsedWarningEntry, &expectedWarningEntry) {
		t.Fatal("Got warning entry does not match expected warning entry")
	}

	if !reflect.DeepEqual(parsedInfoConsumptionEntry, &expectedInfoConsumptionEntry) {
		t.Fatal("Got info entry does not match expected info entry")
	}
}
