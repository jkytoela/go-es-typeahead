package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

const SEARCH_TERM int = 1
const FILE_NAME string = "Amazon Search Terms_Search Terms_US.csv"

func readData() [][]string {
	f, err := os.Open(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	r := csv.NewReader(f)

	// Skip headers
	if _, err := r.Read(); err != nil {
		log.Fatal(err)
	}

	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return rows
}

func main() {
	rows := readData()
	for i, row := range rows {
		term := row[SEARCH_TERM]
		fmt.Printf("%s - %d\n", term, i)
	}
}
