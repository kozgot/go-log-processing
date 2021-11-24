package contentparser

import (
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type InitConnectionEntryParser struct {
	line models.EntryWithLevelAndTimestamp
}

func (i *InitConnectionEntryParser) Parse() *models.InitConnectionParams {
	if strings.Contains(i.line.Rest, formats.InitConnectionPrefix) {
		// the entry looks like this:
		// Initialize DLMS connection fe80::4021:ff:fe00:a:61616 (some file name that is irrelevant)
		url := parseURLFromConnectionEntries(strings.Replace(i.line.Rest, formats.InitConnectionPrefix, "", 1))
		initConnectionParams := models.InitConnectionParams{URL: url}

		return &initConnectionParams
	}

	return nil
}
