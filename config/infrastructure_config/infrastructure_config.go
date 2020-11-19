package config

import (
	"cloud.google.com/go/storage"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type InfrastructureConfig interface {
	GetMongoDBDatabase() *mongo.Database
	GetElasticClient() *elasticsearch.Client
	GetRedisClient() *redis.Client
	GetGoogleStorageClient() *storage.Client
}

type infrastuctureConfig struct {
	mongoInstance   *mongo.Database
	elasticInstance *elasticsearch.Client
	redisInstance   *redis.Client
	gcpInstance     *storage.Client
}

func NewInfrastructureConfig(
	mongoInstance *mongo.Database,
	elasticInstance *elasticsearch.Client,
	redisInstance *redis.Client,
	gcpInstance *storage.Client,
) InfrastructureConfig {
	return &infrastuctureConfig{
		mongoInstance:   mongoInstance,
		elasticInstance: elasticInstance,
		redisInstance:   redisInstance,
		gcpInstance:     gcpInstance,
	}
}

func (i *infrastuctureConfig) GetMongoDBDatabase() *mongo.Database {
	return i.mongoInstance
}

func (i *infrastuctureConfig) GetElasticClient() *elasticsearch.Client {
	return i.elasticInstance
}

func (i *infrastuctureConfig) GetRedisClient() *redis.Client {
	return i.redisInstance
}

func (i *infrastuctureConfig) GetGoogleStorageClient() *storage.Client {
	return i.gcpInstance
}
