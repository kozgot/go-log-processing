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

	// todo: add rabbitmq service with DI
	channel, conn := rabbitmq.OpenChannelAndConnection(rabbitMqURL)
	defer rabbitmq.CloseChannelAndConnection(channel, conn)

	azureFileNames := azure.GetFileNamesFromAzure(azureStorageAccountName, azureStorageAccessKey, azureStorageContainer)

	var wg sync.WaitGroup

	for _, fileName := range azureFileNames {
		fmt.Println(fileName)
		readCloser := azure.DownloadFileFromAzure(
			fileName,
			azureStorageAccountName,
			azureStorageAccessKey,
			azureStorageContainer)

		wg.Add(1)
		go parser.ParseLogFile(readCloser, fileName, &wg, channel)
	}

	wg.Wait()

	// Send a message indicating that this is the end of the processing
	rabbitmq.SendStringMessageToPostProcessor("END", channel)
	log.Printf("  Sent END to Postprocessing service ...")
}
