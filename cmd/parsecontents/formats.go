package parsecontents

// ErrorCodeRegex represents the regular expression that matches the error_code[number] field in an error line
const ErrorCodeRegex = "error_code\\[" + NumberRegex + "\\]"

// NumberRegex matches any number
const NumberRegex = "[0-9]+"

// LongNumberRegex matches any number
const LongNumberRegex = "[0-9]*"

// AbcRegex matches strings containing any lowercase or uppercase letters
const AbcRegex = "(.*?)"

// MessageRegex matches the message[...] field
const MessageRegex = "message\\[" + AbcRegex + "\\]"

// OpeningBracketRegex matches [
const OpeningBracketRegex = "\\["

// ClosingBracketRegex matches ]
const ClosingBracketRegex = "\\]"

// ErrorSeverityRegex represents the regular expression that matches the severity[number] field in an error line
const ErrorSeverityRegex = "severity\\[" + "[0-9]" + "\\]"

// ErrorDescRegex represents the regular expression that matches the description[...] field in an error line
const ErrorDescRegex = "description\\[" + "(.*?)" + "\\]"

// ErrorSourceRegex represents the regular expression that matches the source[...] field in an error line
const ErrorSourceRegex = "source\\[" + "(.*?)" + "\\]"

// WarnRegex represents the regular expression that matches relevant a warn line
const WarnRegex = "Task failed, "

// WarningPriorityRegex represents the regular expression that matches the severity[number] field in an error line
const WarningPriorityRegex = "priority\\[" + LongNumberRegex + "\\]"

// WarningNameRegex matches the name[...] field
const WarningNameRegex = "name\\[" + AbcRegex + "\\]"

const anyCharsExceptParentheses = "[^\\(^\\[]*"

// FileNameRegex matches the (...) field
const FileNameRegex = "\\(" + anyCharsExceptParentheses + "::" + LongNumberRegex + "\\)"

// WarningRetryRegex represents the regular expression that matches the severity[number] field in an error line
const WarningRetryRegex = "retry\\[" + LongNumberRegex + "\\]"

// CreationTimeRegex represents the regular expression that matches the severity[number] field in an error line
const CreationTimeRegex = "creation_time\\[" + anyCharsExceptParentheses + "\\]"

// MinLaunchTimeRegex represents the regular expression that matches the severity[number] field in an error line
const MinLaunchTimeRegex = "min_launch_time\\[" + anyCharsExceptParentheses + "\\]"

// SMCUIDRegex represents the regular expression that matches the severity[number] field in an error line
const SMCUIDRegex = "smc_uid\\[" + anyCharsExceptParentheses + "\\]"

// UIDRegex represents the regular expression that matches the severity[number] field in an error line
const UIDRegex = "uid\\[" + LongNumberRegex + "\\]"
