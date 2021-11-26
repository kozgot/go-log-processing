package elastic

import (
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
)

// EsClient interface exposes elasticsearch functionality used by the uploader service.
type EsClient interface {
	BulkUpload(dataUnits []models.ESDocument, indexName string)
	CreateEsIndex(index string)
}
