package serviceintegrationtests

import (
	"io/ioutil"
	"log"
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

func TestServiceIntegrationWithElasticsearch(t *testing.T) {
	testIndexName := "test"
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
	mockConsumer := mocks.NewRabbitMQConsumerMock(testInputData, allMessagesAcknowledged, testIndexName)
	esClient := elastic.NewEsClientWrapper(testESURL)

	// Start handling messages.
	uploaderService := uploader.NewUploaderService(mockConsumer, esClient, "todo", "todo")
	uploaderService.HandleMessages()

	log.Println(" [TEST] Handling messages...")

	<-allMessagesAcknowledged

	log.Println(" [TEST] All messages handled, waiting for uploading to finish...")

	// We need to wait, because the upload time period is 5 seconds,
	// so to be sure everything is finished uploading, we wait 6 seconds.
	ticker := time.NewTicker(8 * time.Second)
	<-ticker.C

	log.Println(" [TEST] Uploading finished, checking results...")

	// Create a test ES client to query results.
	testESClient := testutils.NewTestEsClientWrapper(testESURL)
	docCount := testESClient.QueryDocCountInIndex(testIndexName)
	testESClient.DeleteIndex(testIndexName) // Clean up test index.

	expectedDocCount := len(testInputData.Events) + len(testInputData.Consumptions)

	if docCount != expectedDocCount {
		t.Fatalf("Expected to have %d documents, actual doc count: %d", expectedDocCount, docCount)
	}
}

func TestServiceIntegrationWithRabbitMQ(t *testing.T) {
	testIndexName := "test"
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
		make(map[string][]models.DataUnit),
		len(testInputData.Consumptions)+len(testInputData.Events))

	uploaderService := uploader.NewUploaderService(rabbitMQConsumer, mockESClient, "todo", "todo")
	uploaderService.HandleMessages()

	testProducer := testutils.NewTestRabbitMqProducer(rabbitMQURL, exchangeName, routingKey)
	testProducer.Connect()
	testProducer.PublishRecreateIndexMessage(testIndexName)

	testProducer.PublishTestInput(testInputData, testIndexName)
	testProducer.PublishDoneMessage()

	log.Println(" [TEST] Waiting for uploading to finish...")
	// We need to wait, because the upload time period is 5 seconds,
	// so to be sure everything is finished uploading, we wait 6 seconds.
	ticker := time.NewTicker(6 * time.Second)
	<-ticker.C

	log.Println(" [TEST] Uploading finished, checking results...")

	if len(mockESClient.Indexes) != 1 {
		t.Fatalf("Expected to create %d indexes, created %d", 1, len(mockESClient.Indexes))
	}

	// cleanup
	testProducer.CloseChannelAndConnection()
	rabbitMQConsumer.CloseChannelAndConnection()
}
