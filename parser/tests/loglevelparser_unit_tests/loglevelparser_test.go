package loglevelparserunittests

import (
	"log"
	"reflect"
	"testing"

	"github.com/kozgot/go-log-processing/parser/internal/loglevelparser"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type loglevelParserTest struct {
	input          string
	expectedOutput *models.EntryWithLogLevel
}

func TestParseLogLevel(t *testing.T) {
	tests := []loglevelParserTest{
		{
			input: "Wed Jun 10 14:56:03 2020 ERROR   : error_code[241] message[DLMS error] severity[3] " +
				"description[n/a] source[dc18-smc18]  (smart_meter_cabinet.cc::129)",
			expectedOutput: &models.EntryWithLogLevel{
				Level: "ERROR",
				Rest: "Wed Jun 10 14:56:03 2020    : error_code[241] message[DLMS error] severity[3] " +
					"description[n/a] source[dc18-smc18]  (smart_meter_cabinet.cc::129)",
			},
		},
		{
			input: "Wed Jun 10 14:55:31 2020 WARN    : Timeout protocol[plc-udp] url[fe80::4021:ff:fe00:10:61616]" +
				" (plc_bridge_connector.cc::227)",
			expectedOutput: &models.EntryWithLogLevel{
				Level: "WARN",
				Rest: "Wed Jun 10 14:55:31 2020     : Timeout protocol[plc-udp] url[fe80::4021:ff:fe00:10:61616]" +
					" (plc_bridge_connector.cc::227)",
			},
		},
		{
			input: "Wed Jun 10 14:56:19 2020 INFO    : <--[consumption]--(SMC) start[1591800000] end[1591800600]" +
				" value[0] service_level_id[9] (abstract_smart_meter_cabinet.h::314)",
			expectedOutput: &models.EntryWithLogLevel{
				Level: "INFO",
				Rest: "Wed Jun 10 14:56:19 2020     : <--[consumption]--(SMC) start[1591800000] end[1591800600]" +
					" value[0] service_level_id[9] (abstract_smart_meter_cabinet.h::314)",
			},
		},
		{
			input:          "Wed Jun 10 09:18:28 2020 VERBOSE : Thread 3055825968 started! (abstract_thread.cc::43)",
			expectedOutput: nil,
		},
		{
			input:          "Wed Jun 10 09:18:28 2020 DEBUG   : Write PID (22369) to file to /opt/bin/dc_main.pid",
			expectedOutput: nil,
		},
		{
			input: "[ 2020-06-10-09:20:15 ]INFO: SMC Join OK [Confirmed] <-- [join_type[LBD] smc_uid[dc18-smc32]" +
				" physical_address[EEBEDDFFFE6210AD] logical_address[FE80::4021:FF:FE00:000a:61616]" +
				" short_address[10] last_joining_date[Wed Jun 10 09:20:14 2020]]--(PLC)",
			expectedOutput: &models.EntryWithLogLevel{
				Level: "INFO",
				Rest: "[ 2020-06-10-09:20:15 ]: SMC Join OK [Confirmed] <-- [join_type[LBD] smc_uid[dc18-smc32]" +
					" physical_address[EEBEDDFFFE6210AD] logical_address[FE80::4021:FF:FE00:000a:61616]" +
					" short_address[10] last_joining_date[Wed Jun 10 09:20:14 2020]]--(PLC)",
			},
		},
		{
			input: "[ 2020-06-10-09:42:41 ]WARNING: SMC Join NOT OK [Rejected] <-- [join_type[LBA]" +
				" smc_uid[dc18-smc26] physical_address[EEBEDDFFFE621098] logical_address[FE80::4021:FF:FE00:0002:61616]" +
				" short_address[2] last_joining_date[Wed Jun 10 09:42:31 2020]]--(PLC) Reason[UNKNOWN_STATUS]",
			expectedOutput: &models.EntryWithLogLevel{
				Level: "WARNING",
				Rest: "[ 2020-06-10-09:42:41 ]: SMC Join NOT OK [Rejected] <-- [join_type[LBA]" +
					" smc_uid[dc18-smc26] physical_address[EEBEDDFFFE621098] logical_address[FE80::4021:FF:FE00:0002:61616]" +
					" short_address[2] last_joining_date[Wed Jun 10 09:42:31 2020]]--(PLC) Reason[UNKNOWN_STATUS]",
			},
		},
	}

	for index, test := range tests {
		output := loglevelparser.ParseLogLevelAndFilter(test.input)

		if test.expectedOutput == nil && output != nil {
			t.Fatalf("Did not filter out irrelevant entry in test no. %d", index+1)
		}

		if test.expectedOutput != nil && output == nil {
			t.Fatalf("Got nil result for relevant entry in test no. %d", index+1)
		}

		if !reflect.DeepEqual(output, test.expectedOutput) {
			t.Fatalf("Output of loglevelparser doest not match expected output in test no. %d", index+1)
		}
	}

	log.Printf("Successfully ran %d test cases.", len(tests))
}
