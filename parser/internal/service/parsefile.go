package service

import (
	"bufio"
	"io"
	"log"
	"sync"

	"github.com/kozgot/go-log-processing/parser/pkg/rabbitmq"
)

func ParseSingleFile(readCloser io.ReadCloser, logFileName string,
	wg *sync.WaitGroup,
	rabbitMQProducer rabbitmq.MessageProducer) {
	defer wg.Done()
	log.Printf("  [PARSER] Parsing log file: %s ...", logFileName)
	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
		line := scanner.Text()
		relevantLine, success := ParseLogLevelAndFilter(line)
		if !success {
			continue
		}

		lineWithTimestamp, ok := ParseTimestamp(*relevantLine)
		if !ok {
			continue
		}

		finalParsedLine := ParseContents(*lineWithTimestamp)
		if finalParsedLine != nil {
			rabbitMQProducer.PublishEntry(*finalParsedLine)
		}
	}

	readCloser.Close()
	log.Printf("  [PARSER] Done parsing log file: %s", logFileName)
}
