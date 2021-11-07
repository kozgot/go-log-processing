package logparserunittests

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/kozgot/go-log-processing/parser/internal/utils"
	"github.com/kozgot/go-log-processing/parser/pkg/logparser"
	"github.com/kozgot/go-log-processing/parser/pkg/mocks"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/parser/tests/testmodels"
)

const updateResourcesEnabled = false

func TestLogParserDCMain(t *testing.T) {
	// Init a mock message producer.
	mockMessageProducer := mocks.MessageProducerMock{Entries: []models.ParsedLogEntry{}}

	// Create mock filedownloader.
	mockFileDownloader := mocks.MockFileDownloader{FileNameToDownload: "./resources/test_dc_main.log"}

	// Run parser
	logParser := logparser.NewLogParser(&mockFileDownloader, &mockMessageProducer)
	logParser.ParseLogfiles()

	// The number of relevant lines in the provided test log file.
	expectedEntryCount := 40
	actualEntryCount := len(mockMessageProducer.Entries)

	log.Printf("Expected %d entries.\n", expectedEntryCount)
	log.Printf("Got %d entries.\n", actualEntryCount)

	testParsedLogFile := testmodels.TestParsedLogFile{Lines: mockMessageProducer.Entries}
	actualBytes := testParsedLogFile.ToJSON()

	updateResourcesIfEnabled("./resources/expected_dc_main.json", actualBytes)

	// Read expected outcome from resource file.
	expectedBytes, err := ioutil.ReadFile("./resources/expected_dc_main.json")
	utils.FailOnError(err, "Could not read test json file.")

	// Assert
	if string(actualBytes) != string(expectedBytes) {
		t.Fatal("Expected json does not match actual json value of created partsed entries.")
	}
	if actualEntryCount != expectedEntryCount {
		t.Fatalf("Expected %d entries, got %d entries.", expectedEntryCount, actualEntryCount)
	}
}

func TestLogParserPLCManager(t *testing.T) {
	// Init a mock message producer.
	mockMessageProducer := mocks.MessageProducerMock{Entries: []models.ParsedLogEntry{}}

	// Create mock filedownloader.
	mockFileDownloader := mocks.MockFileDownloader{FileNameToDownload: "./resources/test_plc_manager.log"}

	// Run parser
	logParser := logparser.NewLogParser(&mockFileDownloader, &mockMessageProducer)
	logParser.ParseLogfiles()

	// The number of relevant lines in the provided test log file.
	expectedEntryCount := 50
	actualEntryCount := len(mockMessageProducer.Entries)

	log.Printf("Expected %d entries.\n", expectedEntryCount)
	log.Printf("Got %d entries.\n", actualEntryCount)

	testParsedLogFile := testmodels.TestParsedLogFile{Lines: mockMessageProducer.Entries}
	actualBytes := testParsedLogFile.ToJSON()

	updateResourcesIfEnabled("./resources/expected_plc_manager.json", actualBytes)

	// Read expected outcome from resource file.
	expectedBytes, err := ioutil.ReadFile("./resources/expected_plc_manager.json")
	utils.FailOnError(err, "Could not read test json file.")

	// Assert
	if string(actualBytes) != string(expectedBytes) {
		t.Fatal("Expected json does not match actual json value of created partsed entries.")
	}
	if actualEntryCount != expectedEntryCount {
		t.Fatalf("Expected %d entries, got %d entries.", expectedEntryCount, actualEntryCount)
	}
}

func updateResourcesIfEnabled(resourceFileName string, newData []byte) {
	if updateResourcesEnabled {
		_ = ioutil.WriteFile(resourceFileName, newData, 0600)
	}
}
