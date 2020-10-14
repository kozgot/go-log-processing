// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	filePath := "dc_main.log"
	fmt.Println("Reading from" + filePath)

	file, ferr := os.Open(filePath)
	if ferr != nil {
		panic(ferr)
	}

	r, _ := regexp.Compile("^(Mon|Tue|Wed|Thu|Fri|Sat|Sun) (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) (0?[1-9]|[12][0-9]|3[01]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9] [0-2][0-9][0-9][0-9]")
	ctr := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ":") {
			date := r.FindString(line)
			// items := strings.Split(line, ":")
			// firstPart := strings.Join(items[0:3], ":")
			fmt.Printf("Line %d: %s\n", ctr, date)
		}
		ctr++
	}
}
