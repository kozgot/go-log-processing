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
	"github.com/kozgot/go-log-processing/elasticuploader/pkg/models"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

const numWorkers = 4
const flushBytes = 1000000
const processedDataExchangeName = "processeddata_direct_durable"

func main() {
	log.Println("Elastic Uploader starting...")
	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("RABBIT_URL:", rabbitMqURL)

	conn, err := amqp.Dial(rabbitMqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		processedDataExchangeName, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"save_data_queue_durable", // name
		true,                      // durable
		false,                     // delete when unused
		true,                      // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// TODO: extract routing key to a single place, eg.: env variables
	err = ch.QueueBind(
		q.Name,                    // queue name
		"save-data",               // routing key
		processedDataExchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	linesByIndexNames := make(map[string][]models.Message)
	forever := make(chan bool)
	dataCountTreshold := 1000
	documentID := 1

	go func() {
		for d := range msgs {
			msgParts := strings.Split(string(d.Body), "|")
			msgPrefix := msgParts[0]
			switch msgPrefix {
			case "CREATEINDEX":
				indexName := strings.Split(string(d.Body), "|")[1]
				createEsIndex(indexName)
			case "DONE":
				for key, value := range linesByIndexNames {
					documentID = BulkIndexerUpload(value, documentID, key)
					log.Printf("  Successfully indexed all %d documents (index name: %s)", documentID-1, key)
				}

				for k := range linesByIndexNames {
					delete(linesByIndexNames, k)
				}

				documentID = 1
			default:
				// save messages until we hit the 1000 line treshold
				data := deserialize(d.Body)
				_, ok := linesByIndexNames[data.IndexName]
				if !ok {
					linesByIndexNames[data.IndexName] = []models.Message{}
				}

				linesByIndexNames[data.IndexName] = append(linesByIndexNames[data.IndexName], models.Message{Content: data.Data})
				if len(linesByIndexNames[data.IndexName]) >= dataCountTreshold {
					documentID = BulkIndexerUpload(linesByIndexNames[data.IndexName], documentID, data.IndexName)
					linesByIndexNames[data.IndexName] = []models.Message{} // clear the buffer after uploading the contents
				}
			}

			// Acknowledge message
			err := d.Ack(false)
			failOnError(err, "Could not acknowledge message")
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// BulkIndexerUpload uploads data to Elasticsearch using BulkIndexer from go-elasticsearch.
func BulkIndexerUpload(lines []models.Message, currentDocumentID int, indexName string) int {
	var (
		countSuccessful uint64

		err error
	)

	fmt.Println(len((lines)))
	fmt.Println(currentDocumentID)
	fmt.Println(indexName)

	// Use a third-party package for implementing the backoff function
	retryBackoff := backoff.NewExponentialBackOff()

	documentID := currentDocumentID
	maxRetries := 10
	flushInterval := 30 * time.Second

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
		MaxRetries: maxRetries,
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Create the BulkIndexer
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,     // The default index name
		Client:        es,            // The Elasticsearch client
		NumWorkers:    numWorkers,    // The number of worker goroutines
		FlushBytes:    flushBytes,    // The flush threshold in bytes
		FlushInterval: flushInterval, // The periodic flush interval
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	start := time.Now().UTC()

	for _, a := range lines {
		// Add an item to the BulkIndexer
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",

				// DocumentID is the (optional) document ID
				DocumentID: strconv.Itoa(documentID),

				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(a.Content),

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

		documentID++

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

	return documentID
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

	log.Println("Deleting index:  ", index, "...")

	// Re-create the index
	if res, err = es.Indices.Delete(
		[]string{index},
		es.Indices.Delete.WithIgnoreUnavailable(true)); err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)
	}
	res.Body.Close()

	log.Println("Creating index:  ", index, "...")
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
	if biStats.NumFailed > 0 {
		// We got some errors while trying to index the documents
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		// Indexed evereything successfully
		log.Printf(
			"Successfully indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}

	log.Println(strings.Repeat("▔", 65))
}

func deserialize(dataBytes []byte) models.ReceivedDataUnit {
	var data models.ReceivedDataUnit
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		fmt.Println("failed to unmarshal:", err)
	}

	return data
}
