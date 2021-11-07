package service

import (
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

func parseSmcConfigUpdate(line string) *models.SmcConfigUpdateParams {
	if strings.Contains(line, formats.SmcConfigUpdatePrefix) {
		smcAddress := parseSmcAddressPayload(line)
		smcUID := parseFieldInBracketsAsString(line, formats.SmcUIDRegex)

		if smcAddress != nil {
			smcConfigUpdate := models.SmcConfigUpdateParams{
				PhysicalAddress: smcAddress.PhysicalAddress,
				LogicalAddress:  smcAddress.LogicalAddress,
				ShortAddress:    smcAddress.ShortAddress,
				SmcUID:          smcUID,
				LastJoiningDate: smcAddress.LastJoiningDate}
			return &smcConfigUpdate
		}
	}

	return nil
}
