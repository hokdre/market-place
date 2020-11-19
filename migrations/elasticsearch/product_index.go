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

func CreateIndexProduct(esClient *elasticsearch.Client) error {
	log.SetOutput(os.Stdout)
	log.Println("Migration Product : starting!")

	mapping := `
	{
		"mappings" : {
			"properties" : {
				"name" : {
					"type" : "search_as_you_type",
					"doc_values" : false,
					"max_shingle_size" : 3
				},
				"weight" : {
					"type" : "float"
				},
				"width" : {
					"type" : "float"
				},
				"height" : {
					"type" : "float"
				},
				"long" : {
					"type" : "float"
				},
				"description" : {
					"type" : "text"
				},
				"etalase" : {
					"type" : "keyword"
				},
				"category" : {
					"properties" : {
						"second_sub" : {
							"type" : "keyword"
						},
						"third_sub" : {
							"type" : "keyword"
						},
						"top" : {
							"type" : "keyword"
						}
					}
				},
				"tags" : {
					"type" : "keyword"
				},
				"colors" : {
					"type" : "keyword"
				},
				"sizes" : {
					"type" : "keyword"
				},
				"photos" : {
					"type" : "keyword"
				},
				"price" : {
					"type" : "float"
				},
				"stock" : {
					"type" : "float"
				},
				"merchant" : {
					"properties" : {
						"_id" : {
							"type" : "keyword"
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
						"location_point" : {
							"type" : "geo_point"
						},
						"name" : {
							"type" : "search_as_you_type",
							"doc_values" : false,
							"max_shingle_size" : 3
						},
						"phone" : {
							"type" : "keyword"
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
						"rating" : {"type" : "float"},
						"num_review" : {"type" : "integer"}
					}
				},
				"rating" : {
					"type" : "float"
				},
				"num_review" : {
					"type" : "integer"
				},
				"created_at" : {
					"type" : "date"
				},
				"updated_at" : {
					"type" : "date"
				},
				"suggest" : {
					"type" : "completion",
					"analyzer" : "simple",
					"preserve_separators" : true,
					"preserve_position_increments" : true,
					"max_input_length" : 50
				}
				
			}
		}
	}
	`
	indexName := "ecommerce.products"
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
	log.Println("Migration Product : success!")
	return nil
}
