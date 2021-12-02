package timestampparserunittests

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/kozgot/go-log-processing/parser/internal/timestampparser"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type timestampParserTest struct {
	input          models.EntryWithLogLevel
	expectedOutput models.EntryWithLevelAndTimestamp
}

func TestParseTimestamp(t *testing.T) {
	tests := []timestampParserTest{
		{
			input: models.EntryWithLogLevel{
				Level: "ERROR",
				Rest: "Wed Jun 10 14:56:03 2020    : error_code[241] message[DLMS error] severity[3] " +
					"description[n/a] source[dc18-smc18]  (smart_meter_cabinet.cc::129)",
			},
			expectedOutput: models.EntryWithLevelAndTimestamp{
				Level:     "ERROR",
				Timestamp: time.Date(2020, time.June, 10, 14, 56, 3, 0, time.UTC),
				Rest: "error_code[241] message[DLMS error] severity[3] description[n/a] source[dc18-smc18]" +
					"  (smart_meter_cabinet.cc::129)",
			},
		},
		{
			input: models.EntryWithLogLevel{
				Level: "WARN",
				Rest: "Wed Jun 10 14:55:31 2020     : Timeout protocol[plc-udp] url[fe80::4021:ff:fe00:10:61616]" +
					" (plc_bridge_connector.cc::227)",
			},
			expectedOutput: models.EntryWithLevelAndTimestamp{
				Level:     "WARN",
				Timestamp: time.Date(2020, time.June, 10, 14, 55, 31, 0, time.UTC),
				Rest:      "Timeout protocol[plc-udp] url[fe80::4021:ff:fe00:10:61616] (plc_bridge_connector.cc::227)",
			},
		},
		{
			input: models.EntryWithLogLevel{
				Level: "INFO",
				Rest: "Wed Jun 10 14:56:19 2020     : <--[consumption]--(SMC) start[1591800000] end[1591800600]" +
					" value[0] service_level_id[9] (abstract_smart_meter_cabinet.h::314)",
			},
			expectedOutput: models.EntryWithLevelAndTimestamp{
				Level:     "INFO",
				Timestamp: time.Date(2020, time.June, 10, 14, 56, 19, 0, time.UTC),

				Rest: "<--[consumption]--(SMC) start[1591800000] end[1591800600] value[0] service_level_id[9]" +
					" (abstract_smart_meter_cabinet.h::314)",
			},
		},
		{
			input: models.EntryWithLogLevel{
				Level: "INFO",
				Rest: "[ 2020-06-10-09:20:15 ]: SMC Join OK [Confirmed] <-- [join_type[LBD] smc_uid[dc18-smc32]" +
					" physical_address[EEBEDDFFFE6210AD] logical_address[FE80::4021:FF:FE00:000a:61616]" +
					" short_address[10] last_joining_date[Wed Jun 10 09:20:14 2020]]--(PLC)",
			},
			expectedOutput: models.EntryWithLevelAndTimestamp{
				Level:     "INFO",
				Timestamp: time.Date(2020, time.June, 10, 9, 20, 15, 0, time.UTC),

				Rest: "SMC Join OK [Confirmed] <-- [join_type[LBD] smc_uid[dc18-smc32]" +
					" physical_address[EEBEDDFFFE6210AD] logical_address[FE80::4021:FF:FE00:000a:61616]" +
					" short_address[10] last_joining_date[Wed Jun 10 09:20:14 2020]]--(PLC)",
			},
		},
	}

	for index, test := range tests {
		output := timestampparser.ParseTimestamp(test.input)

		if !reflect.DeepEqual(*output, test.expectedOutput) {
			t.Fatalf("Parsed entry does not match expected parsed entry in test no. %d", index)
		}
	}

	log.Printf("Successfully ran %d test cases.", len(tests))
}
