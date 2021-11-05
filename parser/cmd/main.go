package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kozgot/go-log-processing/parser/internal/azure"
	"github.com/kozgot/go-log-processing/parser/internal/rabbitmq"
	parser "github.com/kozgot/go-log-processing/parser/internal/service"
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

	// todo: add rabbitmq service with DI
	channel, conn := rabbitmq.OpenChannelAndConnection(rabbitMqURL, logEntriesExchangeName)
	defer rabbitmq.CloseChannelAndConnection(channel, conn)

	azureFileDownloader := azure.SetupDownloader(azureStorageAccountName, azureStorageAccessKey, azureStorageContainer)

	azureFileNames := azureFileDownloader.GetFileNamesFromAzure()

	var wg sync.WaitGroup

	for _, fileName := range azureFileNames {
		fmt.Println(fileName)
		readCloser := azureFileDownloader.DownloadFileFromAzure(fileName)

		wg.Add(1)
		go parser.ParseLogFile(readCloser, fileName, &wg, channel, logEntriesExchangeName, processEntryRoutingKey)
	}

	wg.Wait()

	// Send a message indicating that this is the end of the processing
	rabbitmq.SendStringMessageToPostProcessor("END", channel, logEntriesExchangeName, processEntryRoutingKey)
	log.Printf("  Sent END to Postprocessing service ...")
}
