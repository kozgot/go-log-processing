package mocks

import "github.com/kozgot/go-log-processing/elasticuploader/pkg/models"

// EsClientMock is a mock Elasticsearch client used in tests.
type EsClientMock struct {
	Indexes map[string][]models.DataUnit
}

// BulkUpload mocks the BulkUpload function of the EsClient interface.
func (m *EsClientMock) BulkUpload(dataUnits []models.DataUnit, indexName string) {
	/*
		_, ok := m.Indexes[indexName]
		if !ok {
			m.Indexes[indexName] = []models.DataUnit{}
		}*/

	m.Indexes[indexName] = append(m.Indexes[indexName], dataUnits...)
}

// RecreateEsIndex mocks the RecreateEsIndex function of the EsClient interface.
func (m *EsClientMock) RecreateEsIndex(index string) {
	m.Indexes[index] = []models.DataUnit{}
}
