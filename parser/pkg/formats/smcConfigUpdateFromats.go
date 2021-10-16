package formats

// This file contains regular expressions and string constants needed to
// parse log entries that start like this (after the timestamp and log level):
// Update SMC configuration in DB

// SmcConfigUpdatePrefix is the prefix of smc config update log entries.
const SmcConfigUpdatePrefix = "Update SMC configuration in DB "

// The rest of the parameters of this log entry are used in different entries as well,
// so they are defined in another file.
// eg.: logical address, short address, smc uid...
