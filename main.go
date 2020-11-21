package main

import (
	"log"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	infrastructureConfig "github.com/market-place/config/infrastructure_config"
	repoConfig "github.com/market-place/config/repo_config"
	usecaseConfig "github.com/market-place/config/usecase_config"
	"github.com/market-place/infrastructure/cloudstorage"
	"github.com/market-place/infrastructure/database"
	"github.com/market-place/infrastructure/http_api"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed loading .env file : %s \n", err)
		return
	}

	//cloud storage
	log.SetOutput(os.Stdout)
	log.Println("Setup Dependencies : starting!")
	gStorage, err := cloudstorage.NewGCPStorage()
	if err != nil {
		log.Fatalf("Failed setting google storage : %s \n", err)
		return
	}

	//secondary db
	esClient, err := database.NewElasticSearchDatabase()
	if err != nil {
		log.Fatalf("Failed setting secondary db: %s \n", err)
		return
	}
	redisClient := database.NewRedisDatabase()

	//main db
	mongoDatabase, err := database.NewMongoDatabase()
	if err != nil {
		log.Fatalf("Failed setting main db: %s \n", err)
		return
	}

	//infrastructure config
	infrastructureConf := infrastructureConfig.NewInfrastructureConfig(
		mongoDatabase,
		esClient,
		redisClient,
		gStorage,
	)

	//repository
	repoConf := repoConfig.NewRepoConfig(infrastructureConf)

	//logic or usecase
	usecaseConfig := usecaseConfig.NewUsecaseConfig(repoConf)

	//http
	r := mux.NewRouter()
	httpConfig := http_api.NewHttpAPI(r, usecaseConfig, infrastructureConf)
	log.Println("Setup Dependencies : finish!")

	log.Println("Server: starting!")
	port := ":80"
	readTimeOut := 2 * time.Second
	writeTimeOut := 2 * time.Second
	log.Printf("Starting server on port : %s", port)
	if err := httpConfig.StartServer(port, readTimeOut, writeTimeOut); err != nil {
		log.Fatalf("Failed starting http server : %s", err)
	}
}
