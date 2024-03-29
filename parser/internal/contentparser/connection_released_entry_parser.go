package contentparser

import (
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type ConnectionReleasedEntryParser struct {
	line models.EntryWithLevelAndTimestamp
}

func (c *ConnectionReleasedEntryParser) Parse() *models.ConnectionReleasedParams {
	if strings.Contains(c.line.Rest, formats.ConnectionReleasedPrefix) {
		// the entry looks like this:
		// Successfully Released DLMS connection fe80::4021:ff:fe00:23:61616 (smart_meter_cabinet_initializer.cc::115)
		url := parseURLFromConnectionEntries(strings.Replace(c.line.Rest, formats.ConnectionReleasedPrefix, "", 1))
		connectionReleasedParams := models.ConnectionReleasedParams{URL: url}

		return &connectionReleasedParams
	}

	return nil
}

func parseURLFromConnectionEntries(entryWithoutPrefix string) string {
	// We need to trim off the source file name in parentheses from the end.
	minLengthIfContainsSeparator := 2
	lineParts := strings.Split(entryWithoutPrefix, " (")
	if len(lineParts) < minLengthIfContainsSeparator {
		// If the entry does not contain the file name in parentheses,
		// then it probably is not in the correct format.
		return ""
	}

	return lineParts[0]
}
