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

// TestLogParser calls logparser.ParseLogfiles()
// with a mock file downloader that passes a dc main test file and a real rabbitmq producer,
// checking for valid messages.
func TestLogParser(t *testing.T) {
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
	mockFileDownloader := mocks.MockFileDownloader{FileNameToDownload: "../resources/test_dc_main.log"}

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
	_ = ioutil.WriteFile("test.json", testParsedLogFile.ToJSON(), 0644)

	if len(entries) != 40 {
		t.Fatal("Expected some entries")
	}
}
