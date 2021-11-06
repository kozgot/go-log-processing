package formats

// This file contains regular expressions used to parse data from log lines (plc-manager-log) in the following format:
//  SMC Join OK [Confirmed] <-- [join_type[LBD] smc_uid[dc18-smc37] physical_address[EEBEDDFFFE62114D]
//  logical_address[FE80::4021:FF:FE00:0003:61616] short_address[3] last_joining_date[Wed Jun 10 09:26:41 2020]]--(PLC)

// SmcJoinRegex matches a log line that contains data regarding an SMC join event.
const SmcJoinRegex = "SMC Join "

// JoinStatusResponseRegex matches the response status of an SMC join event.
const JoinStatusResponseRegex = "(NOT OK|OK)"

// JoinTypeRegex matches the join type field of a log entry.
const JoinTypeRegex = "join_type" + AnyLettersBetweenBrackets

// SmcUIDRegex matches the smc uid field of a log line, formatted like: smc_uid[dc18-smc18].
const SmcUIDRegex = "smc_uid" + AnyLettersBetweenBrackets

// PhysicalAddressRegex matches the physical address field of a log line,
// formatted like: physical_address[EEBEDDFFFE6210AA].
const PhysicalAddressRegex = "physical_address" + AnyLettersBetweenBrackets

// LogicalAddressRegex matches the logical address field of a log line,
// formatted like: logical_address[FE80::4021:FF:FE00:0015:61616].
const LogicalAddressRegex = "logical_address" + AnyLettersBetweenBrackets

// ShortAddressRegex matches the short address field of a log line, formatted like: short_address[21].
const ShortAddressRegex = "short_address" + LongNumberBetweenBracketsRegex
