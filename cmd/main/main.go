// main.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	filter "github.com/kozgot/go-log-processing/cmd/filterlines"
	contentparser "github.com/kozgot/go-log-processing/cmd/parsecontents"
	"github.com/kozgot/go-log-processing/cmd/parsedates"
	"github.com/streadway/amqp"
)

func main() {
	// expects the file path from a command line argument (only works for dc_main.log files for now)
	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("Communicationg with RabbitMQ at: ", rabbitMqURL)

	if len(os.Args) == 0 {
		log.Fatalf("ERROR: Missing log file path param!")
	}

	filePath := os.Args[1]
	file, ferr := os.Open(filePath)
	if ferr != nil {
		panic(ferr)
	}

	scanner := bufio.NewScanner(file)
	indexName := "dc_main"
	// Send the name of the index
	sendStringMessageToElastic(rabbitMqURL, "[INDEXNAME] "+indexName)
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
		sendLinesToElastic(rabbitMqURL, *finalParsedLine)
	}

	// Send a message indicating that this is the end of the current index
	sendStringMessageToElastic(rabbitMqURL, "[DONE]")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func serializeLine(line contentparser.ParsedLine) []byte {
	bytes, err := json.Marshal(line)
	if err != nil {
		fmt.Println("Can't serialize", line)
	}

	return bytes
}

func sendData(rabbitMqURL string, data []byte) {
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

	body := data

	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")

	// log.Printf(" [PARSER] Sent a line to RabbitMQ")
	failOnError(err, "Failed to publish a message")
}

func sendLinesToElastic(rabbitMqURL string, line contentparser.ParsedLine) {
	byteData := serializeLine(line)
	sendData(rabbitMqURL, byteData)
}

func sendStringMessageToElastic(rabbitMqURL string, indexName string) {
	bytes, err := json.Marshal(indexName)
	fmt.Println(string(bytes))
	if err != nil {
		fmt.Println("Can't serialize", indexName)
	}
	sendData(rabbitMqURL, bytes)
}
