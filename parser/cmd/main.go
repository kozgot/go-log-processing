package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kozgot/go-log-processing/parser/internal/service"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
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
	files, err := Unzip(filePath, unzippedInputFileFolderName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))

	for _, file := range files {
		processFile(file, rabbitMqURL)
	}
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return make([]string, 0), err
	}
	defer r.Close()

	filenames := make([]string, 0, len(r.File))

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			error := os.MkdirAll(fpath, os.ModePerm)
			if error != nil {
				panic(error)
			}

			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func processFile(filePath string, rabbitMqURL string) {
	file, ferr := os.Open(filePath)
	if ferr != nil {
		panic(ferr)
	}

	_, shortFileName := filepath.Split(filePath)

	log.Printf("  Processing log file: %s ...", shortFileName)
	scanner := bufio.NewScanner(file)

	// Send the name of the index
	sendStringMessageToElastic(rabbitMqURL, "[INDEXNAME] "+shortFileName)
	log.Printf("  Creating index: %s ...", shortFileName)
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
			sendLinesToElastic(rabbitMqURL, *finalParsedLine)
		}
	}

	// Send a message indicating that this is the end of the current index
	sendStringMessageToElastic(rabbitMqURL, "[DONE]")
	log.Printf("  Done processing log file: %s", shortFileName)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func serializeLine(line models.ParsedLine) []byte {
	bytes, err := json.Marshal(line)
	if err != nil {
		fmt.Println("Can't serialize", line)
	}

	return bytes
}

func sendData(rabbitMqURL string, data []byte) {
	conn, err := amqp.Dial(rabbitMqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// create the channel
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

	body := data

	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")

	failOnError(err, "Failed to publish a message")
}

func sendLinesToElastic(rabbitMqURL string, line models.ParsedLine) {
	byteData := serializeLine(line)
	sendData(rabbitMqURL, byteData)
}

func sendStringMessageToElastic(rabbitMqURL string, indexName string) {
	bytes, err := json.Marshal(indexName)
	if err != nil {
		fmt.Println("Can't serialize", indexName)
	}
	sendData(rabbitMqURL, bytes)
}
