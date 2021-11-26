package mocks

import "github.com/kozgot/go-log-processing/elasticuploader/pkg/models"

// EsClientMock is a mock Elasticsearch client used in tests.
type EsClientMock struct {
	Indexes      map[string][]models.DataUnit
	allDocsCount int
}

// BulkUpload mocks the BulkUpload function of the EsClient interface.
func (m *EsClientMock) BulkUpload(dataUnits []models.DataUnit, indexName string) {
	m.Indexes[indexName] = append(m.Indexes[indexName], dataUnits...)
	m.allDocsCount += len(dataUnits)
}

// CreateEsIndex mocks the CreateEsIndex function of the EsClient interface.
func (m *EsClientMock) CreateEsIndex(index string) {
	m.Indexes[index] = []models.DataUnit{}
}

func NewESClientMock(indexes map[string][]models.DataUnit, expectedCount int) *EsClientMock {
	esMock := EsClientMock{
		Indexes:      indexes,
		allDocsCount: 0,
	}

	return &esMock
}
