package service

import (
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func parseSmcInternalDiagnosticsEntry(line string) *models.InternalDiagnosticsData {
	if strings.Contains(line, formats.SmcInternalDiagnosticsPrefix) {
		// the entry looks like this:
		// SMC internal diagnostics smc_uid[dc18-smc32] last_successful_dlms_response_date[n/a] (file_name...)
		diagnosticsData := models.InternalDiagnosticsData{}
		smcUID := parseFieldInBracketsAsString(line, formats.SMCUIDRegex)
		lastSuccessfulDlmsResponseDate := parseDateTimeField(line, formats.LastSuccessfulDlmsResponseDateRegex)
		diagnosticsData.SmcUID = smcUID
		diagnosticsData.LastSuccessfulDlmsResponseDate = lastSuccessfulDlmsResponseDate

		return &diagnosticsData
	}

	return nil
}
