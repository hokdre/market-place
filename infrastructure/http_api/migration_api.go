package http_api

import (
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/market-place/infrastructure/http_api/http_response"
	migrations "github.com/market-place/migrations/elasticsearch"
)

type MigrationAPI interface {
	ElasticProductIndex(w http.ResponseWriter, r *http.Request)
	ElasticMerchantIndex(w http.ResponseWriter, r *http.Request)
}

type migrationAPI struct {
	esClient *elasticsearch.Client
}

func NewMigrationAPI(
	esClient *elasticsearch.Client,
) MigrationAPI {
	return &migrationAPI{
		esClient: esClient,
	}
}

func (m *migrationAPI) ElasticProductIndex(w http.ResponseWriter, r *http.Request) {
	if err := migrations.CreateIndexProduct(m.esClient); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	data := map[string]interface{}{
		"message": "product index is success created!",
	}
	http_response.SendOkJSON(w, http.StatusCreated, data)
}

func (m *migrationAPI) ElasticMerchantIndex(w http.ResponseWriter, r *http.Request) {
	if err := migrations.CreateIndexMerchant(m.esClient); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	data := map[string]interface{}{
		"message": "merchant index is success created!",
	}
	http_response.SendOkJSON(w, http.StatusCreated, data)
}
