package serviceunittests

import (
	"reflect"
	"testing"
	"time"

	"github.com/kozgot/go-log-processing/parser/internal/service"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func TestParseTimestamp(t *testing.T) {
	errorEntry := models.EntryWithLogLevel{
		Level: "ERROR",
		Rest: "Wed Jun 10 14:56:03 2020    : error_code[241] message[DLMS error] severity[3] " +
			"description[n/a] source[dc18-smc18]  (smart_meter_cabinet.cc::129)",
	}

	warningEntry := models.EntryWithLogLevel{
		Level: "WARN",
		Rest: "Wed Jun 10 14:55:31 2020     : Timeout protocol[plc-udp] url[fe80::4021:ff:fe00:10:61616]" +
			" (plc_bridge_connector.cc::227)",
	}

	infoConsumptionEntry := models.EntryWithLogLevel{
		Level: "INFO",
		Rest: "Wed Jun 10 14:56:19 2020     : <--[consumption]--(SMC) start[1591800000] end[1591800600]" +
			" value[0] service_level_id[9] (abstract_smart_meter_cabinet.h::314)",
	}

	expectedErrorEntry := models.EntryWithLevelAndTimestamp{
		Level:     "ERROR",
		Timestamp: time.Date(2020, time.June, 10, 14, 56, 3, 0, time.UTC),
		Rest: "error_code[241] message[DLMS error] severity[3] description[n/a] source[dc18-smc18]" +
			"  (smart_meter_cabinet.cc::129)",
	}

	expectedWarningEntry := models.EntryWithLevelAndTimestamp{
		Level:     "WARN",
		Timestamp: time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
		Rest:      "Timeout protocol[plc-udp] url[fe80::4021:ff:fe00:10:61616] (plc_bridge_connector.cc::227)",
	}

	expectedInfoConsumptionEntry := models.EntryWithLevelAndTimestamp{
		Level:     "INFO",
		Timestamp: time.Date(2020, time.June, 10, 14, 56, 19, 0, time.UTC),

		Rest: "<--[consumption]--(SMC) start[1591800000] end[1591800600] value[0] service_level_id[9]" +
			" (abstract_smart_meter_cabinet.h::314)",
	}

	parsedErrorEntry := service.ParseTimestamp(errorEntry)
	parsedWarningEntry := service.ParseTimestamp(warningEntry)
	parsedInfoConsumptionEntry := service.ParseTimestamp(infoConsumptionEntry)

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
