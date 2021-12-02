package mocks

import "github.com/kozgot/go-log-processing/elasticuploader/pkg/models"

// EsClientMock is a mock Elasticsearch client used in tests.
type EsClientMock struct {
	Indexes map[string][]models.ESDocument
}

// BulkUpload mocks the BulkUpload function of the EsClient interface.
func (m *EsClientMock) BulkUpload(dataUnits []models.ESDocument, indexName string) {
	m.Indexes[indexName] = append(m.Indexes[indexName], dataUnits...)
}

// CreateEsIndex mocks the CreateEsIndex function of the EsClient interface.
func (m *EsClientMock) CreateEsIndex(index string) {
	m.Indexes[index] = []models.ESDocument{}
}

func NewESClientMock(indexes map[string][]models.ESDocument) *EsClientMock {
	esMock := EsClientMock{
		Indexes: indexes,
	}

	return &esMock
}
