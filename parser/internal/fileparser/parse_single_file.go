package fileparser

import (
	"bufio"
	"io"
	"log"
	"sync"

	"github.com/kozgot/go-log-processing/parser/internal/contentparser"
	"github.com/kozgot/go-log-processing/parser/internal/loglevelparser"
	"github.com/kozgot/go-log-processing/parser/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/parser/internal/timestampparser"
)

func ParseSingleFile(readCloser io.ReadCloser, logFileName string,
	wg *sync.WaitGroup,
	rabbitMQProducer rabbitmq.MessageProducer) {
	defer wg.Done()
	log.Printf("  [PARSER] Parsing log file: %s ...", logFileName)
	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse the log level, and filter out irrelevant lines eg.: VERBOSE log level.
		relevantLine := loglevelparser.ParseLogLevelAndFilter(line)
		if relevantLine == nil {
			continue
		}

		// Parse the timestamp of the log entry.
		lineWithTimestamp := timestampparser.ParseTimestamp(*relevantLine)
		if lineWithTimestamp == nil {
			continue
		}

		// Parse the remaining contents of the log entry depending on the log level.
		finalParsedLine := contentparser.ParseEntryContents(*lineWithTimestamp)
		if finalParsedLine == nil {
			continue
		}

		rabbitMQProducer.PublishEntry(*finalParsedLine)
	}

	readCloser.Close()
	log.Printf("  [PARSER] Done parsing log file: %s", logFileName)
}
