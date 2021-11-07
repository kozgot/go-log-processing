package logparserintegrationtests

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/kozgot/go-log-processing/parser/internal/utils"
	"github.com/kozgot/go-log-processing/parser/pkg/logparser"
	"github.com/kozgot/go-log-processing/parser/pkg/mocks"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/parser/pkg/rabbitmq"
	"github.com/kozgot/go-log-processing/parser/tests/testmodels"
	"github.com/kozgot/go-log-processing/parser/tests/testutils"
)

// TestLogParserDCMain calls logparser.ParseLogfiles()
// with a mock file downloader that passes a dc main test file and a real rabbitmq producer,
// checking for valid messages.
func TestLogParserDCMain(t *testing.T) {
	rabbitMqURL := "amqp://guest:guest@rabbitmq:5672/"
	testRoutingKey := "test-routing-key"
	testExchangeName := "test_exchange"
	testQueueName := "test_queue"

	// Initialize rabbitMQ producer.
	rabbitMqProducer := rabbitmq.NewAmqpProducer(
		testRoutingKey,
		testExchangeName,
		rabbitMqURL)

	// Open a connection and a channel to send the log entries to.
	rabbitMqProducer.OpenChannelAndConnection()
	defer rabbitMqProducer.CloseChannelAndConnection()

	// Init test consumer.
	testConsumer := testutils.NewTestRabbitConsumer(rabbitMqURL, testRoutingKey, testExchangeName, testQueueName)
	testConsumer.Connect()
	defer testConsumer.CloseConnectionAndChannel()

	// Register test consumer.
	msgs, err := testConsumer.ConsumeMessages()
	utils.FailOnError(err, "Could not register test consumer")

	// Create mock filedownloader.
	mockFileDownloader := mocks.MockFileDownloader{FileNameToDownload: "./resources/test_dc_main.log"}

	// Parse
	logParser := logparser.NewLogParser(&mockFileDownloader, rabbitMqProducer)
	logParser.ParseLogfiles()

	entries := []models.ParsedLogEntry{}

	for d := range msgs {
		if strings.Contains(string(d.Body), "END") {
			log.Println("End of entries...")

			// Acknowledge the message after it has been processed.
			err := d.Ack(false)
			utils.FailOnError(err, "Could not acknowledge END message")
			break
		}
		entry := models.ParsedLogEntry{}
		entry.FromJSON(d.Body)
		entries = append(entries, entry)
		err := d.Ack(false)
		utils.FailOnError(err, "Could not acknowledge")
	}

	log.Printf("Got %d entries.\n", len(entries))
	testParsedLogFile := testmodels.TestParsedLogFile{Lines: entries}
	actualBytes := testParsedLogFile.ToJSON()
	expectedBytes, err := ioutil.ReadFile("./resources/expected_parsed_log.json")
	utils.FailOnError(err, "Could not read test json file.")

	if string(actualBytes) != string(expectedBytes) {
		t.Fatal("Expected json does not match actual json value of created partsed entries.")
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

	// Initialize rabbitMQ producer.
	rabbitMqProducer := rabbitmq.NewAmqpProducer(
		testRoutingKey,
		testExchangeName,
		rabbitMqURL)

	// Open a connection and a channel to send the log entries to.
	rabbitMqProducer.OpenChannelAndConnection()
	defer rabbitMqProducer.CloseChannelAndConnection()

	// Init test consumer.
	testConsumer := testutils.NewTestRabbitConsumer(rabbitMqURL, testRoutingKey, testExchangeName, testQueueName)
	testConsumer.Connect()
	defer testConsumer.CloseConnectionAndChannel()

	// Register test consumer.
	msgs, err := testConsumer.ConsumeMessages()
	utils.FailOnError(err, "Could not register test consumer")

	// Create mock filedownloader.
	mockFileDownloader := mocks.MockFileDownloader{FileNameToDownload: "./resources/test_plc_manager.log"}

	// Parse
	logParser := logparser.NewLogParser(&mockFileDownloader, rabbitMqProducer)
	logParser.ParseLogfiles()

	entries := []models.ParsedLogEntry{}

	for d := range msgs {
		if strings.Contains(string(d.Body), "END") {
			log.Println("End of entries...")

			// Acknowledge the message after it has been processed.
			err := d.Ack(false)
			utils.FailOnError(err, "Could not acknowledge END message")
			break
		}
		entry := models.ParsedLogEntry{}
		entry.FromJSON(d.Body)
		entries = append(entries, entry)
		err := d.Ack(false)
		utils.FailOnError(err, "Could not acknowledge")
	}

	log.Printf("Got %d entries.\n", len(entries))
	testParsedLogFile := testmodels.TestParsedLogFile{Lines: entries}
	actualBytes := testParsedLogFile.ToJSON()
	expectedBytes, err := ioutil.ReadFile("./resources/expected_plc_manager.json")
	utils.FailOnError(err, "Could not read test json file.")

	if string(actualBytes) != string(expectedBytes) {
		t.Fatal("Expected json does not match actual json value of created partsed entries.")
	}
	if len(entries) != 50 {
		t.Fatalf("Expected 50 entries, got %d entries.", len(entries))
	}
}
