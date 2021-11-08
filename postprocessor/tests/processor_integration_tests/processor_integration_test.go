package processorintegrationtests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/kozgot/go-log-processing/postprocessor/internal/processing"
	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"
	"github.com/kozgot/go-log-processing/postprocessor/tests/testmodels"
	"github.com/kozgot/go-log-processing/postprocessor/tests/testutils"
	"github.com/streadway/amqp"
)

const updateResourcesEnabled = false

// TestLogParserDCMain calls logparser.ParseLogfiles()
// with a mock file downloader that passes a dc main test file and a real rabbitmq producer,
// checking for valid messages.
func TestLogParserDCMain(t *testing.T) {
	rabbitMqURL := "amqp://guest:guest@rabbitmq:5672/"
	testInputRoutingKey := "test-input-routing-key"
	testInputExchangeName := "test_input_exchange"
	testInputQueueName := "test_input_queue"

	testOutputRoutingKey := "test-output-routing-key"
	testOutputExchangeName := "test_output_exchange"
	testOutputQueueName := "test_output_queue"

	testEventIdxName := "smctest"
	testConsumptionIdxName := "consumptiontest"

	rabbitMqOutputProducer, rabbitMqInputConsumer, testOutputConsumer, testInputProducer := setupDependencies(
		rabbitMqURL,
		testInputRoutingKey,
		testInputExchangeName,
		testInputQueueName,
		testOutputRoutingKey,
		testOutputExchangeName,
		testOutputQueueName,
	)

	defer tearDownDependecies(
		rabbitMqOutputProducer,
		rabbitMqInputConsumer,
		testOutputConsumer,
		testInputProducer,
	)

	msgs := testOutputConsumer.ConsumeMessages()

	processor := processing.NewEntryProcessor(
		rabbitMqOutputProducer,
		rabbitMqInputConsumer,
		testEventIdxName,
		testConsumptionIdxName)
	processor.HandleEntries()

	// todo produce input messages
	// Read expected outcome from resource file.
	parsedInputBytes, err := ioutil.ReadFile("./resources/parsed_test_dc_main.json")
	utils.FailOnError(err, "Could not open test input ./resources/parsed_test_dc_main.json")
	testparsedFile := testmodels.TestParsedLogFile{}
	testparsedFile.FromJSON(parsedInputBytes)
	sendTestInput(testInputProducer, testparsedFile)

	processedData := getSentProcessedData(msgs, testEventIdxName, testConsumptionIdxName)
	if len(processedData.IndexNames) != 2 {
		t.Fatal("Expected to create 2 indices")
	}

	actualProcessedDataBytes := processedData.ToJSON()
	updateResourcesIfEnabled("./resources/expected_processed_dc_main.json", actualProcessedDataBytes)

	// Read expected outcome from resource file.
	expectedBytes, err := ioutil.ReadFile("./resources/expected_processed_dc_main.json")
	utils.FailOnError(err, "Could not read expected json file.")

	// Assert
	if string(actualProcessedDataBytes) != string(expectedBytes) {
		t.Fatal("Expected json does not match actual json value of processed data.")
	}

	// todo more assertions
}

func setupDependencies(
	rabbitMqURL string,
	testInputRoutingKey string,
	testInputExchangeName string,
	testInputQueueName string,
	testOutputRoutingKey string,
	testOutputExchangeName string,
	testOutputQueueName string,
) (
	*rabbitmq.AmqpProducer,
	*rabbitmq.AmqpConsumer,
	*testutils.TestRabbitConsumer,
	*testutils.TestRabbitMqProducer,
) {
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
		case "CREATEINDEX":
			indexName := strings.Split(string(delivery.Body), "|")[1]
			testdata.IndexNames = append(testdata.IndexNames, indexName)
			// service.esClient.CreateEsIndex(indexName)
		default:
			data := DeserializeDataUnit(delivery.Body)
			switch data.IndexName {
			case testEventIdxName:
				smcEvent := models.SmcEvent{}
				smcEvent.Deserialize(data.Data)
				testdata.Events = append(testdata.Events, smcEvent)
			case testConsumptionIdxName:
				consumption := models.ConsumtionValue{}
				consumption.Deserialize(data.Data)
				testdata.Consumptions = append(testdata.Consumptions, consumption)
			}
			// uploadBuffer.AppendAndUploadIfNeeded(models.Message{Content: data.Data}, data.IndexName, uploadTicker)
		}

		// Acknowledge message
		err := delivery.Ack(false)
		utils.FailOnError(err, "Could not acknowledge message")
	}

	return testdata
}

func updateResourcesIfEnabled(resourceFileName string, newData []byte) {
	if updateResourcesEnabled {
		_ = ioutil.WriteFile(resourceFileName, newData, 0600)
	}
}

// DeserializeDataUnit deserializes a received data unit.
func DeserializeDataUnit(dataBytes []byte) models.DataUnit {
	var data models.DataUnit
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		fmt.Println("Failed to unmarshal:", err)
	}

	return data
}

func sendTestInput(
	testInputProducer *testutils.TestRabbitMqProducer,
	testparsedFile testmodels.TestParsedLogFile) {
	for _, parsedEntry := range testparsedFile.Lines {
		testInputProducer.PublishEntry(parsedEntry)
	}

	// Send a message indicating that this is the end of the processing
	testInputProducer.PublishStringMessage("END")
}
