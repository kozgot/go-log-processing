package loglevelparser

import (
	"regexp"
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// ParseLogLevelAndFilter decides if a line is relevant in the input file.
func ParseLogLevelAndFilter(line string) *models.EntryWithLogLevel {
	levelRegex, _ := regexp.Compile(formats.LogLevelsRegex)
	logLevel := levelRegex.FindString(line)
	if logLevel != "" {
		restOfLine := strings.Replace(line, logLevel, "", 1)
		return &models.EntryWithLogLevel{Level: logLevel, Rest: restOfLine}
	}

	// could not parse log level
	return nil
}
