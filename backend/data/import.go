package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

const SEARCH_TERM int = 1
const FILE_NAME string = "Amazon Search Terms_Search Terms_US.csv"

func getWords() []string {
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

	// Read rest as rows
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var words []string
	for _, row := range rows {
		words = append(words, row[SEARCH_TERM])
	}

	return words
}

func getESClient() *elasticsearch.Client {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println(res)

	defer res.Body.Close()
	return es
}

const mapping = `
{
	"mappings": {
    "properties": {
      "word": {
        "type": "search_as_you_type"
      }
    }
  }
}`

func main() {
	// Create ES client
	es := getESClient()

	ctx := context.Background()

	// Read CSV file to rows
	words := getWords()

	// Configure request
	createIndexRequest := esapi.IndicesCreateRequest{
		Index: "words",
		Body:  bytes.NewReader([]byte(mapping)),
	}

	// Execute the request
	resCreateIndex, err := createIndexRequest.Do(ctx, es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	// Handle closing after this function exits
	defer resCreateIndex.Body.Close()

	// Create a new bulk indexer
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      "words", // The index name
		Client:     es,      // The ES client
		NumWorkers: 50,      // The number of workers
	})

	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	// For calculating the process duration
	start := time.Now().UTC()

	wordsLen := len(words)

	for i, word := range words {
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.Itoa(i),
				Body:       bytes.NewReader([]byte(`{"word" : "` + word + `"}`)),
			},
		)

		fmt.Printf("Indexing word %d/%d: %s\n", i+1, wordsLen, word)

		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	// Indexer statistics
	biStats := bi.Stats()

	log.Println(strings.Repeat("=", 75))

	duration := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed %d words with %d errors in %v (%d docs/sec)",
			int64(biStats.NumFlushed),
			int64(biStats.NumFailed),
			duration.Truncate(time.Second),
			int64(1000.0/float64(duration/time.Millisecond)*float64(biStats.NumFlushed)),
		)
	} else {
		log.Printf(
			"Indexed %d words in %v (%d docs/sec)",
			int64(biStats.NumFlushed),
			duration.Truncate(time.Millisecond),
			int64(1000.0/float64(duration/time.Millisecond)*float64(biStats.NumFlushed)),
		)
	}
}
