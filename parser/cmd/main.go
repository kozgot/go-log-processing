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
	parser "github.com/kozgot/go-log-processing/parser/internal/service"
)

const unzippedInputFileFolderName = "unzipped_input_files"

func main() {
	rabbitMqURL := os.Getenv("RABBIT_URL")
	fmt.Println("Communicationg with RabbitMQ at: ", rabbitMqURL)

	if len(os.Args) == 0 {
		log.Fatalf("ERROR: Missing log file path param!")
	}

	zipFilePath := os.Args[1]
	inputFiles, err := files.Unzip(zipFilePath, unzippedInputFileFolderName)
	if err != nil {
		log.Fatal(err)
	}

	// todo: get the files from azure here
	// azure.Cucc()

	fmt.Println("Unzipped:\n" + strings.Join(inputFiles, "\n"))

	channel, conn := rabbitmq.OpenChannelAndConnection(rabbitMqURL)
	defer rabbitmq.CloseChannelAndConnection(channel, conn)

	var wg sync.WaitGroup

	for _, filePath := range inputFiles {
		file, ferr := os.Open(filePath)
		if ferr != nil {
			panic(ferr)
		}

		_, shortFileName := filepath.Split(filePath)

		scanner := bufio.NewScanner(file)

		wg.Add(1)
		go parser.ParseLogFile(scanner, shortFileName, &wg, channel)
	}

	wg.Wait()

	// Send a message indicating that this is the end of the processing
	rabbitmq.SendStringMessageToPostProcessor("END", channel)
	log.Printf("  Sent END to Postprocessing service ...")
}
