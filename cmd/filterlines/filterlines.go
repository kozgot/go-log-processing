package filterlines

import (
	"regexp"
	"strings"
)

// Filter decides if a line is relevant in the input file
func Filter(line string) (*Line, bool) {
	levelRegex, _ := regexp.Compile(LogLevelsRegex)
	logLevel := levelRegex.FindString(line)
	if logLevel != "" {
		restOfLine := strings.Replace(line, logLevel, "", 1)
		return &Line{Level: logLevel, Rest: restOfLine}, true
	}

	// this line is not important
	return nil, false
}
