// main.go
package main

import (
	"bufio"
	"log"
	"os"
	"time"

	elasticuploader "github.com/kozgot/go-log-processing/cmd/elasticsearch"
	filter "github.com/kozgot/go-log-processing/cmd/filterlines"
	contentparser "github.com/kozgot/go-log-processing/cmd/parsecontents"
	"github.com/kozgot/go-log-processing/cmd/parsedates"
	"github.com/streadway/amqp"
)

func main() {
	// expects the file path from a command line argument (only works for dc_main.log files for now)
	if len(os.Args) >= 2 {
		/////////////////////////////////////////////////
		// temporary rabbit	MQ tutorial code
		for {
			sendHello()
			time.Sleep(1 * time.Second)
		}
		///////////////////////////////////////////////////////

		//log.Fatalf("ERROR: Missing log file path param!")
	}

	filePath := os.Args[1]
	file, ferr := os.Open(filePath)
	if ferr != nil {
		panic(ferr)
	}

	scanner := bufio.NewScanner(file)
	relevantLines := []contentparser.ParsedLine{}
	for scanner.Scan() {
		line := scanner.Text()
		relevantLine, success := filter.Filter(line)
		if !success {
			continue
		}

		parsedLine, ok := parsedates.ParseDate(*relevantLine)
		if !ok {
			continue
		}

		finalParsedLine := contentparser.ParseContents(*parsedLine)
		relevantLines = append(relevantLines, *finalParsedLine)
	}

	elasticuploader.BulkIndexerUpload(relevantLines)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendHello() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// create the channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// queue
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
