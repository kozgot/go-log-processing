package contentparser

import (
	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type StatusEntryParser struct {
	line models.EntryWithLevelAndTimestamp
}

func (s *StatusEntryParser) Parse() *models.StatusMessageParams {
	statusLine := models.StatusMessageParams{}
	statusLine.StatusByte = common.ParseFieldInBracketsAsString(s.line.Rest, formats.StatusByteRegex)
	if statusLine.StatusByte == "" {
		return nil
	}

	statusLine.Message = common.ParseFieldAsString(s.line.Rest, formats.StatusMessageRegex)

	return &statusLine
}
