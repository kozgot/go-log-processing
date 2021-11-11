package processorunittests

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/kozgot/go-log-processing/postprocessor/internal/processing"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/utils"
	"github.com/kozgot/go-log-processing/postprocessor/tests/mocks"
	"github.com/kozgot/go-log-processing/postprocessor/tests/testmodels"
)

const updateResourcesEnabled = false

func TestLogParserPLCManager(t *testing.T) {
	smcEventIndexName := "testSmcEvents"
	consumtionIndexName := "testConsumptions"

	testInputFileName := "./resources/parsed_test_dc_main.json"
	expectedDataFileName := "./resources/expected_processed_dc_main.json"

	done := make(chan string)

	// Init a mock message producer.
	mockMessageProducer := mocks.MockMessageProducer{
		Data: testmodels.TestProcessedData{
			Events:       []models.SmcEvent{},
			Consumptions: []models.ConsumtionValue{},
			IndexNames:   []string{},
		},
		Done: done,
	}

	// Read test input from resource file.
	parsedInputBytes, err := ioutil.ReadFile(testInputFileName)
	utils.FailOnError(err, "Could not open test input "+testInputFileName)

	testData := testmodels.TestParsedLogFile{}
	testData.FromJSON(parsedInputBytes)
	mockMessageConsumer := mocks.MockMessageConsumer{TestParsedLogFile: testData}

	// Run processor
	processor := processing.NewEntryProcessor(
		&mockMessageProducer,
		&mockMessageConsumer,
		smcEventIndexName,
		consumtionIndexName,
	)
	processor.HandleEntries()

	doneMessage := <-done
	fmt.Println(doneMessage)

	actualProcessedDataBytes := mockMessageProducer.Data.ToJSON()
	updateResourcesIfEnabled(expectedDataFileName, actualProcessedDataBytes)

	// Read expected outcome from resource file.
	expectedBytes, err := ioutil.ReadFile(expectedDataFileName)
	utils.FailOnError(err, "Could not read file "+expectedDataFileName)

	/*
		expectedData := testmodels.TestProcessedData{}
		expectedData.FromJSON(expectedBytes)
	*/

	// Assert
	if string(actualProcessedDataBytes) != string(expectedBytes) {
		t.Fatal("Expected json does not match actual json value of processed data.")
	}
}

func updateResourcesIfEnabled(resourceFileName string, newData []byte) {
	if updateResourcesEnabled {
		err := ioutil.WriteFile(resourceFileName, newData, 0600)
		utils.FailOnError(err, "Could not update resource file "+resourceFileName)
	}
}
