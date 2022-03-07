package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"typeahead/helpers"

	"github.com/gin-gonic/gin"
)

func getSearchResults(searchString string) []string {
	es := helpers.GetESClient()
	query := helpers.BuildSearchQuery(searchString)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("words"),
		es.Search.WithBody(&buf),
		es.Search.WithPretty(),
	)

	defer res.Body.Close()

	if err != nil {
		log.Fatalf("Error while searching: %s", err)
	}

	var result []string
	r := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error while decoding the response body: %s", err)
	}

	for _, res := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		result = append(result, res.(map[string]interface{})["_source"].(map[string]interface{})["word"].(string))
	}
	return result
}

func main() {
	r := gin.Default()
	r.GET("/search", func(c *gin.Context) {
		queryStr := c.Query("query")
		if len(queryStr) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Query string not provided",
			})
			return
		}

		searchResults := getSearchResults(queryStr)

		c.JSON(200, gin.H{
			"data": searchResults,
		})
	})
	r.Run()
}
