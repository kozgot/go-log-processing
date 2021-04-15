package service

import (
	"regexp"
	"strings"

	"github.com/kozgot/go-log-processing/parser/pkg/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// Filter decides if a line is relevant in the input file
func Filter(line string) (*models.Line, bool) {
	levelRegex, _ := regexp.Compile(formats.LogLevelsRegex)
	logLevel := levelRegex.FindString(line)
	if logLevel != "" {
		restOfLine := strings.Replace(line, logLevel, "", 1)
		return &models.Line{Level: logLevel, Rest: restOfLine}, true
	}

	// could not parse log level
	return nil, false
}
