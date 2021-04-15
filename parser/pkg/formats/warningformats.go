package formats

// WarnRegex represents the regular expression that matches relevant a warn line
const WarnRegex = "Task failed, "

// WarningPriorityRegex represents the regular expression that matches the severity[number] field in an error line
const WarningPriorityRegex = "priority" + LongNumberBetweenBracketsRegex

// WarningNameRegex matches the name[...] field
const WarningNameRegex = "name" + AnyLettersBetweenBrackets

// WarningRetryRegex represents the regular expression that matches the severity[number] field in an error line
const WarningRetryRegex = "retry" + LongNumberBetweenBracketsRegex
