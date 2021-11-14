package serviceintegrationtests

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/elastic"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/service"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/mocks"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/testmodels"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/testutils"
)

func TestServiceIntegrationWithElasticsearch(t *testing.T) {
	testIndexName := "test"
	inputFileName := "./resources/input_data.json"

	// Read test input from resource file.
	testInput, err := ioutil.ReadFile(inputFileName)
	utils.FailOnError(err, "Could not read file "+inputFileName)
	testInputData := testmodels.TestProcessedData{}
	testInputData.FromJSON(testInput)

	// A channel to indicate that all mocked deliveries are acknowledged (handled).
	allMessagesAcknowledged := make(chan bool)

	// Create a mock rabbitMQ consumer and actual ES client as dependencies.
	mockConsumer := mocks.NewRabbitMQConsumerMock(testInputData, allMessagesAcknowledged, testIndexName)
	esClient := elastic.NewEsClientWrapper()

	// Start handling messages.
	uploaderService := service.NewUploaderService(mockConsumer, esClient)
	uploaderService.HandleMessages()

	fmt.Println("Handling messages...")

	<-allMessagesAcknowledged

	fmt.Println("All messages handled, waiting for uploading to finish...")

	// We need to wait, because the uploading period is 5 seconds.
	ticker := time.NewTicker(6 * time.Second)
	<-ticker.C

	fmt.Println("Uploading finished, checking results...")

	// Create a test ES client to query results.
	testESClient := testutils.NewTestEsClientWrapper()
	docCount := testESClient.QueryDocCountInIndex(testIndexName)
	testESClient.DeleteIndex(testIndexName) // Clean up test index.

	expectedDocCount := len(testInputData.Events) + len(testInputData.Consumptions)

	if docCount != expectedDocCount {
		t.Fatalf("Expected to have %d documents, actual doc count: %d", expectedDocCount, docCount)
	}
}

func TestServiceIntegrationWithRabbitMQ(t *testing.T) {

}
