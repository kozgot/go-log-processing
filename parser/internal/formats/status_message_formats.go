package formats

// This file contains the regular expressions used to parse log lines in the following format:
// LOADNG_SEQ_NUM_REPORTED status_byte[0xA5]<--[Network status]--(PLC)

// StatusMessageRegex represents the regular expression that matches a log line
// that contains data regarding a status message.
const StatusMessageRegex = AnyStringUntilFirstSpaceRegex

// StatusByteRegex  represents the regular expression that matches the status byte field of a log entry.
const StatusByteRegex = "status_byte" + AnyLettersBetweenBrackets
