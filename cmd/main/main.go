// main.go
package main

import (
	"bufio"
	"fmt"
	"os"

	filter "github.com/kozgot/go-log-processing/cmd/filterlines"
	"github.com/kozgot/go-log-processing/cmd/parsedates"
)

func main() {
	// expects the file path from a command line argument (only works for dc_main.log files for now)
	filePath := os.Args[1]
	file, ferr := os.Open(filePath)
	if ferr != nil {
		panic(ferr)
	}

	scanner := bufio.NewScanner(file)

	relevantLines := []parsedates.LineWithDate{}
	for scanner.Scan() {
		line := scanner.Text()
		relevantLine, success := filter.Filter(line)
		if !success {
			continue
		}

		parsedLine, ok := parsedates.ParseDate(*relevantLine)
		if !ok {
			continue
		}

		relevantLines = append(relevantLines, *parsedLine)
	}

	fmt.Println(len(relevantLines))
}
