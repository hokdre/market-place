package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func NewMongoDatabase() (*mongo.Database, error) {
	log.SetOutput(os.Stdout)
	log.Println("Settup Mongodb : starting!")

	if DB == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		url := os.Getenv("MONGO_URL")

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", url)))
		if err != nil {
			log.Printf("Setup Mongodb: failed cause, %s", err)
			return nil, err
		}

		DB = client.Database("ecommerce")
	}
	log.Println("Settup Mongodb : success!")
	return DB, nil
}
