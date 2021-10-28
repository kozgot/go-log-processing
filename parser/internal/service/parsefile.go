package service

import (
	"bufio"
	"log"
	"sync"

	"github.com/kozgot/go-log-processing/parser/internal/rabbitmq"
	"github.com/streadway/amqp"
)

func ParseLogFile(scanner *bufio.Scanner, shortFileName string, wg *sync.WaitGroup, channel *amqp.Channel) {
	defer wg.Done()
	log.Printf("  Parsing log file: %s ...", shortFileName)

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
			rabbitmq.SendLineToPostProcessor(*finalParsedLine, channel)
		}
	}

	log.Printf("  Done parsing log file: %s", shortFileName)
}
