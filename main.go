package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	separator = flag.String("s", ";", "separator for CSV file")

	skipHeader = flag.Bool("h", false, "skip the first line of the CSV file")

	columns = flag.String("c", "", "comma-separated list of columns to print (e.g. 1,2,3)")
)

func main() {
	flag.Parse()
	reader := csv.NewReader(os.Stdin)
	reader.Comma = ';'

	if separator != nil && rune((*separator)[0]) != ';' {
		reader.Comma = rune((*separator)[0])
	}

	if skipHeader != nil && *skipHeader {
		if _, err := reader.Read(); err != nil {
			log.Fatal(err)
		}
	}

	all := true
	colIndices := make([]int, 0)
	if columns != nil && *columns != "" {
		all = false
		for _, n := range strings.Split(*columns, ",") {
			i, err := strconv.Atoi(n)
			if err != nil {
				log.Fatalf("invalid column index: %s", n)
			}

			colIndices = append(colIndices, i-1)
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		line := make([]string, 0)
		if all {
			line = record
			fmt.Fprintf(os.Stdout, "%s\n", strings.Join(line, ";"))
			continue
		}

		for _, i := range colIndices {
			if i < 0 || i >= len(record) {
				log.Fatalf("column index %d out of range %d", i, len(record))
			}

			line = append(line, record[i])
		}

		fmt.Fprintf(os.Stdout, "%s\n", strings.Join(line, ";"))
	}

}
