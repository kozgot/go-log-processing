package formats

// AnyStringUntilFirstSpaceRegex represents the regular expression that matches the first 'word' in a string until the first space
const AnyStringUntilFirstSpaceRegex = "([^\\s]+)"

// LongNumberRegex matches any number
const LongNumberRegex = "[0-9]*"

// LongNumberBetweenBracketsRegex matches any number between square brackets
const LongNumberBetweenBracketsRegex = "\\[" + LongNumberRegex + "\\]"

// AnyLettersBetweenBrackets matches strings containing any lowercase or uppercase letters between square brackets
const AnyLettersBetweenBrackets = "\\[(.*?)\\]"

// NumberBetweenBrackets matches strings containing a single number between square brackets
const NumberBetweenBrackets = "\\[[0-9]\\]"

const anyCharsExceptOpeningParentheses = "[^\\(^\\[]*"

// FileNameRegex matches the (...) field
const FileNameRegex = "\\(" + anyCharsExceptOpeningParentheses + "::" + LongNumberRegex + "\\)"

// CreationTimeRegex represents the regular expression that matches the creation_time[Wed Jun 10 09:18:39 2020] field in an error line
const CreationTimeRegex = "creation_time\\[" + anyCharsExceptOpeningParentheses + "\\]"

// MinLaunchTimeRegex represents the regular expression that matches the min_launch_time[Wed Jun 10 09:18:39 2020] field in an error line
const MinLaunchTimeRegex = "min_launch_time\\[" + anyCharsExceptOpeningParentheses + "\\]"

// SMCUIDRegex represents the regular expression that matches the smc_uid[smcidentifier] field in a log message
const SMCUIDRegex = "smc_uid\\[" + anyCharsExceptOpeningParentheses + "\\]"

// UIDRegex represents the regular expression that matches the uid[number] field in a log message
const UIDRegex = "uid" + LongNumberBetweenBracketsRegex

// InComingArrow represents the direction of messages sent to the dc main
const InComingArrow = "<--"

// OutGoingArrow represents the direction of messages sent by the dc main
const OutGoingArrow = "-->"
