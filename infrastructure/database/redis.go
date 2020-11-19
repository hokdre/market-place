package database

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func NewRedisDatabase() *redis.Client {
	log.SetOutput(os.Stdout)
	log.Println("Settup Redis : starting!")

	url := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr: url,
	})

	log.Println("Setup Redis: succes!")
	return client
}
