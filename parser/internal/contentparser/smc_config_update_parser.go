package contentparser

import (
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// SmcConfigUpdateParser parses an SMC config update log entry.
type SmcConfigUpdateParser struct {
	line models.EntryWithLevelAndTimestamp
}

// Parse parses an SMC config update log entry.
func (s *SmcConfigUpdateParser) Parse() *models.SmcConfigUpdateParams {
	if strings.Contains(s.line.Rest, formats.SmcConfigUpdatePrefix) {
		smcConfigUpdate := models.SmcConfigUpdateParams{}

		smcConfigUpdate.SmcUID = common.ParseFieldInBracketsAsString(s.line.Rest, formats.SmcUIDRegex)
		smcConfigUpdate.PhysicalAddress = common.ParseFieldInBracketsAsString(s.line.Rest, formats.PhysicalAddressRegex)
		smcConfigUpdate.LogicalAddress = common.ParseFieldInBracketsAsString(s.line.Rest, formats.LogicalAddressRegex)

		smcConfigUpdate.ShortAddress = common.TryParseIntFromString(
			common.ParseFieldInBracketsAsString(s.line.Rest, formats.ShortAddressRegex))

		smcConfigUpdate.LastJoiningDate = common.ParseDateTimeField(s.line.Rest, formats.LastJoiningDateRegex)
		smcUID := common.ParseFieldInBracketsAsString(s.line.Rest, formats.SmcUIDRegex)
		smcConfigUpdate.SmcUID = smcUID

		return &smcConfigUpdate
	}

	return nil
}
