package processingunittests

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/kozgot/go-log-processing/postprocessor/internal/processing"
	"github.com/kozgot/go-log-processing/postprocessor/internal/utils"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/tests/mocks"
	"github.com/kozgot/go-log-processing/postprocessor/tests/testmodels"
)

const updateResourcesEnabled = false

type postProcessorTest struct {
	inputDataFile            string
	expectedDataFile         string
	expectedEventCount       int
	expectedConsumptionCount int
}

func TestProcessEntries(t *testing.T) {
	postProcessorTests := []postProcessorTest{
		{
			inputDataFile:            "./resources/parsed_test_dc_main.json",
			expectedDataFile:         "./resources/expected_processed_dc_main.json",
			expectedEventCount:       23,
			expectedConsumptionCount: 0,
		},
		{
			inputDataFile:            "./resources/parsed_test_plc_manager.json",
			expectedDataFile:         "./resources/expected_processed_plc_manager.json",
			expectedEventCount:       19,
			expectedConsumptionCount: 0,
		},
	}

	for index, test := range postProcessorTests {
		done := make(chan string)

		// Init a mock message producer.
		mockMessageProducer := mocks.NewMockMessageProducer(
			testmodels.TestProcessedData{
				Events:       []models.SmcEvent{},
				Consumptions: []models.ConsumtionValue{},
			},
			done,
			test.expectedEventCount+test.expectedConsumptionCount,
		)

		// Read test input from resource file.
		parsedInputBytes, err := ioutil.ReadFile(test.inputDataFile)
		utils.FailOnError(err, "Could not open test input "+test.inputDataFile)

		testData := testmodels.TestParsedLogFile{}
		testData.FromJSON(parsedInputBytes)
		mockMessageConsumer := mocks.MockMessageConsumer{TestParsedLogFile: testData}

		// Run processor
		processor := processing.NewEntryProcessor(
			mockMessageProducer,
			&mockMessageConsumer,
		)
		processor.HandleEntries()

		<-done

		actualProcessedDataBytes := mockMessageProducer.Data.ToJSON()
		updateResourcesIfEnabled(test.expectedDataFile, actualProcessedDataBytes)

		// Read expected outcome from resource file.
		expectedBytes, err := ioutil.ReadFile(test.expectedDataFile)
		utils.FailOnError(err, "Could not read file "+test.expectedDataFile)

		// Assert
		if string(actualProcessedDataBytes) != string(expectedBytes) {
			t.Fatalf("Expected json does not match actual json value of processed data in test number %d", index)
		}
	}

	log.Printf("Successfully ran %d test cases", len(postProcessorTests))
}

func updateResourcesIfEnabled(resourceFileName string, newData []byte) {
	if updateResourcesEnabled {
		err := ioutil.WriteFile(resourceFileName, newData, 0600)
		utils.FailOnError(err, "Could not update resource file "+resourceFileName)
	}
}
