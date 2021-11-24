package formats

// DateFormatRegex matches the timestamp in this format: 'Wed Jun 10 14:56:19 2020'.
const DateFormatRegex = "^(Mon|Tue|Wed|Thu|Fri|Sat|Sun)" +
	" (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)" +
	" (0?[1-9]|[12][0-9]|3[01]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9] [0-2][0-9][0-9][0-9]"

// DateFormatRegexShort matches the timestamp in this format: '2020-06-10-09:24:18'.
const DateFormatRegexShort = "([12]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01])" +
	"-(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))"

// DateSurroundingRegex matches the whitespaces, square brackets,
// and colon that might sourround the timestamp in a line of a log file.
// This is used to trim off the already parsed parts from the line.
const DateSurroundingRegex = "\\[*( )*\\]*(:)*"
