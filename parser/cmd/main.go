package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kozgot/go-log-processing/parser/internal/filedownloader"
	"github.com/kozgot/go-log-processing/parser/internal/logparser"
	"github.com/kozgot/go-log-processing/parser/internal/rabbitmq"
)

func main() {
	http.HandleFunc("/process/", handler)
	log.Printf("  [PARSER] Application started, listening on port 8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	rabbitMqURL := os.Getenv("RABBIT_URL")
	log.Println("RabbitMQ URL: ", rabbitMqURL)

	azureStorageAccountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	log.Println("Azure storage account name: ", azureStorageAccountName)

	azureStorageContainer := os.Getenv("AZURE_STORAGE_CONTAINER")
	log.Println("Azure storage container: ", azureStorageContainer)

	azureStorageAccessKey := os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	log.Println("Azure storage access key: ", azureStorageAccessKey[0:5]+"...")

	if len(azureStorageAccountName) == 0 || len(azureStorageAccessKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}

	logEntriesExchangeName := os.Getenv("LOG_ENTRIES_EXCHANGE")
	fmt.Println("LOG_ENTRIES_EXCHANGE:", logEntriesExchangeName)
	if len(logEntriesExchangeName) == 0 {
		log.Fatal("The LOG_ENTRIES_EXCHANGE environment variable is not set")
	}

	processEntryRoutingKey := os.Getenv("PROCESS_ENTRY_ROUTING_KEY")
	fmt.Println("PROCESS_ENTRY_ROUTING_KEY:", processEntryRoutingKey)
	if len(processEntryRoutingKey) == 0 {
		log.Fatal("The PROCESS_ENTRY_ROUTING_KEY environment variable is not set")
	}

	fmt.Fprintf(w, "<div>Storage account: %s, container: %s</div>",
		azureStorageAccountName,
		azureStorageContainer)

	// Initialize rabbitMQ producer.
	rabbitMqProducer := rabbitmq.NewAmqpProducer(processEntryRoutingKey, logEntriesExchangeName, rabbitMqURL)

	// Open a connection and a channel to send the log entries to.
	rabbitMqProducer.OpenChannelAndConnection()
	defer rabbitMqProducer.CloseChannelAndConnection()

	// Initialize file downloader.
	azureFileDownloader := filedownloader.NewAzureDownloader(
		azureStorageAccountName,
		azureStorageAccessKey,
		azureStorageContainer)

	// Init and run parser.
	logParser := logparser.NewLogParser(azureFileDownloader, rabbitMqProducer)
	logParser.ParseLogfiles()

	fmt.Fprint(w, "<div>Finished parsing log files, allow a few seconds for the processing to finish...</div>")
	fmt.Fprintf(w, "<a href=\"http://localhost:5601/app/home#/\">Check results in Kibana</a>")
}
