package contentparser

import (
	"regexp"
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// WarningParser is responsible for parsing log entries with WARN or WARNING log level.
type WarningParser struct {
	line models.EntryWithLevelAndTimestamp
}

func NewWarningParser(line models.EntryWithLevelAndTimestamp) *WarningParser {
	return &WarningParser{line: line}
}

// ParseWarning parses an smc join log entry with a WARNING log level.
func (warningParser *WarningParser) ParseWarning() *models.WarningParams {
	warningParams := models.WarningParams{}
	smcJoinEntryParser := SmcJoinEntryParser{line: warningParser.line}
	warningParams.JoinMessageParams = smcJoinEntryParser.Parse()

	return &warningParams
}

// ParseWarn parses WARN level log entries from the dc_main.log file.
func (warningParser *WarningParser) ParseWarn() *models.WarningParams {
	warningParams := models.WarningParams{}

	if strings.Contains(warningParser.line.Rest, formats.LostConnectionPrefix) {
		warningParams.LostConnectionParams = warningParser.parseLostConnectionEntry()
		return &warningParams
	}

	if strings.Contains(warningParser.line.Rest, formats.TimeoutWarnPrefix) {
		// This is a Timeout warn entry.
		timeoutParams := warningParser.parseTimeoutEntry()
		if timeoutParams != nil {
			warningParams.TimeoutParams = timeoutParams
		}

		return &warningParams
	}

	warnRegex, _ := regexp.Compile(formats.WarnRegex) // we only care for for Task failed warnings from here
	warn := warnRegex.FindString(warningParser.line.Rest)
	if warn == "" {
		return nil
	}

	// parse SMC UID
	smcUID := warningParser.parseWarningSMCUID()
	warningParams.SmcUID = smcUID

	// parse UID
	uid := warningParser.parseWarningUID()
	warningParams.UID = uid

	// parse Priority
	priority := warningParser.parsePriority()
	warningParams.Priority = priority

	// parse Name
	name := warningParser.parseWarningName()
	warningParams.Name = name

	// parse FileName
	fileName := warningParser.parseFileName()
	warningParams.FileName = fileName

	// parse Retry
	retry := warningParser.parseRetry()
	warningParams.Retry = retry

	// parse Creation
	creationTime := warningParser.parseWarningCreationTime()
	warningParams.Creation = creationTime

	// parse MinLaunchTime
	minLaunchTime := warningParser.parseWarningMinLaunchTime()
	warningParams.MinLaunchTime = minLaunchTime

	// parse inner error params
	errorParser := NewErrorParser(warningParser.line)
	warningParams.Details = errorParser.ParseError()

	return &warningParams
}

func (warningParser *WarningParser) parseLostConnectionEntry() *models.LostConnectionParams {
	result := warningParser.parseLostConnectionParams()

	// Get the reason from this: '...lost due to <unknown reason> (mqtt_connector.cc::54)'
	if strings.Contains(warningParser.line.Rest, "lost due to ") {
		reasonPart := strings.Split(warningParser.line.Rest, "lost due to ")[1]

		// Trim off the file name in the parentheses.
		if strings.Contains(reasonPart, " (") {
			result.Reason = strings.Split(reasonPart, " (")[0]
		}
	}

	return result
}

func (
	warningParser *WarningParser,
) parseLostConnectionParams() *models.LostConnectionParams {
	resultType := common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.ConnectOrDisconnectTypeRegex))
	clientID := common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.ClientIDRegex)
	connected := common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.ConnectedRegex)) == 1
	URL := common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.URLRegex)
	topic := common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.TopicRegex)
	timeout := common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.TimeoutRegex))

	if clientID != "" || resultType != 0 || URL != "" || topic != "" || timeout != 0 {
		result := models.LostConnectionParams{
			ClientID:  clientID,
			Type:      resultType,
			URL:       URL,
			Topic:     topic,
			Timeout:   timeout,
			Connected: connected}
		return &result
	}

	return nil
}

func (warningParser *WarningParser) parseTimeoutEntry() *models.TimelineOutParams {
	if !strings.Contains(warningParser.line.Rest, formats.TimeoutWarnPrefix) {
		// This is not a timeout warn log entry.
		return nil
	}

	result := models.TimelineOutParams{}
	result.Protocol = common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.TimeoutProtocolRegex)
	result.URL = common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.TimeoutURLRegex)

	return &result
}

func (warningParser *WarningParser) parsePriority() int {
	return common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.WarningPriorityRegex))
}

func (warningParser *WarningParser) parseRetry() int {
	return common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.WarningRetryRegex))
}

func (warningParser *WarningParser) parseWarningUID() int {
	return common.TryParseIntFromString(common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.UIDRegex))
}

func (warningParser *WarningParser) parseWarningName() string {
	return common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.WarningNameRegex)
}

func (warningParser *WarningParser) parseWarningSMCUID() string {
	return common.ParseFieldInBracketsAsString(warningParser.line.Rest, formats.SMCUIDRegex)
}

func (warningParser *WarningParser) parseWarningCreationTime() time.Time {
	return common.ParseDateTimeField(warningParser.line.Rest, formats.CreationTimeRegex)
}

func (warningParser *WarningParser) parseWarningMinLaunchTime() time.Time {
	return common.ParseDateTimeField(warningParser.line.Rest, formats.MinLaunchTimeRegex)
}

func (warningParser *WarningParser) parseFileName() string {
	fileNameFieldRegex, _ := regexp.Compile(formats.FileNameRegex)
	fileNameField := fileNameFieldRegex.FindString(warningParser.line.Rest)

	if fileNameField != "" {
		fileNameField = strings.Replace(fileNameField, ")", "", 1)
		fileNameField = strings.Replace(fileNameField, "(", "", 1)
		return fileNameField
	}

	return ""
}
