package formats

// This file contains regular expressions needed to
// parse log entries that look like this (after the timestamp and log level):
// Initialize DLMS connection fe80::4021:ff:fe00:a:61616 (smart_meter_cabinet_initializer.cc::115)

// InitConnectionPrefix is the prefix of initialize dlms connection entries.
const InitConnectionPrefix = "Initialize DLMS connection "
