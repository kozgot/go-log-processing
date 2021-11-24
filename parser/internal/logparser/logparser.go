package logparser

import (
	"log"
	"sync"

	"github.com/kozgot/go-log-processing/parser/internal/filedownloader"
	"github.com/kozgot/go-log-processing/parser/internal/fileparser"
	"github.com/kozgot/go-log-processing/parser/internal/rabbitmq"
)

// LogParser encapsulates parser data and logic.
type LogParser struct {
	fileDownloader   filedownloader.FileDownloader
	rabbitMqProducer rabbitmq.MessageProducer
}

// NewLogParser creates a new LogParser.
func NewLogParser(fileDownloader filedownloader.FileDownloader, rabbitMqProducer rabbitmq.MessageProducer) *LogParser {
	logparser := LogParser{
		fileDownloader:   fileDownloader,
		rabbitMqProducer: rabbitMqProducer,
	}

	return &logparser
}

// ParseLogfiles downloads log files from the given filedownloader, parses the log entries
// and forwards them to the provided rabbitMQ producer.
func (logparser *LogParser) ParseLogfiles() {
	azureFileNames := logparser.fileDownloader.ListFileNames()

	// Send a message indicating that we are starting to parse the log entries.
	logparser.rabbitMqProducer.PublishStringMessage("START")

	var wg sync.WaitGroup
	for _, fileName := range azureFileNames {
		readCloser := logparser.fileDownloader.DownloadFile(fileName)

		wg.Add(1)
		go fileparser.ParseSingleFile(readCloser, fileName, &wg, logparser.rabbitMqProducer)
	}
	wg.Wait()

	// Send a message indicating that we have reached the end of the log files.
	logparser.rabbitMqProducer.PublishStringMessage("END")
	log.Printf("  [PARSER] Sent END to Postprocessing service ...")

	log.Printf("  [PARSER] Finished parsing all files")
}
