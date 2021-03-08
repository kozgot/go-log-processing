package parsecontents

// ErrorCodeRegex represents the regular expression that matches the error_code[number] field in an error line
const ErrorCodeRegex = "error_code\\[" + NumberRegex + "\\]"

// NumberRegex matches any number
const NumberRegex = "[0-9]+"

// LongNumberRegex matches any number
const LongNumberRegex = "[0-9]*"

// LongNumberBetweenBracketsRegex matches any number between square brackets
const LongNumberBetweenBracketsRegex = "\\[" + LongNumberRegex + "\\]"

// AnyLettersBetweenBrackets matches strings containing any lowercase or uppercase letters
const AnyLettersBetweenBrackets = "\\[(.*?)\\]"

// MessageRegex matches the message[...] field
const MessageRegex = "message" + AnyLettersBetweenBrackets

// ErrorSeverityRegex represents the regular expression that matches the severity[number] field in an error line
const ErrorSeverityRegex = "severity\\[" + "[0-9]" + "\\]"

// ErrorDescRegex represents the regular expression that matches the description[...] field in an error line
const ErrorDescRegex = "description\\[" + "(.*?)" + "\\]"

// ErrorSourceRegex represents the regular expression that matches the source[...] field in an error line
const ErrorSourceRegex = "source\\[" + "(.*?)" + "\\]"

// WarnRegex represents the regular expression that matches relevant a warn line
const WarnRegex = "Task failed, "

// WarningPriorityRegex represents the regular expression that matches the severity[number] field in an error line
const WarningPriorityRegex = "priority" + LongNumberBetweenBracketsRegex

// WarningNameRegex matches the name[...] field
const WarningNameRegex = "name" + AnyLettersBetweenBrackets

const anyCharsExceptOpeningParentheses = "[^\\(^\\[]*"

// FileNameRegex matches the (...) field
const FileNameRegex = "\\(" + anyCharsExceptOpeningParentheses + "::" + LongNumberRegex + "\\)"

// WarningRetryRegex represents the regular expression that matches the severity[number] field in an error line
const WarningRetryRegex = "retry" + LongNumberBetweenBracketsRegex

// CreationTimeRegex represents the regular expression that matches the severity[number] field in an error line
const CreationTimeRegex = "creation_time\\[" + anyCharsExceptOpeningParentheses + "\\]"

// MinLaunchTimeRegex represents the regular expression that matches the severity[number] field in an error line
const MinLaunchTimeRegex = "min_launch_time\\[" + anyCharsExceptOpeningParentheses + "\\]"

// SMCUIDRegex represents the regular expression that matches the smc_uid[smcidentifier] field in a log message
const SMCUIDRegex = "smc_uid\\[" + anyCharsExceptOpeningParentheses + "\\]"

// UIDRegex represents the regular expression that matches the uid[number] field in a log message
const UIDRegex = "uid" + LongNumberBetweenBracketsRegex
