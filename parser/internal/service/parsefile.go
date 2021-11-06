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
	log.Printf("  Parsing log file: %s ...", logFileName)
	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
		line := scanner.Text()
		relevantLine, success := Filter(line)
		if !success {
			continue
		}

		parsedLine, ok := ParseDate(*relevantLine)
		if !ok {
			continue
		}

		finalParsedLine := ParseContents(*parsedLine)
		if finalParsedLine != nil {
			rabbitMQProducer.PublishEntry(*finalParsedLine)
		}
	}

	readCloser.Close()
	log.Printf("  Done parsing log file: %s", logFileName)
}
