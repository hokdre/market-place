package database

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

func NewElasticSearchDatabase() (*elasticsearch.Client, error) {
	log.SetOutput(os.Stdout)
	log.Println("Settup ElasticSearch : starting!")
	url := os.Getenv("ELASTIC_URL")

	cfg := elasticsearch.Config{
		MaxRetries: 10,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 10 * time.Second,
		},
		Addresses: []string{
			fmt.Sprintf("http://%s", url),
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Setup ElasticSearch: failed cause, %s", err)
		return nil, err
	}

	_, err = es.Info()
	if err != nil {
		log.Printf("Setup ElasticSearch: failed cause, %s", err)
		return nil, err
	}

	log.Println("Settup ElasticSearch : success!")
	return es, nil
}
