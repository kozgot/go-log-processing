package formats

// MessageRegex matches the message[...] field
const MessageRegex = "message" + AnyLettersBetweenBrackets

// ErrorSeverityRegex represents the regular expression that matches the severity[number] field in an error line
const ErrorSeverityRegex = "severity" + NumberBetweenBrackets

// ErrorDescRegex represents the regular expression that matches the description[...] field in an error line
const ErrorDescRegex = "description\\[" + "(.*?)" + "\\]"

// ErrorSourceRegex represents the regular expression that matches the source[...] field in an error line
const ErrorSourceRegex = "source\\[" + "(.*?)" + "\\]"

// ErrorCodeRegex represents the regular expression that matches the error_code[number] field in an error line
const ErrorCodeRegex = "error_code" + NumberBetweenBrackets
