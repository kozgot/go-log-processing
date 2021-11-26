package uploaderintegrationtests

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/kozgot/go-log-processing/elasticuploader/internal/elastic"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/rabbit"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/uploader"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/mocks"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/testmodels"
	"github.com/kozgot/go-log-processing/elasticuploader/tests/testutils"
)

// TestServiceIntegrationWithElasticsearch uses a real ES client to upload data consumed from a mock RabbitMQ consumer.
// Expected to create exactly 2 indexes, with only the events index containing any documents.
func TestServiceIntegrationWithElasticsearch(t *testing.T) {
	inputFileName := "./resources/input_data.json"
	testESURL := "http://elasticsearch:9200"

	// Read test input from resource file.
	testInput, err := ioutil.ReadFile(inputFileName)
	utils.FailOnError(err, "Could not read file "+inputFileName)
	testInputData := testmodels.TestProcessedData{}
	testInputData.FromJSON(testInput)

	// A channel to indicate that all mocked deliveries are acknowledged (handled).
	allMessagesAcknowledged := make(chan bool)

	// Create a mock rabbitMQ consumer and actual ES client as dependencies.
	mockConsumer := mocks.NewRabbitMQConsumerMock(
		testInputData,
		allMessagesAcknowledged,
		0,
		len(testInputData.Consumptions)+len(testInputData.Events),
	)
	esClient := elastic.NewEsClientWrapper(testESURL)

	// Start handling messages.
	uploaderService := uploader.NewUploaderService(
		mockConsumer,
		esClient,
		"test_events",
		"test_consumptions",
		"@midnight", // todo
	)
	uploaderService.HandleMessages()

	log.Println(" [TEST] Handling messages...")

	<-allMessagesAcknowledged

	log.Println(" [TEST] All messages handled, waiting for uploading to finish...")

	// We need to wait, because the upload time period is 5 seconds,
	// so to be sure everything is finished uploading, we wait 6 seconds.
	ticker := time.NewTicker(8 * time.Second)
	<-ticker.C

	log.Println(" [TEST] Uploading finished, checking results...")

	if len(esClient.CreatedIndexNames) != 2 {
		t.Fatalf("Expected to create 2 indexes, actual index count: %d", len(esClient.CreatedIndexNames))
	}

	// Create a test ES client to query results.
	testESClient := testutils.NewTestEsClientWrapper(testESURL)

	// Check event documents.
	docCount := testESClient.QueryDocCountInIndex(esClient.CreatedIndexNames[0])
	testESClient.DeleteIndex(esClient.CreatedIndexNames[0]) // Clean up test index.
	expectedDocCount := len(testInputData.Events) + len(testInputData.Consumptions)
	if docCount != expectedDocCount {
		t.Fatalf("Expected to have %d documents, actual doc count: %d", expectedDocCount, docCount)
	}

	// Check consumption documents.
	consumptionDocCount := testESClient.QueryDocCountInIndex(esClient.CreatedIndexNames[1])
	testESClient.DeleteIndex(esClient.CreatedIndexNames[1]) // Clean up test index.
	if consumptionDocCount != 0 {
		t.Fatalf("Expected to have %d documents in the consumption index, actual doc count: %d",
			0,
			consumptionDocCount,
		)
	}
}

func TestServiceIntegrationWithRabbitMQ(t *testing.T) {
	inputFileName := "./resources/input_data.json"
	rabbitMQURL := "amqp://guest:guest@rabbitmq:5672/"
	exchangeName := "test_exchange"
	routingKey := "test-key"
	queueName := "test_queue"

	// Read test input from resource file.
	testInput, err := ioutil.ReadFile(inputFileName)
	utils.FailOnError(err, "Could not read file "+inputFileName)
	testInputData := testmodels.TestProcessedData{}
	testInputData.FromJSON(testInput)

	rabbitMQConsumer := rabbit.NewAmqpConsumer(rabbitMQURL, exchangeName, routingKey, queueName)
	rabbitMQConsumer.Connect()

	mockESClient := mocks.NewESClientMock(
		make(map[string][]models.ESDocument),
		len(testInputData.Consumptions)+len(testInputData.Events))

	uploaderService := uploader.NewUploaderService(
		rabbitMQConsumer,
		mockESClient,
		"test_events",
		"test_consumptions",
		"@midnight", // todo
	)

	uploaderService.HandleMessages()

	testProducer := testutils.NewTestRabbitMqProducer(rabbitMQURL, exchangeName, routingKey)
	testProducer.Connect()

	testProducer.PublishTestInput(testInputData)

	log.Println(" [TEST] Waiting for uploading to finish...")
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

	// cleanup
	testProducer.CloseChannelAndConnection()
	rabbitMQConsumer.CloseChannelAndConnection()
}
