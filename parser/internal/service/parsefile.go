package service

import (
	"bufio"
	"io"
	"log"
	"sync"

	"github.com/kozgot/go-log-processing/parser/internal/rabbitmq"
	"github.com/streadway/amqp"
)

func ParseLogFile(
	readCloser io.ReadCloser,
	logFileName string,
	wg *sync.WaitGroup,
	channel *amqp.Channel,
	exchangeName string,
	routingKey string) {
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
			rabbitmq.SendLineToPostProcessor(*finalParsedLine, channel, routingKey, exchangeName)
		}
	}

	readCloser.Close()
	log.Printf("  Done parsing log file: %s", logFileName)
}
