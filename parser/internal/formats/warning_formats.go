package formats

// TaskFailedWarnRegex represents the regular expression that matches relevant a warn line.
const TaskFailedWarnRegex = "Task failed, "

// WarningPriorityRegex represents the regular expression that matches the severity[number] field in an error line.
const WarningPriorityRegex = "priority" + LongNumberBetweenBracketsRegex

// WarningNameRegex matches the name[...] field.
const WarningNameRegex = "name" + AnyLettersBetweenBrackets

// WarningRetryRegex represents the regular expression that matches the severity[number] field in an error line.
const WarningRetryRegex = "retry" + LongNumberBetweenBracketsRegex

// TimeoutWarnPrefix is the prefix of timeout warn level log entries.
// eg.: Timeout protocol[plc-udp] url[fe80::4021:ff:fe00:a:61616] (plc_bridge_connector.cc::227).
const TimeoutWarnPrefix = "Timeout "

// TimeoutProtocolRegex matches the protocol field in a timeout warn level log entry.
const TimeoutProtocolRegex = "protocol" + AnyLettersBetweenBrackets

// TimeoutURLRegex matches the url field in a timeout warn level log entry.
const TimeoutURLRegex = "url" + AnyLettersBetweenBrackets

// LostConnectionPrefix matches the start of connection lost log entries.
const LostConnectionPrefix = "Connection of type"
