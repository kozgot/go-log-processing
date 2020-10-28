package elasticsearch

import (
	"log"

	contentparser "github.com/kozgot/go-log-processing/cmd/parsecontents"

	"github.com/elastic/go-elasticsearch/v7"
)

// UploadData uploads the parsed log data to elasticsearch
func UploadData(data []contentparser.ParsedLine) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, eerr := elasticsearch.NewClient(cfg)
	if eerr != nil {
		log.Fatalf("Error creating the client: %s", eerr)
	}
	log.Println(elasticsearch.Version)

	info, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	/*
		NOTE: It is critical to both close the response body and to consume it,
		in order to re-use persistent TCP connections in the default HTTP transport.
		If you're not interested in the response body, call io.Copy(ioutil.Discard, res.Body).
	*/
	defer info.Body.Close()
	log.Println(info)
	log.Println(len(data))
}
