package migrations

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

func CreateIndexMerchant(esClient *elasticsearch.Client) error {
	log.SetOutput(os.Stdout)
	log.Println("Migration Merchant : starting!")

	mapping := `
	{
		"mappings" : {
		   "properties" : {
			 "name" : {
			   "type" : "search_as_you_type",
			   "doc_values" : false,
			   "max_shingle_size" : 3
			 },
			 "address" : {
			   "properties" : {
				 "_id" : {
				   "type" : "keyword"
				 },
				 "city" : {
				   "properties" : {
					 "city_id" : {
					   "type" : "keyword"
					 },
					 "city_name" : {
					   "type" : "keyword"
					 },
					 "postal_code" : {
					   "type" : "keyword"
					 },
					 "province" : {
					   "type" : "keyword"
					 },
					 "province_id" : {
					   "type" : "keyword"
					 }
				   }
				 },
				 "number" : {
				   "type" : "keyword"
				 },
				 "street" : {
				   "type" : "keyword"
				 }
			   }
			 },
			 "avatar" : {
			   "type" : "keyword"
			 },
			 "phone" : {
			   "type" : "keyword"
			 },
			  "description" : {
			   "type" : "text"
			 },
			  "etalase" : {
			   "type" : "keyword"
			 },
			 "rating" : {
			   "type" : "float"
			 },
			 "num_review" : {
			   "type": "float"
			 },
			 "shippings" : {
			   "properties" : {
				 "_id" : {
				   "type" : "keyword"
				 },
				 "created_at" : {
				   "type" : "date"
				 },
				 "name" : {
				   "type" : "keyword"
				 },
				 "updated_at" : {
				   "type" : "date"
				 }
			   }
			 },
			 "bank_accounts" : {
			   "properties" : {
				 "_id" : {
				   "type" : "keyword"
				 },
				 "bank_code" : {
				   "type" : "keyword"
				 },
				 "number" : {
				   "type" : "keyword"
				 }
			   }
			 },
			 "location_point" : {
			   "type" : "geo_point"
			 },
			 "created_at" : {
			   "type" : "date"
			 },
			 "updated_at" : {
			   "type" : "date"
			 }
		   }
		}
	}
	`
	indexName := "ecommerce.merchants"
	ctxCreateIndex, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body := strings.NewReader(mapping)
	res, err := esClient.Indices.Create(
		indexName,
		esClient.Indices.Create.WithBody(body),
		esClient.Indices.Create.WithContext(ctxCreateIndex),
	)

	if err != nil {
		log.Printf("Migration product : failed cause, %s \n", err)
		return err
	}

	if res.IsError() {
		log.Printf("Migration product : failed cause, %s \n", res.Status())
		return errors.New("mapping product failed!")
	}
	log.Println("Migration Merchant : success!")
	return nil
}
