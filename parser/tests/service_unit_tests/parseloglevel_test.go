package serviceunittests

import (
	"reflect"
	"testing"

	"github.com/kozgot/go-log-processing/parser/internal/loglevelparser"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func TestParseLogLevel(t *testing.T) {
	errorEntry := "Wed Jun 10 14:56:03 2020 ERROR   : error_code[241] message[DLMS error] severity[3] " +
		"description[n/a] source[dc18-smc18]  (smart_meter_cabinet.cc::129)"

	warningEntry := "Wed Jun 10 14:55:31 2020 WARN    : Timeout protocol[plc-udp] url[fe80::4021:ff:fe00:10:61616]" +
		" (plc_bridge_connector.cc::227)"

	infoConsEntry := "Wed Jun 10 14:56:19 2020 INFO    : <--[consumption]--(SMC) start[1591800000] end[1591800600]" +
		" value[0] service_level_id[9] (abstract_smart_meter_cabinet.h::314)"

	expectedErrorEntry := models.EntryWithLogLevel{
		Level: "ERROR",
		Rest: "Wed Jun 10 14:56:03 2020    : error_code[241] message[DLMS error] severity[3] " +
			"description[n/a] source[dc18-smc18]  (smart_meter_cabinet.cc::129)",
	}

	expectedWarningEntry := models.EntryWithLogLevel{
		Level: "WARN",
		Rest: "Wed Jun 10 14:55:31 2020     : Timeout protocol[plc-udp] url[fe80::4021:ff:fe00:10:61616]" +
			" (plc_bridge_connector.cc::227)",
	}

	expectedInfoConsumptionEntry := models.EntryWithLogLevel{
		Level: "INFO",
		Rest: "Wed Jun 10 14:56:19 2020     : <--[consumption]--(SMC) start[1591800000] end[1591800600]" +
			" value[0] service_level_id[9] (abstract_smart_meter_cabinet.h::314)",
	}

	parsedErrorEntry := loglevelparser.ParseLogLevelAndFilter(errorEntry)
	if parsedErrorEntry == nil {
		t.Fatal("Did not mark relevant ERROR line as relevant")
	}

	parsedWarningEntry := loglevelparser.ParseLogLevelAndFilter(warningEntry)
	if parsedWarningEntry == nil {
		t.Fatal("Did not mark relevant WARN line as relevant")
	}

	parsedInfoConsumptionEntry := loglevelparser.ParseLogLevelAndFilter(infoConsEntry)
	if parsedInfoConsumptionEntry == nil {
		t.Fatal("Did not mark relevant INFO line as relevant")
	}

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
