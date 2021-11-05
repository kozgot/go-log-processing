package service

import (
	"bufio"
	"io"
	"log"
	"sync"

	"github.com/kozgot/go-log-processing/parser/internal/rabbitmq"
)

func ParseLogFile(
	readCloser io.ReadCloser,
	logFileName string,
	wg *sync.WaitGroup,
	producer *rabbitmq.AmqpProducer) {
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
			producer.SendDataToPostProcessor(*finalParsedLine)
		}
	}

	readCloser.Close()
	log.Printf("  Done parsing log file: %s", logFileName)
}
