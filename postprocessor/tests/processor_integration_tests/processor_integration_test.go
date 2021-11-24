package processorintegrationtests

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/kozgot/go-log-processing/postprocessor/internal/processing"
	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/postprocessor/internal/utils"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/tests/testmodels"
	"github.com/kozgot/go-log-processing/postprocessor/tests/testutils"
	"github.com/streadway/amqp"
)

// If set to true, running the tests automatically updates the expeted resource files.
const updateResourcesEnabled = true

// TestProcessDCMain calls processor.HandleEntries()
// with a real rabbitM message consumer that consumes parsed log entries from a dc_main.log file.
func TestProcessDCMain(t *testing.T) {
	testEventIdxName := "smctest"
	testConsumptionIdxName := "consumptiontest"

	testInputFileName := "./resources/parsed_test_dc_main.json"
	expectedDataFileName := "./resources/expected_processed_dc_main.json"

	rabbitMqOutputProducer, rabbitMqInputConsumer, testOutputConsumer, testInputProducer := setupDependencies()

	defer tearDownDependecies(
		rabbitMqOutputProducer,
		rabbitMqInputConsumer,
		testOutputConsumer,
		testInputProducer,
	)

	// Register a consumer for the output of the processing.
	msgs := testOutputConsumer.ConsumeMessages()

	processor := processing.NewEntryProcessor(
		rabbitMqOutputProducer,
		rabbitMqInputConsumer,
		testEventIdxName,
		testConsumptionIdxName)
	processor.HandleEntries()

	// Read test input from resource file.
	parsedInputBytes, err := ioutil.ReadFile(testInputFileName)
	utils.FailOnError(err, "Could not open test input "+testInputFileName)

	testparsedFile := testmodels.TestParsedLogFile{}
	testparsedFile.FromJSON(parsedInputBytes)

	// Create test input for the processor.
	sendTestInput(testInputProducer, testparsedFile)

	// Handle output created by the processor.
	processedData := getSentProcessedData(msgs, testEventIdxName, testConsumptionIdxName)
	if len(processedData.IndexNames) != 2 {
		t.Fatalf("Expected to create 2 indices, actual: %d", len(processedData.IndexNames))
	}

	actualProcessedDataBytes := processedData.ToJSON()
	updateResourcesIfEnabled(expectedDataFileName, actualProcessedDataBytes)

	// Read expected outcome from resource file.
	expectedBytes, err := ioutil.ReadFile(expectedDataFileName)
	utils.FailOnError(err, "Could not read file "+expectedDataFileName)

	// Assert
	if string(actualProcessedDataBytes) != string(expectedBytes) {
		t.Fatal("Expected json does not match actual json value of processed data.")
	}
}

// TestProcessPLCManager calls processor.HandleEntries()
// with a real rabbitM message consumer that consumes parsed log entries from a plc-manager.log file.
func TestProcessPLCManager(t *testing.T) {
	testEventIdxName := "smctest"
	testConsumptionIdxName := "consumptiontest"

	testInputFileName := "./resources/parsed_test_plc_manager.json"
	expectedDataFileName := "./resources/expected_processed_plc_manager.json"

	rabbitMqOutputProducer, rabbitMqInputConsumer, testOutputConsumer, testInputProducer := setupDependencies()

	defer tearDownDependecies(
		rabbitMqOutputProducer,
		rabbitMqInputConsumer,
		testOutputConsumer,
		testInputProducer,
	)

	// Register a consumer for the output of the processing.
	msgs := testOutputConsumer.ConsumeMessages()

	processor := processing.NewEntryProcessor(
		rabbitMqOutputProducer,
		rabbitMqInputConsumer,
		testEventIdxName,
		testConsumptionIdxName)
	processor.HandleEntries()

	// Read test input from resource file.
	parsedInputBytes, err := ioutil.ReadFile(testInputFileName)
	utils.FailOnError(err, "Could not open test input "+testInputFileName)

	testparsedFile := testmodels.TestParsedLogFile{}
	testparsedFile.FromJSON(parsedInputBytes)

	// Create test input for the processor.
	sendTestInput(testInputProducer, testparsedFile)

	// Handle output created by the processor.
	processedData := getSentProcessedData(msgs, testEventIdxName, testConsumptionIdxName)
	if len(processedData.IndexNames) != 2 {
		t.Fatalf("Expected to create 2 indices, got %d", len(processedData.IndexNames))
	}

	actualProcessedDataBytes := processedData.ToJSON()
	updateResourcesIfEnabled(expectedDataFileName, actualProcessedDataBytes)

	// Read expected outcome from resource file.
	expectedBytes, err := ioutil.ReadFile(expectedDataFileName)
	utils.FailOnError(err, "Could not read file "+expectedDataFileName)

	// Assert
	if string(actualProcessedDataBytes) != string(expectedBytes) {
		t.Fatal("Expected json does not match actual json value of processed data.")
	}
}

