package parsecontents

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kozgot/go-log-processing/cmd/parsedates"
)

// ParseContents extracts the params of an error level line in the log file
func ParseContents(line parsedates.LineWithDate) *ParsedLine {
	parsedLine := ParsedLine{Level: line.Level, Timestamp: line.Timestamp}
	switch line.Level {
	case "ERROR":
		errorParams := parseError(line)
		parsedLine.ErrorParams = *errorParams

	case "WARN":
		warning := parseWarn(line)
		if warning == nil {
			return nil
		}
		parsedLine.WarningParams = *warning
	}

	return &parsedLine
}

func parseError(line parsedates.LineWithDate) *ErrorParams {
	errorCode := parseErrorCode(line.Rest)
	errorMessage := parseErrorMessage(line.Rest)
	errorSeverity := parseErrorSeverity(line.Rest)
	errorDesc := parseErrorDesc(line.Rest)
	errorSource := parseErrorSource(line.Rest)

	errorParams := ErrorParams{ErrorCode: errorCode, Message: errorMessage, Severity: errorSeverity, Description: errorDesc, Source: errorSource}

	return &errorParams
}

func parseErrorCode(line string) int {
	errorCodeFieldRegex, _ := regexp.Compile(ErrorCodeRegex)
	errorCodeField := errorCodeFieldRegex.FindString(line)

	if errorCodeField != "" {
		errorCodeRegex, _ := regexp.Compile(NumberRegex)
		errorCodeString := errorCodeRegex.FindString(line)
		errorCode, err := strconv.Atoi(errorCodeString)
		if err != nil {
			panic(err)
		}

		return errorCode
	}

	return 0
}

func parseErrorMessage(line string) string {
	errorMessageFieldRegex, _ := regexp.Compile(MessageRegex)
	errorMessageField := errorMessageFieldRegex.FindString(line)

	if errorMessageField != "" {
		// todo: match this with regex
		errorMessageString := strings.Split(errorMessageField, "[")[1]
		errorMessageString = strings.Replace(errorMessageString, "]", "", 1)

		return errorMessageString
	}

	return ""
}

func parseErrorDesc(line string) string {
	errorDescFieldRegex, _ := regexp.Compile(ErrorDescRegex)
	errorDescField := errorDescFieldRegex.FindString(line)

	if errorDescField != "" {
		// todo: match this with regex
		errorDescString := strings.Split(errorDescField, "[")[1]
		errorDescString = strings.Replace(errorDescString, "]", "", 1)

		return errorDescString
	}

	return ""
}

func parseErrorSource(line string) string {
	errorSourceFieldRegex, _ := regexp.Compile(ErrorSourceRegex)
	errorSourceField := errorSourceFieldRegex.FindString(line)

	if errorSourceField != "" {
		// todo: match this with regex
		errorSourceString := strings.Split(errorSourceField, "[")[1]
		errorSourceString = strings.Replace(errorSourceString, "]", "", 1)

		return errorSourceString
	}

	return ""
}

func parseErrorSeverity(line string) int {
	errorSeverityFieldRegex, _ := regexp.Compile(ErrorSeverityRegex)
	errorSeverityField := errorSeverityFieldRegex.FindString(line)

	if errorSeverityField != "" {
		errorSeverityRegex, _ := regexp.Compile(NumberRegex)
		errorSeverityString := errorSeverityRegex.FindString(errorSeverityField)
		errorSeverity, err := strconv.Atoi(errorSeverityString)
		if err != nil {
			panic(err)
		}

		return errorSeverity
	}

	return 0
}

func parseWarn(line parsedates.LineWithDate) *WarningParams {
	warningParams := WarningParams{}

	warnRegex, _ := regexp.Compile(WarnRegex)
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
	warningParams.ErrorParams = *errorParams

	// could not parse log level
	return &warningParams
}

func parseWarningPriority(line string) int {
	warningPriorityFieldRegex, _ := regexp.Compile(WarningPriorityRegex)
	warningPriorityField := warningPriorityFieldRegex.FindString(line)

	if warningPriorityField != "" {
		warningPriorityRegex, _ := regexp.Compile(NumberRegex)
		warningPriorityString := warningPriorityRegex.FindString(warningPriorityField)
		warningPriority, err := strconv.Atoi(warningPriorityString)
		if err != nil {
			panic(err)
		}

		return warningPriority
	}

	return 0
}

