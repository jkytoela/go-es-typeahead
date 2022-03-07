package helpers

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func GetESClient() *elasticsearch.Client {
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

func BuildSearchQuery(searchString string) map[string]interface{} {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": searchString,
				"type":  "bool_prefix",
				"fields": []string{
					"word",
					"word._2gram",
					"word._3gram",
				},
			},
		},
	}
	return query
}
