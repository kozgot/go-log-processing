package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kozgot/go-log-processing/parser/internal/files"
	"github.com/kozgot/go-log-processing/parser/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/parser/internal/service"
	"github.com/streadway/amqp"
)

const unzippedInputFileFolderName = "unzipped_input_files"

func main() {
	// expects the file path from a command line argument (only works for dc_main.log files for now)
	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("Communicationg with RabbitMQ at: ", rabbitMqURL)

	if len(os.Args) == 0 {
		log.Fatalf("ERROR: Missing log file path param!")
	}

	filePath := os.Args[1]
	inputFiles, err := files.Unzip(filePath, unzippedInputFileFolderName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unzipped:\n" + strings.Join(inputFiles, "\n"))

	channel, conn := rabbitmq.OpenChannelAndConnection(rabbitMqURL)
	defer rabbitmq.CloseChannelAndConnection(channel, conn)

	var wg sync.WaitGroup

	for _, file := range inputFiles {
		wg.Add(1)
		go processFile(file, &wg, channel)
	}

	wg.Wait()

	// Send a message indicating that this is the end of the processing
	rabbitmq.SendStringMessageToElastic("DONE|", channel)
}

func processFile(filePath string, wg *sync.WaitGroup, channel *amqp.Channel) {
	defer wg.Done()

	file, ferr := os.Open(filePath)
	if ferr != nil {
		panic(ferr)
	}

	_, shortFileName := filepath.Split(filePath)

	log.Printf("  Processing log file: %s ...", shortFileName)
	scanner := bufio.NewScanner(file)
	indexName := shortFileName

	// Send the name of the index
	rabbitmq.SendStringMessageToElastic("CREATEINDEX|"+indexName, channel)
	log.Printf("  Creating index: %s ...", indexName)
	for scanner.Scan() {
		line := scanner.Text()
		relevantLine, success := service.Filter(line)
		if !success {
			continue
		}

		parsedLine, ok := service.ParseDate(*relevantLine)
		if !ok {
			continue
		}

		finalParsedLine := service.ParseContents(*parsedLine)
		if finalParsedLine != nil {
			rabbitmq.SendLineToElastic(*finalParsedLine, channel, indexName)
		}
	}

	log.Printf("  Done processing log file: %s", shortFileName)
}
