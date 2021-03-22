package parsedates

// DateFormatRegex represents the regular expression that matches the timestamp in a line of the dc_main.log file
const DateFormatRegex = "^(Mon|Tue|Wed|Thu|Fri|Sat|Sun) (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) (0?[1-9]|[12][0-9]|3[01]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9] [0-2][0-9][0-9][0-9]"

// DateFormatRegexShort represents the regular expression that matches the timestamp in a line of the plc_manager.log file
const DateFormatRegexShort = "([12]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01])-(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))"

// DateSurroundingRegex represents the regular expression that matches the whitespaces, square brackets, and colon that might sourround the timestamp in a line of a log file
const DateSurroundingRegex = "\\[*( )*\\]*(:)*"

// DateBetweenBracketsRegex matches a date between [], in the following format: [Wed Jun 10 09:18:33 2020]
const DateBetweenBracketsRegex = "\\[" + DateFormatRegex + "\\]"
