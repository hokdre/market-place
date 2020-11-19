package cloudstorage

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func NewGCPStorage() (*storage.Client, error) {
	log.SetOutput(os.Stdout)
	log.Println("Settup GCP : starting!")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile("gcp.json"))
	if err != nil {
		log.Printf("Setup GCP: failed cause, %s", err)
		return nil, err
	}

	log.Println("Settup GCP : success!")
	return client, nil
}
