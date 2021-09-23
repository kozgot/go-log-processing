package service

import (
	"fmt"

	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// ParseContents extracts the params of an error level line in the log file
func ParseContents(line models.LineWithDate) *models.ParsedLine {
	parsedLine := models.ParsedLine{Level: line.Level, Timestamp: line.Timestamp}
	switch line.Level {
	case "ERROR":
		// todo: add these methods to the appropriate types,
		// eg.: ErrorParams should have a method parse() that parses error params
		errorParams := parseError(line)
		parsedLine.ErrorParams = *errorParams

	case "WARN":
		warning := parseWarn(line)
		if warning == nil {
			return nil
		}
		parsedLine.WarningParams = *warning

	case "WARNING":
		fmt.Println(line) // todo parse warning lines

	case "INFO":
		info := parseInfo(line)
		if info == nil {
			return nil
		}
		parsedLine.InfoParams = *info
	}

	return &parsedLine
}
