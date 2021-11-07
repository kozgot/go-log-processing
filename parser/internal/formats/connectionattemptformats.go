package formats

// This file contains regular expressions needed to
// parse log entries that look like this (after the timestamp and log level):
// Attempt to connect to SMC_dc18-smc27 (@ 0021) at URL fe80::4021:ff:fe00:21:61616 ..

// ConnectionAttemptPrefix is the string that connection attempt entries start with.
const ConnectionAttemptPrefix = "Attempt to connect to SMC_"

// AtRegex matches the "(@ 0021)" field of connaction attempt entries.
const AtRegex = "\\(" + anyCharsExceptOpeningParentheses + "\\)"

// URLPrefix is the string that the url fields of connection attempt entries start with.
const URLPrefix = "at URL "
