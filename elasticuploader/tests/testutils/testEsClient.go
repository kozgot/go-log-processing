package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/kozgot/go-log-processing/elasticuploader/internal/utils"
)

// TestEsClientWrapper is used in tests to query the uploaded data from ES, and do some cleanup after the test run.
type TestEsClientWrapper struct {
	esClient *elasticsearch.Client
}

// NewTestEsClientWrapper creates a new EsClientWrapper.
func NewTestEsClientWrapper() *TestEsClientWrapper {
	// Create the ES client
	// Use a third-party package for implementing the backoff function
	retryBackoff := backoff.NewExponentialBackOff()
	elasticSearchClient, err := elasticsearch.NewClient(elasticsearch.Config{
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
	utils.FailOnError(err, " [ESClient] Error creating the client")
	clientWrapper := TestEsClientWrapper{esClient: elasticSearchClient}

	return &clientWrapper
}

func (testEsClient *TestEsClientWrapper) QueryDocCountInIndex(indexName string) int {
	var (
		r map[string]interface{}
	)

	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	/*
			"query": {
		    "match_phrase": {
		      "EventType": "12"
		    }
		  }
	*/
	err := json.NewEncoder(&buf).Encode(query)
	utils.FailOnError(err, "Error encoding query")

	// Perform the search request.
	res, err := testEsClient.esClient.Search(
		testEsClient.esClient.Search.WithContext(context.Background()),
		testEsClient.esClient.Search.WithIndex(indexName),
		testEsClient.esClient.Search.WithBody(&buf),
		testEsClient.esClient.Search.WithTrackTotalHits(true),
	)
	utils.FailOnError(err, "Error getting response")

	if res.IsError() {
		var e map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&e)
		utils.FailOnError(err, "Error parsing the response body")

		// Print the response status and error information.
		log.Fatalf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	utils.FailOnError(err, "Error parsing the response body")

	hits := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		hits,
		int(r["took"].(float64)),
	)
	/*


		// Print the ID and document source for each hit.
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		}
	*/

	log.Println(strings.Repeat("=", 37))
	res.Body.Close()
	return hits
}

func (testEsClient *TestEsClientWrapper) DeleteIndex(indexName string) {
	var (
		res *esapi.Response
		err error
	)

	log.Println(" [TEST ES] Deleting index: ", indexName)

	// Delete the index
	if res, err = testEsClient.esClient.Indices.Delete(
		[]string{indexName},
		testEsClient.esClient.Indices.Delete.WithIgnoreUnavailable(true)); err != nil || res.IsError() {
		log.Fatalf(" [TEST ES] Failed to delete index: %s", err)
	}

	res.Body.Close()
}
