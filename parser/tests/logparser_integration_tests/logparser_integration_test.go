package logparserintegrationtests

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/kozgot/go-log-processing/parser/internal/utils"
	"github.com/kozgot/go-log-processing/parser/pkg/logparser"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/parser/pkg/rabbitmq"
	"github.com/kozgot/go-log-processing/parser/tests/mocks"
	"github.com/kozgot/go-log-processing/parser/tests/testmodels"
	"github.com/kozgot/go-log-processing/parser/tests/testutils"
	"github.com/streadway/amqp"
)

const updateResourcesEnabled = false

// TestLogParserDCMain calls logparser.ParseLogfiles()
// with a mock file downloader that passes a dc main test file and a real rabbitmq producer,
// checking for valid messages.
func TestLogParserDCMain(t *testing.T) {
	rabbitMqURL := "amqp://guest:guest@rabbitmq:5672/"
	testRoutingKey := "test-routing-key"
	testExchangeName := "test_exchange"
	testQueueName := "test_queue"

	rabbitProducer, testConsumer, mockDownloader := setupDependencies(
		rabbitMqURL,
		testRoutingKey,
		testExchangeName,
		testQueueName,
		"./resources/test_dc_main.log")
	defer tearDownDependecies(rabbitProducer, testConsumer)

	// Register test consumer.
	msgs, err := testConsumer.ConsumeMessages()
	utils.FailOnError(err, "Could not register test consumer")

	// Parse
	logParser := logparser.NewLogParser(mockDownloader, rabbitProducer)
	logParser.ParseLogfiles()

	// Get entries sent to RabbitMq.
	entries := getSentParsedEntries(msgs)

	log.Printf("Got %d entries.\n", len(entries))

	testParsedLogFile := testmodels.TestParsedLogFile{Lines: entries}
	actualBytes := testParsedLogFile.ToJSON()

	updateResourcesIfEnabled("./resources/expected_parsed_log.json", actualBytes)

	// Read expected outcome from resource file.
	expectedBytes, err := ioutil.ReadFile("./resources/expected_parsed_log.json")
	utils.FailOnError(err, "Could not read test json file.")

	// Assert
	if string(actualBytes) != string(expectedBytes) {
		t.Fatal("Expected json does not match actual json value of created parsed entries.")
	}
	if len(entries) != 40 {
		t.Fatalf("Expected 40 entries, got %d entries", len(entries))
	}
}

// TestLogParserPLCManager calls logparser.ParseLogfiles()
// with a mock file downloader that passes a dc main test file and a real rabbitmq producer,
// checking for valid messages.
func TestLogParserPLCManager(t *testing.T) {
	rabbitMqURL := "amqp://guest:guest@rabbitmq:5672/"
	testRoutingKey := "test-routing-key"
	testExchangeName := "test_exchange"
	testQueueName := "test_queue"

	rabbitProducer, testConsumer, mockDownloader := setupDependencies(
		rabbitMqURL,
		testRoutingKey,
		testExchangeName,
		testQueueName,
		"./resources/test_plc_manager.log")
	defer tearDownDependecies(rabbitProducer, testConsumer)

	// Register test consumer.
	msgs, err := testConsumer.ConsumeMessages()
	utils.FailOnError(err, "Could not register test consumer")

	// Parse
	logParser := logparser.NewLogParser(mockDownloader, rabbitProducer)
	logParser.ParseLogfiles()

	// Get entries sent to RabbitMq.
	entries := getSentParsedEntries(msgs)

	log.Printf("Got %d entries.\n", len(entries))
	testParsedLogFile := testmodels.TestParsedLogFile{Lines: entries}
	actualBytes := testParsedLogFile.ToJSON()

	updateResourcesIfEnabled("./resources/expected_plc_manager.json", actualBytes)

	// Read expected outcome from resource file.
	expectedBytes, err := ioutil.ReadFile("./resources/expected_plc_manager.json")
	utils.FailOnError(err, "Could not read test json file.")

	// Assert
	if string(actualBytes) != string(expectedBytes) {
		t.Fatal("Expected json does not match actual json value of created parsed entries.")
	}
	if len(entries) != 50 {
		t.Fatalf("Expected 50 entries, got %d entries.", len(entries))
	}
}

func setupDependencies(
	rabbitMqURL string,
	testRoutingKey string,
	testExchangeName string,
	testQueueName string,
	testLogfileName string,
) (
	*rabbitmq.AmqpProducer,
	*testutils.TestRabbitConsumer,
	*mocks.MockFileDownloader,
) {
	// Initialize rabbitMQ producer.
	rabbitMqProducer := rabbitmq.NewAmqpProducer(
		testRoutingKey,
		testExchangeName,
		rabbitMqURL)

	// Init test consumer.
	testConsumer := testutils.NewTestRabbitConsumer(rabbitMqURL, testRoutingKey, testExchangeName, testQueueName)
	testConsumer.Connect()

	// Open a connection and a channel to send the log entries to.
	rabbitMqProducer.OpenChannelAndConnection()

	// Create mock filedownloader.
	mockFileDownloader := mocks.MockFileDownloader{FileNameToDownload: testLogfileName}
	return rabbitMqProducer, testConsumer, &mockFileDownloader
}

func tearDownDependecies(
	rabbitMqProducer *rabbitmq.AmqpProducer,
	testConsumer *testutils.TestRabbitConsumer,
) {
	testConsumer.CloseConnectionAndChannel()
	rabbitMqProducer.CloseChannelAndConnection()
}

func getSentParsedEntries(deliveries <-chan amqp.Delivery) []models.ParsedLogEntry {
	entries := []models.ParsedLogEntry{}

	for d := range deliveries {
		if strings.Contains(string(d.Body), "END") {
			log.Println("End of entries...")

			// Acknowledge the message after it has been processed.
			err := d.Ack(false)
			utils.FailOnError(err, "Could not acknowledge END message")
			break
		} else if strings.Contains(string(d.Body), "START") {
			// Acknowledge the message after it has been processed.
			err := d.Ack(false)
			utils.FailOnError(err, "Could not acknowledge START message")
			continue
		}
		entry := models.ParsedLogEntry{}
		entry.FromJSON(d.Body)
		entries = append(entries, entry)
		err := d.Ack(false)
		utils.FailOnError(err, "Could not acknowledge")
	}

	return entries
}

func updateResourcesIfEnabled(resourceFileName string, newData []byte) {
	if updateResourcesEnabled {
		_ = ioutil.WriteFile(resourceFileName, newData, 0600)
	}
}
