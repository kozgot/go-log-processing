package serviceunittests

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/uploader"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/mocks"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/testmodels"
)

func TestUploderService(t *testing.T) {
	testIndexName := "test"
	inputFileName := "./resources/input_data.json"

	// Read test input from resource file.
	testInput, err := ioutil.ReadFile(inputFileName)
	utils.FailOnError(err, "Could not read file "+inputFileName)
	testInputData := testmodels.TestProcessedData{}
	testInputData.FromJSON(testInput)

	expectedDocCount := len(testInputData.Events) + len(testInputData.Consumptions)

	// A channel to indicate that all mocked deliveries are acknowledged (handled).
	allMessagesAcknowledged := make(chan bool)

	// Create a mock rabbitMQ consumer and actual ES client as dependencies.
	mockConsumer := mocks.NewRabbitMQConsumerMock(testInputData, allMessagesAcknowledged, testIndexName)

	mockESClient := mocks.NewESClientMock(
		make(map[string][]models.DataUnit),
		expectedDocCount)

	uploaderService := uploader.NewUploaderService(
		mockConsumer,
		mockESClient,
		"todo",
		"todo",
	)
	uploaderService.HandleMessages()

	log.Println(" [TEST] Handling messages...")

	<-allMessagesAcknowledged

	log.Println(" [TEST] All messages handled, waiting for uploading to finish...")

	// We need to wait, because the upload time period is 5 seconds,
	// so to be sure everything is finished uploading, we wait 6 seconds.
	ticker := time.NewTicker(6 * time.Second)
	<-ticker.C

	log.Println(" [TEST] Uploading finished, checking results...")

	if len(mockESClient.Indexes) != 1 {
		t.Fatalf("Expected to create %d indexes, created %d", 1, len(mockESClient.Indexes))
	}

	if len(mockESClient.Indexes[testIndexName]) != expectedDocCount {
		t.Fatalf("Expected doc count %d, actual doc count %d",
			expectedDocCount,
			len(mockESClient.Indexes[testIndexName]))
	}
}
