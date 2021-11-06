package service

import (
	"regexp"
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func parseConnectionAttempt(line string) *models.ConnectionAttemptParams {
	// parse entries like this:
	// Attempt to connect to SMC_dc18-smc27 (@ 0021) at URL fe80::4021:ff:fe00:21:61616 ..

	if !strings.Contains(line, formats.ConnectionAttemptPrefix) {
		// This is not a connection attempt entry.
		return nil
	}

	attempt := models.ConnectionAttemptParams{}

	// parse the SMC UID
	attempt.SmcUID = parseUIDFromConnectionAttempt(line)

	// Parse the (@ 0021) like part
	atFieldRegex, _ := regexp.Compile(formats.AtRegex)
	atFieldString := atFieldRegex.FindString(line)

	if atFieldString != "" {
		attempt.At = atFieldString
	}

	attempt.URL = parseURLFromConnectionAttempt(line)

	return &attempt
}

func parseURLFromConnectionAttempt(line string) string {
	// parse the fe80::4021:ff:fe00:21:61616 part from these kinds of entries:
	// Attempt to connect to SMC_dc18-smc27 (@ 0021) at URL fe80::4021:ff:fe00:21:61616 ..
	minLengthIfContainsSeparator := 2
	parts := strings.Split(line, formats.URLPrefix)
	if len(parts) < minLengthIfContainsSeparator {
		return ""
	}

	parts = strings.Split(parts[1], " ")
	if len(parts) < minLengthIfContainsSeparator {
		return ""
	}

	return parts[0]
}

func parseUIDFromConnectionAttempt(line string) string {
	// parse the dc18-smc27 part from these kinds of entries:
	// Attempt to connect to SMC_dc18-smc27 (@ 0021) at URL fe80::4021:ff:fe00:21:61616 ..
	minLengthIfContainsSeparator := 2
	parts := strings.Split(line, "(@")
	if len(parts) < minLengthIfContainsSeparator {
		return ""
	}

	// Remove the 'Attempt to connect to SMC_' part
	result := strings.Replace(parts[0], formats.ConnectionAttemptPrefix, "", 1)

	// Remove trailing whitespace
	result = strings.Trim(result, " ")

	return result
}
