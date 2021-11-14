package logparser

import (
	"log"
	"sync"

	"github.com/kozgot/go-log-processing/parser/internal/service"
	"github.com/kozgot/go-log-processing/parser/pkg/filedownloader"
	"github.com/kozgot/go-log-processing/parser/pkg/rabbitmq"
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

	var wg sync.WaitGroup
	for _, fileName := range azureFileNames {
		readCloser := logparser.fileDownloader.DownloadFile(fileName)

		wg.Add(1)
		go service.ParseSingleFile(readCloser, fileName, &wg, logparser.rabbitMqProducer)
	}
	wg.Wait()

	forever := make(chan bool)
	// Send a message indicating that we have reached the end of the log files.
	logparser.rabbitMqProducer.PublishStringMessage("END")
	log.Printf("  [PARSER] Sent END to Postprocessing service ...")

	log.Printf("  [PARSER] Done parsing files, press CTRL+C to exit ...")
	<-forever
}
