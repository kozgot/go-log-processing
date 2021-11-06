package service

import (
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func parseInitConnectionLogEntry(line string) *models.InitConnectionParams {
	if strings.Contains(line, formats.InitConnectionPrefix) {
		// the entry looks like this:
		// Initialize DLMS connection fe80::4021:ff:fe00:a:61616 (some file name that is irrelevant)
		url := parseURLFromConnectionEntries(strings.Replace(line, formats.InitConnectionPrefix, "", 1))
		initConnectionParams := models.InitConnectionParams{URL: url}

		return &initConnectionParams
	}

	return nil
}
