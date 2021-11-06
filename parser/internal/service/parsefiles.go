package service

import (
	"fmt"
	"log"
	"sync"

	"github.com/kozgot/go-log-processing/parser/internal/filedownloader"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

// RabbitMQProducer encapsulates methods used to communicate with rabbitMQ server.
type RabbitMQProducer interface {
	PublishStringMessage(indexName string)
	PublishEntry(line models.ParsedLogEntry)
	OpenChannelAndConnection()
	CloseChannelAndConnection()
}

// LogParser encapsulates parser data and logic.
type LogParser struct {
	fileDownloader   filedownloader.FileDownloader
	rabbitMqProducer RabbitMQProducer
}

// NewLogParser creates a new LogParser.
func NewLogParser(fileDownloader filedownloader.FileDownloader, rabbitMqProducer RabbitMQProducer) *LogParser {
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
		fmt.Println(fileName)
		readCloser := logparser.fileDownloader.DownloadFile(fileName)

		wg.Add(1)
		go ParseSingleFile(readCloser, fileName, &wg, logparser.rabbitMqProducer)
	}
	wg.Wait()

	// Send a message indicating that this is the end of the processing
	logparser.rabbitMqProducer.PublishStringMessage("END")
	log.Printf("  Sent END to Postprocessing service ...")
}
