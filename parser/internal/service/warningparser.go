package service

import (
	"regexp"
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/parser/pkg/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// Parses an smc join log entry with a WARNING log level.
func parseWarning(line models.LineWithDate) *models.WarningParams {
	warningParams := models.WarningParams{}
	smcJoinParams := parseSmcJoinLine(line.Rest)
	warningParams.JoinMessageParams = *smcJoinParams

	return &warningParams
}

// Parses WARN level log entries from the dc_main.log file.
func parseWarn(line models.LineWithDate) *models.WarningParams {
	warningParams := models.WarningParams{}

	warnRegex, _ := regexp.Compile(formats.WarnRegex)
	warn := warnRegex.FindString(line.Rest)
	if warn == "" {
		return nil
	}

	// parse SMC UID
	smcUID := parseWarningSMCUID(line.Rest)
	warningParams.SmcUID = smcUID

	// parse UID
	uid := parseWarningUID(line.Rest)
	warningParams.UID = uid

	// parse Priority
	priority := parseWarningPriority(line.Rest)
	warningParams.Priority = priority

	// parse Name
	name := parseWarningName(line.Rest)
	warningParams.Name = name

	// parse FileName
	fileName := parseFileName(line.Rest)
	warningParams.FileName = fileName

	// parse Retry
	retry := parseWarningRetry(line.Rest)
	warningParams.Retry = retry

	// parse Creation
	creationTime := parseWarningCreationTime(line.Rest)
	warningParams.Creation = creationTime

	// parse MinLaunchTime
	minLaunchTime := parseWarningMinLaunchTime(line.Rest)
	warningParams.MinLaunchTime = minLaunchTime

	// parse inner error params
	errorParams := parseError(line)
	warningParams.Details = *errorParams

	return &warningParams
}

func parseWarningPriority(line string) int {
	return tryParseIntFromString(parseFieldInBracketsAsString(line, formats.WarningPriorityRegex))
}

func parseWarningRetry(line string) int {
	return tryParseIntFromString(parseFieldInBracketsAsString(line, formats.WarningRetryRegex))
}

func parseWarningUID(line string) int {
	return tryParseIntFromString(parseFieldInBracketsAsString(line, formats.UIDRegex))
}

func parseWarningName(line string) string {
	return parseFieldInBracketsAsString(line, formats.WarningNameRegex)
}

func parseWarningSMCUID(line string) string {
	return parseFieldInBracketsAsString(line, formats.SMCUIDRegex)
}

func parseWarningCreationTime(line string) time.Time {
	return parseDateTimeField(line, formats.CreationTimeRegex)
}

func parseWarningMinLaunchTime(line string) time.Time {
	return parseDateTimeField(line, formats.MinLaunchTimeRegex)
}

func parseDateTimeField(line string, regex string) time.Time {
	timeFieldRegex, _ := regexp.Compile(regex)
	timeField := timeFieldRegex.FindString(line)

	if timeField != "" {
		timeString := strings.Split(timeField, "[")[1]
		timeString = strings.Replace(timeString, "]", "", 1)

		dateTime := parseDateTime(timeString)
		return dateTime
	}

	return time.Time{}
}

func parseDateTime(timeString string) time.Time {
	dateRegex, _ := regexp.Compile(formats.DateFormatRegex)

	dateString := dateRegex.FindString(timeString)
	if dateString != "" {
		date, err := time.Parse(formats.DateLayoutString, dateString)
		if err != nil {
			// Do not die here, log instead
			panic(err)
		}

		return date
	}

	return time.Time{}
}

func parseFileName(line string) string {
	fileNameFieldRegex, _ := regexp.Compile(formats.FileNameRegex)
	fileNameField := fileNameFieldRegex.FindString(line)

	if fileNameField != "" {
		fileNameField = strings.Replace(fileNameField, ")", "", 1)
		fileNameField = strings.Replace(fileNameField, "(", "", 1)
		return fileNameField
	}

	return ""
}
