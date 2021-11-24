package formats

// This file contains regular expressions needed to
// parse log entries that start like this (after the timestamp and log level):
// SMC internal diagnostics

// SmcInternalDiagnosticsPrefix is the prefix of smc internal diagnostics entries.
const SmcInternalDiagnosticsPrefix = "SMC internal diagnostics "

// LastSuccessfulDlmsResponseDateRegex matches the last successful dlms response date filed of
// an internal diagnostics log entry. It is not the same as LastSuccessfulRespDateRegex in dcmessageformats file,
// as they differ in a space before the [ character.
// eg.: last_successful_dlms_response_date[Wed Jun 10 12:51:33 2020].
const LastSuccessfulDlmsResponseDateRegex = "last_successful_dlms_response_date\\[" +
	anyCharsExceptOpeningParentheses +
	"\\]"
