package formats

// This file contains regular expressions and string constants needed to
// parse log entries that start like this (after the timestamp and log level):
// Successfully Released DLMS connection fe80::4021:ff:fe00:23:61616 (sm...

// ConnectionReleasedPrefix is the prefix of connection released log entries.
const ConnectionReleasedPrefix = "Successfully Released DLMS connection "
