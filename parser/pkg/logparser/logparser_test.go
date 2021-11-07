package logparser

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/kozgot/go-log-processing/parser/internal/utils"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/parser/pkg/rabbitmq"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	name := "Gladys"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Hello("Gladys")
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}

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
	testConsumer := NewTestRabbitConsumer(rabbitMqURL, testRoutingKey, testExchangeName, testQueueName)
	testConsumer.Connect()
	defer testConsumer.CloseConnectionAndChannel()

	// Register test consumer.
	msgs, err := testConsumer.ConsumeMessages()
	utils.FailOnError(err, "Could not register test consumer")

	// Create mock filedownloader.
	mockFileDownloader := MockFileDownloader{FileNameToDownload: "test_dc_main.log"}

	// Parse
	logParser := NewLogParser(&mockFileDownloader, rabbitMqProducer)
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
		entry := deserializeParsedLogEntry(d.Body)
		entries = append(entries, entry)
		err := d.Ack(false)
		utils.FailOnError(err, "Could not acknowledge")
	}

	log.Printf("Got %d entries.\n", len(entries))

	if len(entries) == 0 {
		t.Fatal("Expected some entries")
	}
}

func deserializeParsedLogEntry(bytes []byte) models.ParsedLogEntry {
	var parsedEntry models.ParsedLogEntry
	err := json.Unmarshal(bytes, &parsedEntry)
	utils.FailOnError(err, "Failed to unmarshal log entry")
	return parsedEntry
}