func setupDependencies() (
	*rabbitmq.AmqpProducer,
	*rabbitmq.AmqpConsumer,
	*testutils.TestRabbitConsumer,
	*testutils.TestRabbitMqProducer,
) {
	rabbitMqURL := "amqp://guest:guest@rabbitmq:5672/"
	testInputRoutingKey := "test-input-routing-key"
	testInputExchangeName := "test_input_exchange"
	testInputQueueName := "test_input_queue"

	testOutputRoutingKey := "test-output-routing-key"
	testOutputExchangeName := "test_output_exchange"
	testOutputQueueName := "test_output_queue"

	// Initialize rabbitMQ producer. This will be passed to the processor as a parameter to send processed data to.
	rabbitMqOutputProducer := rabbitmq.NewAmqpProducer(
		rabbitMqURL,
		testOutputExchangeName,
		testOutputRoutingKey,
	)

	// Initialize rabbitMQ consumer. This will be passed to the processor as a parameter to consume parsed entries from.
	rabbitMqInputConsumer := rabbitmq.NewAmqpConsumer(
		rabbitMqURL,
		testInputRoutingKey,
		testInputExchangeName,
		testInputQueueName,
	)

	// Init test consumer. This will be used to consume messages from the output of the processor to validate the output.
	testOutputConsumer := testutils.NewTestRabbitConsumer(
		rabbitMqURL,
		testOutputRoutingKey,
		testOutputExchangeName,
		testOutputQueueName)
	testOutputConsumer.Connect()

	// Init test producer. This will be used to produce messages as input for the processor.
	testInputProducer := testutils.NewTestProducer(rabbitMqURL, testInputExchangeName, testInputRoutingKey)
	testInputProducer.Connect()

	// Open a connection and a channel to send processed data to.
	rabbitMqOutputProducer.Connect()

	// Open a connection and a channel to consume parsed entries from.
	rabbitMqInputConsumer.Connect()

	// Create mock filedownloader.
	return rabbitMqOutputProducer, rabbitMqInputConsumer, testOutputConsumer, testInputProducer
}

func tearDownDependecies(
	rabbitMqOutputProducer *rabbitmq.AmqpProducer,
	rabbitMqInputConsumer *rabbitmq.AmqpConsumer,
	testOutputConsumer *testutils.TestRabbitConsumer,
	testInputProducer *testutils.TestRabbitMqProducer,
) {
	rabbitMqOutputProducer.CloseChannelAndConnection()
	rabbitMqInputConsumer.CloseConnectionAndChannel()
	testOutputConsumer.CloseConnectionAndChannel()
	testInputProducer.CloseChannelAndConnection()
}

// getSentProcessedData reads and returns the processed data sent to a rabbitMQ exchange by the postprocessor.
func getSentProcessedData(
	deliveries <-chan amqp.Delivery,
	testEventIdxName string,
	testConsumptionIdxName string) testmodels.TestProcessedData {
	testdata := testmodels.TestProcessedData{
		IndexNames:   []string{},
		Events:       []models.SmcEvent{},
		Consumptions: []models.ConsumtionValue{},
	}
	for delivery := range deliveries {
		msgParts := strings.Split(string(delivery.Body), "|")
		msgPrefix := msgParts[0]
		switch msgPrefix {
		case "DONE":
			return testdata
		case "RECREATEINDEX":
			indexName := strings.Split(string(delivery.Body), "|")[1]
			testdata.IndexNames = append(testdata.IndexNames, indexName)
		default:
			dataUnit := models.DataUnit{}
			dataUnit.Deserialize(delivery.Body)
			switch dataUnit.IndexName {
			case testEventIdxName:
				smcEvent := models.SmcEvent{}
				smcEvent.Deserialize(dataUnit.Data)
				testdata.Events = append(testdata.Events, smcEvent)
			case testConsumptionIdxName:
				consumption := models.ConsumtionValue{}
				consumption.Deserialize(dataUnit.Data)
				testdata.Consumptions = append(testdata.Consumptions, consumption)
			}
		}

		// Acknowledge message
		err := delivery.Ack(false)
		utils.FailOnError(err, "Could not acknowledge message")
	}

	return testdata
}

// sendTestInput publishes test parsed log entries to a rabbitMQ exchange for the processor to consume.
func sendTestInput(
	testInputProducer *testutils.TestRabbitMqProducer,
	testparsedFile testmodels.TestParsedLogFile) {
	// Send a message indicating that this is the start of the entries.
	testInputProducer.PublishStringMessage("START")

	for _, parsedEntry := range testparsedFile.Lines {
		testInputProducer.PublishEntry(parsedEntry)
	}

	// Send a message indicating that this is the end of the entries.
	testInputProducer.PublishStringMessage("END")
}

func updateResourcesIfEnabled(resourceFileName string, newData []byte) {
	if updateResourcesEnabled {
		_ = ioutil.WriteFile(resourceFileName, newData, 0600)
	}
}
