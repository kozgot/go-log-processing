package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	parsermodels "github.com/kozgot/go-log-processing/parser/pkg/models"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/streadway/amqp"
)

const logEntriesExchangeName = "logentries"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	log.Println("PostProcessor service starting...")
	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("RABBIT_URL:", rabbitMqURL)

	conn, err := amqp.Dial(rabbitMqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		logEntriesExchangeName, // name
		"fanout",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,                 // queue name
		"",                     // routing key
		logEntriesExchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			if strings.Contains(string(d.Body), "START") {
				fmt.Println("Start of entries...")
				continue
			} else if strings.Contains(string(d.Body), "END") {
				fmt.Println("End of entries...")
				continue
			}

			entry := deserializeMessage(d.Body)
			process(entry)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func deserializeMessage(message []byte) parsermodels.ParsedLine {
	var data parsermodels.ParsedLine
	if err := json.Unmarshal(message, &data); err != nil {
		fmt.Println("failed to unmarshal:", err)
	}

	return data
}

func process(logEntry parsermodels.ParsedLine) models.ProcessedEntries {
	entriesBySmcUID := make(map[string][]models.SmcEntry)
	routingEntries := []models.RoutingEntry{}
	statusEntries := []models.StatusEntry{}

	result := models.ProcessedEntries{}

	switch logEntry.Level {
	case "INFO":
		smcEntry, routingEntry, statusEntry := processInfo(logEntry)
		if smcEntry != nil {
			uid := smcEntry.UID
			_, ok := entriesBySmcUID[uid]
			if !ok {
				entriesBySmcUID[uid] = []models.SmcEntry{}
			}

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
		}

		if routingEntry != nil {
			routingEntries = append(routingEntries, *routingEntry)
		}
		if statusEntry != nil {
			statusEntries = append(statusEntries, *statusEntry)
		}
	case "WARN":
		smcEntry := processWarn(logEntry)
		if smcEntry != nil {
			uid := smcEntry.UID
			_, ok := entriesBySmcUID[uid]
			if !ok {
				entriesBySmcUID[uid] = []models.SmcEntry{}
			}

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
		}
	case "WARNING":
		smcEntry := processWarning(logEntry)
		if smcEntry != nil {
			uid := smcEntry.UID
			_, ok := entriesBySmcUID[uid]
			if !ok {
				entriesBySmcUID[uid] = []models.SmcEntry{}
			}

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
		}
	case "ERROR":
		smcEntry := processError(logEntry)

		if smcEntry != nil {
			uid := smcEntry.UID
			_, ok := entriesBySmcUID[uid]
			if !ok {
				entriesBySmcUID[uid] = []models.SmcEntry{}
			}

			entriesBySmcUID[uid] = append(entriesBySmcUID[uid], *smcEntry)
		}
	default:
		fmt.Printf("Unknown log level %s", logEntry.Level)
	}

	result.RoutingEntries = routingEntries
	result.StatusEntries = statusEntries
	result.SmcEntries = entriesBySmcUID

	return result
}

func processInfo(logEntry parsermodels.ParsedLine) (*models.SmcEntry, *models.RoutingEntry, *models.StatusEntry) {
	// one of 'ROUTING', 'JOIN', 'STATUS', or 'DC'
	switch logEntry.InfoParams.MessageType {
	case "ROUTING":
		routingEntry := processRoutingMessage(logEntry)
		return nil, routingEntry, nil
	case "JOIN":
		joinEntry := processJoinMessage(logEntry)
		return joinEntry, nil, nil
	case "STATUS":
		statusEntry := processStatusMessage(logEntry)
		return nil, nil, statusEntry
	case "DC":
		dcMessage := processDCMessage(logEntry)
		return dcMessage, nil, nil
	default:
		break
	}

	return nil, nil, nil
}

func processWarn(logEntry parsermodels.ParsedLine) *models.SmcEntry {
	result := models.SmcEntry{}

	// todo
	return &result
}

func processWarning(logEntry parsermodels.ParsedLine) *models.SmcEntry {
	result := models.SmcEntry{}

	// todo
	return &result
}

func processError(logEntry parsermodels.ParsedLine) *models.SmcEntry {
	result := models.SmcEntry{}

	// todo
	return &result
}

func processDCMessage(logEntry parsermodels.ParsedLine) *models.SmcEntry {
	result := models.SmcEntry{}

	return &result
}

func processJoinMessage(logEntry parsermodels.ParsedLine) *models.SmcEntry {
	result := models.SmcEntry{}

	return &result
}

func processStatusMessage(logEntry parsermodels.ParsedLine) *models.StatusEntry {
	result := models.StatusEntry{}

	return &result
}

func processRoutingMessage(logEntry parsermodels.ParsedLine) *models.RoutingEntry {
	result := models.RoutingEntry{}

	return &result
}
