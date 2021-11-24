package contentparser

import (
	"regexp"
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// Parses an smc join log entry with a WARNING log level.
func parseWarning(line models.EntryWithLevelAndTimestamp) *models.WarningParams {
	warningParams := models.WarningParams{}
	smcJoinParams := parseSmcJoinLine(line.Rest)
	warningParams.JoinMessageParams = smcJoinParams

	return &warningParams
}

// Parses WARN level log entries from the dc_main.log file.
func parseWarn(line models.EntryWithLevelAndTimestamp) *models.WarningParams {
	warningParams := models.WarningParams{}

	if strings.Contains(line.Rest, formats.LostConnectionPrefix) {
		warningParams.LostConnectionParams = parseLostConnectionParams(line.Rest)
		return &warningParams
	}

	if strings.Contains(line.Rest, formats.TimeoutWarnPrefix) {
		// This is a Timeout warn entry.
		timeoutParams := parseTimeoutEntry(line.Rest)
		if timeoutParams != nil {
			warningParams.TimeoutParams = timeoutParams
		}

		return &warningParams
	}

	warnRegex, _ := regexp.Compile(formats.WarnRegex) // we only care for for Task failed warnings from here
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
	innerErrorParams := ParseError(line)
	warningParams.Details = innerErrorParams

	return &warningParams
}

func parseLostConnectionParams(line string) *models.LostConnectionParams {
	result := models.LostConnectionParams{}
	params := parseConnectOrDisconnectPayload(line)
	if params != nil {
		result.ClientID = params.ClientID
		result.Connected = params.Connected
		result.Timeout = params.Timeout
		result.Topic = params.Topic
		result.URL = params.URL
		result.Type = params.Type
	}

	// Get the reason from this: '...lost due to <unknown reason> (mqtt_connector.cc::54)'
	if strings.Contains(line, "lost due to ") {
		reasonPart := strings.Split(line, "lost due to ")[1]

		// Trim off the file name in the parentheses.
		if strings.Contains(reasonPart, " (") {
			result.Reason = strings.Split(reasonPart, " (")[0]
		}
	}

	return &result
}

func parseTimeoutEntry(line string) *models.TimelineOutParams {
	if !strings.Contains(line, formats.TimeoutWarnPrefix) {
		// This is not a timeout warn log entry.
		return nil
	}

	result := models.TimelineOutParams{}
	result.Protocol = common.ParseFieldInBracketsAsString(line, formats.TimeoutProtocolRegex)
	result.URL = common.ParseFieldInBracketsAsString(line, formats.TimeoutURLRegex)

	return &result
}

func parseWarningPriority(line string) int {
	return common.TryParseIntFromString(common.ParseFieldInBracketsAsString(line, formats.WarningPriorityRegex))
}

func parseWarningRetry(line string) int {
	return common.TryParseIntFromString(common.ParseFieldInBracketsAsString(line, formats.WarningRetryRegex))
}

func parseWarningUID(line string) int {
	return common.TryParseIntFromString(common.ParseFieldInBracketsAsString(line, formats.UIDRegex))
}

func parseWarningName(line string) string {
	return common.ParseFieldInBracketsAsString(line, formats.WarningNameRegex)
}

func parseWarningSMCUID(line string) string {
	return common.ParseFieldInBracketsAsString(line, formats.SMCUIDRegex)
}

func parseWarningCreationTime(line string) time.Time {
	return common.ParseDateTimeField(line, formats.CreationTimeRegex)
}

func parseWarningMinLaunchTime(line string) time.Time {
	return common.ParseDateTimeField(line, formats.MinLaunchTimeRegex)
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
