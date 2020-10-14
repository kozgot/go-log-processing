// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

const logLevelsRegex = "(INFO|WARN)"
const dateFormatRegex = "^(Mon|Tue|Wed|Thu|Fri|Sat|Sun) (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) (0?[1-9]|[12][0-9]|3[01]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9] [0-2][0-9][0-9][0-9]"
const dateLayoutString = "Mon Jan  2 15:04:05 2006"

func main() {
	filePath := "dc_main.log"
	file, ferr := os.Open(filePath)
	if ferr != nil {
		panic(ferr)
	}

	r, _ := regexp.Compile(dateFormatRegex)
	levelRegex, _ := regexp.Compile(logLevelsRegex)

	scanner := bufio.NewScanner(file)
	relevantLines := []string{""}
	for scanner.Scan() {
		line := scanner.Text()
		dateString := r.FindString(line)
		if dateString != "" {

			date, err := time.Parse(dateLayoutString, dateString)
			if err != nil {
				panic(err)
			}

			logLevel := levelRegex.FindString(line)
			if logLevel != "" {
				relevantLines = append(relevantLines, line)
				fmt.Printf("Date %s: %s\n", date, logLevel)
			}
		}
	}
	fmt.Println(len(relevantLines))
}
