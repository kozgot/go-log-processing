package processing

import (
	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// ProcessError processes a log entry with ERROR log level.
func ProcessError(logEntry parsermodels.ParsedLine) *models.SmcEntry {
	result := models.SmcEntry{}

	// todo
	return &result
}
