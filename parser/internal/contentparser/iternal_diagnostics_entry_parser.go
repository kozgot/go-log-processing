package contentparser

import (
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type InternalDiagnosticsEntryParser struct {
	line models.EntryWithLevelAndTimestamp
}

func (i *InternalDiagnosticsEntryParser) Parse() *models.InternalDiagnosticsData {
	if strings.Contains(i.line.Rest, formats.SmcInternalDiagnosticsPrefix) {
		// the entry looks like this:
		// SMC internal diagnostics smc_uid[dc18-smc32] last_successful_dlms_response_date[n/a] (file_name...)
		diagnosticsData := models.InternalDiagnosticsData{}
		smcUID := common.ParseFieldInBracketsAsString(i.line.Rest, formats.SMCUIDRegex)
		lastSuccessfulDlmsResponseDate := common.ParseDateTimeField(i.line.Rest, formats.LastSuccessfulDlmsResponseDateRegex)
		diagnosticsData.SmcUID = smcUID
		diagnosticsData.LastSuccessfulDlmsResponseDate = lastSuccessfulDlmsResponseDate

		return &diagnosticsData
	}

	return nil
}
