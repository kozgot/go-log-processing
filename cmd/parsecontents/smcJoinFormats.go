package parsecontents

// This file contains regular expressions used to parse data from log lines (plc-manager-log) in the following format:
//  SMC Join OK [Confirmed] <-- [join_type[LBD] smc_uid[dc18-smc37] physical_address[EEBEDDFFFE62114D] logical_address[FE80::4021:FF:FE00:0003:61616] short_address[3] last_joining_date[Wed Jun 10 09:26:41 2020]]--(PLC)

// SmcJoinRegex represents the regular expression that matches a log line that contains data regarding an SMC join event
const SmcJoinRegex = "SMC Join "

// JoinStatusResponseRegex represents the regular expression that matches the response status of an SMC join event
const JoinStatusResponseRegex = "(NOT OK|OK)"

// InComingMessageArrow represents the direction of messages sent to the dc main
const InComingMessageArrow = "<--"

// OutGoingMessageArrow represents the direction of messages sent by the dc main
const OutGoingMessageArrow = "-->"

const JoinTypeRegex = "join_type" + AnyLettersBetweenBrackets

const SmcUidRegex = "smc_uid" + AnyLettersBetweenBrackets

const PhysicalAddressRegex = "physical_address" + AnyLettersBetweenBrackets

const LogicalAddressRegex = "logical_address" + AnyLettersBetweenBrackets

const ShortAddressRegex = "short_address" + NumberBetweenBrackets
