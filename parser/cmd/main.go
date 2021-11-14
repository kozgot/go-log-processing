package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kozgot/go-log-processing/parser/pkg/filedownloader"
	"github.com/kozgot/go-log-processing/parser/pkg/logparser"
	"github.com/kozgot/go-log-processing/parser/pkg/rabbitmq"
)

func main() {
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

	log.Printf("  [PARSER] Press CTRL+C to exit...")

	forever := make(chan bool)
	<-forever
}
