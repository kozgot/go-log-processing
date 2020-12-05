package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// ErrorParams contains the parsed error parameters
type ErrorParams struct {
	ErrorCode   int
	Message     string
	Severity    int
	Description string
	Source      string
}

// ParsedLine contains a parsed line from the log file
type ParsedLine struct {
	Timestamp   time.Time
	Level       string
	ErrorParams *ErrorParams
}

var numWorkers = 4
var flushBytes = 1000000

// will be overwritten
var indexName = "index"

func main() {
	log.Println("Reciever started and listening...")
	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("RABBIT_URL:", rabbitMqURL)

	conn, err := amqp.Dial(rabbitMqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

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
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
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

	lines := []ParsedLine{}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			if strings.Contains(string(d.Body), "[INDEXNAME] ") {
				indexNameMessageString := string(d.Body)
				runes := []rune(indexNameMessageString)
				indexName = string(runes[len("[INDEXNAME] ")+1 : len(indexNameMessageString)-1])
				log.Println("Creating index:  ", indexName)
				createEsIndex(indexName)
			} else if strings.Contains(string(d.Body), "[DONE]") {
				BulkIndexerUpload(lines)
				lines = []ParsedLine{} // clear the buffer after uploading the contents
				log.Println("----------    DONE    -----------")
			} else {
				// save messages
				line := deserializeLines(d.Body)
				lines = append(lines, line)
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func deserializeLines(bytes []byte) ParsedLine {
	var line ParsedLine
	err := json.Unmarshal(bytes, &line)
	if err != nil {
		fmt.Println("ERROR ----- Can't deserislize message: " + string(bytes))
	}

	return line
}

// BulkIndexerUpload uploads data to Elasticsearch using BulkIndexer from go-elasticsearch
func BulkIndexerUpload(lines []ParsedLine) {
	var (
		countSuccessful uint64

		err error
	)

	numItems := len(lines)

	log.Printf(
		"\x1b[1mBulkIndexer\x1b[0m: documents [%s] workers [%d] flush [%s]",
		humanize.Comma(int64(numItems)), numWorkers, humanize.Bytes(uint64(flushBytes)))

	// Use a third-party package for implementing the backoff function
	retryBackoff := backoff.NewExponentialBackOff()

	// Create the Elasticsearch client
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		// Retry on 429 TooManyRequests statuses
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 10,
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Create the BulkIndexer
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,        // The default index name
		Client:        es,               // The Elasticsearch client
		NumWorkers:    numWorkers,       // The number of worker goroutines
		FlushBytes:    int(flushBytes),  // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	start := time.Now().UTC()

	for i, a := range lines {
		data, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("Cannot encode line %d: %s", i, err)
		}

		// Add an item to the BulkIndexer
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",

				// DocumentID is the (optional) document ID
				DocumentID: strconv.Itoa(i),

				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(data),

				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}

	// Close the indexer
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	dur := time.Since(start)
	biStats := bi.Stats()

	reportBulkIndexerStats(biStats, dur)
}

func createEsIndex(index string) {
	var (
		res *esapi.Response
		err error
	)

	// Use a third-party package for implementing the backoff function
	retryBackoff := backoff.NewExponentialBackOff()

	// Create the Elasticsearch client
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		// Retry on 429 TooManyRequests statuses
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 10,
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Re-create the index
	if res, err = es.Indices.Delete([]string{index}, es.Indices.Delete.WithIgnoreUnavailable(true)); err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)
	}
	res.Body.Close()
	res, err = es.Indices.Create(index)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}
	res.Body.Close()

}

func reportBulkIndexerStats(biStats esutil.BulkIndexerStats, dur time.Duration) {
	log.Println(strings.Repeat("▔", 65))

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}

	log.Println(strings.Repeat("▔", 65))
}
