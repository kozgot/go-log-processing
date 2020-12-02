// main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	elasticuploader "github.com/kozgot/go-log-processing/cmd/elasticsearch"
	filter "github.com/kozgot/go-log-processing/cmd/filterlines"
	contentparser "github.com/kozgot/go-log-processing/cmd/parsecontents"
	"github.com/kozgot/go-log-processing/cmd/parsedates"
	"github.com/streadway/amqp"
)

func main() {
	// expects the file path from a command line argument (only works for dc_main.log files for now)
	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("Communicationg with RabbitMQ at: ", rabbitMqURL)

	count := 2
	if len(os.Args) >= count {
		// temporary rabbit	MQ tutorial code
		seq := 1
		for {
			sendHello(rabbitMqURL, seq)
			time.Sleep(1 * time.Second)
			seq++
		}

		// log.Fatalf("ERROR: Missing log file path param!")
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

func sendHello(rabbitMqURL string, seq int) {
	conn, err := amqp.Dial(rabbitMqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// create the channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body := "Hello World!  " + strconv.Itoa(seq)

	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