func parseWarningRetry(line string) int {
	warningRetryFieldRegex, _ := regexp.Compile(WarningRetryRegex)
	warningRetryField := warningRetryFieldRegex.FindString(line)

	if warningRetryField != "" {
		warningRetryRegex, _ := regexp.Compile(NumberRegex)
		warningRetryString := warningRetryRegex.FindString(warningRetryField)
		warningRetry, err := strconv.Atoi(warningRetryString)
		if err != nil {
			panic(err)
		}

		return warningRetry
	}

	return 0
}

func parseWarningUID(line string) int {
	warningUIDFieldRegex, _ := regexp.Compile(UIDRegex)
	warningUIDField := warningUIDFieldRegex.FindString(line)

	if warningUIDField != "" {
		warningUIDRegex, _ := regexp.Compile(NumberRegex)
		warningUIDString := warningUIDRegex.FindString(warningUIDField)
		warningUID, err := strconv.Atoi(warningUIDString)
		if err != nil {
			panic(err)
		}

		return warningUID
	}

	return 0
}

func parseWarningName(line string) string {
	warningNameFieldRegex, _ := regexp.Compile(WarningNameRegex)
	warningNameField := warningNameFieldRegex.FindString(line)

	if warningNameField != "" {
		warningNameString := strings.Split(warningNameField, "[")[1]
		warningNameString = strings.Replace(warningNameString, "]", "", 1)

		return warningNameString
	}

	return ""
}

func parseWarningSMCUID(line string) string {
	warningSMCUIDFieldRegex, _ := regexp.Compile(SMCUIDRegex)
	warningSMCUIDField := warningSMCUIDFieldRegex.FindString(line)

	if warningSMCUIDField != "" {
		// todo: match this with regex
		warningSMCUIDString := strings.Split(warningSMCUIDField, "[")[1]
		warningSMCUIDString = strings.Replace(warningSMCUIDString, "]", "", 1)

		return warningSMCUIDString
	}

	return ""
}

func parseWarningCreationTime(line string) time.Time {
	warningCreationTimeFieldRegex, _ := regexp.Compile(CreationTimeRegex)
	warningCreationTimeField := warningCreationTimeFieldRegex.FindString(line)

	if warningCreationTimeField != "" {
		// todo: match this with regex
		warningCreationTimeString := strings.Split(warningCreationTimeField, "[")[1]
		warningCreationTimeString = strings.Replace(warningCreationTimeString, "]", "", 1)

		dateTime := parseDateTime(warningCreationTimeString)
		return dateTime
	}

	return time.Time{}
}

func parseWarningMinLaunchTime(line string) time.Time {
	warningMinLaunchTimeFieldRegex, _ := regexp.Compile(MinLaunchTimeRegex)
	warningMinLaunchTimeField := warningMinLaunchTimeFieldRegex.FindString(line)

	if warningMinLaunchTimeField != "" {
		// todo: match this with regex
		warningMinLaunchTimeString := strings.Split(warningMinLaunchTimeField, "[")[1]
		warningMinLaunchTimeString = strings.Replace(warningMinLaunchTimeString, "]", "", 1)

		dateTime := parseDateTime(warningMinLaunchTimeString)
		return dateTime
	}

	return time.Time{}
}

func parseDateTime(timeString string) time.Time {
	dateRegex, _ := regexp.Compile(parsedates.DateFormatRegex)

	dateString := dateRegex.FindString(timeString)
	if dateString != "" {
		date, err := time.Parse(parsedates.DateLayoutString, dateString)
		if err != nil {
			// Do not die here, log instead
			panic(err)
		}

		return date
	}

	return time.Time{}
}

func parseFileName(line string) string {
	fileNameFieldRegex, _ := regexp.Compile(FileNameRegex)
	fileNameField := fileNameFieldRegex.FindString(line)

	if fileNameField != "" {
		fileNameField = strings.Replace(fileNameField, ")", "", 1)
		fileNameField = strings.Replace(fileNameField, "(", "", 1)
		return fileNameField
	}

	return ""
}
