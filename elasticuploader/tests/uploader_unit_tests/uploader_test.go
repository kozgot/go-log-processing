package uploaderunittests

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/uploader"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/mocks"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/testmodels"
)

// TestUploderService tests the uploader service
// with a large index recreation time period.
func TestUploderService(t *testing.T) {
	inputFileName := "./resources/input_data.json"

	// Read test input from resource file.
	testInput, err := ioutil.ReadFile(inputFileName)
	utils.FailOnError(err, "Could not read file "+inputFileName)
	testInputData := testmodels.TestProcessedData{}
	testInputData.FromJSON(testInput)

	expectedDocCount := len(testInputData.Events) + len(testInputData.Consumptions)

	// A channel to indicate that all mocked deliveries are acknowledged (handled).
	allMessagesAcknowledged := make(chan bool)

	// Create a mock rabbitMQ consumer.
	mockConsumer := mocks.NewRabbitMQConsumerMock(testInputData, allMessagesAcknowledged, 0, expectedDocCount)

	// Create a mock ES client.
	mockESClient := mocks.NewESClientMock(
		make(map[string][]models.ESDocument),
		expectedDocCount)

	uploaderService := uploader.NewUploaderService(
		mockConsumer,
		mockESClient,
		"test_events",       // event index name
		"test_consumptions", // consumption index name
		"@midnight",         // index recreation time, in a non-test environment it would be every midnight
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

	if len(mockESClient.Indexes) != 2 {
		t.Fatalf("Expected to create %d indexes, created %d", 2, len(mockESClient.Indexes))
	}

	for key, data := range mockESClient.Indexes {
		if strings.Contains(key, "test_events") && len(data) != 23 {
			t.Fatalf(
				"Expected to have %d documents in the events index, actual doc count %d",
				23,
				len(data),
			)
		}

		if strings.Contains(key, "test_consumptions") && len(data) != 0 {
			t.Fatalf(
				"Expected to have %d documents in the consumptions index, actual doc count %d",
				0,
				len(data),
			)
		}
	}
}
