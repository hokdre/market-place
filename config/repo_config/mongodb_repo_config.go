package config

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/market-place/usecase/repository"
	elasticRepo "github.com/market-place/usecase/repository/elasticsearch"
	mongoRepo "github.com/market-place/usecase/repository/mongodb"
	redisRepo "github.com/market-place/usecase/repository/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepoConfig struct {
	userRepo      repository.CustomerRepository
	adminRepo     repository.AdminRepository
	merchantRepo  repository.MerchantRepository
	productRepo   repository.ProductRepository
	cartRepo      repository.CartRepository
	orderRepo     repository.OrderRepository
	returRepo     repository.ReturRepository
	tBuyerRepo    repository.TBuyerRepository
	tSellerRepo   repository.TSellerRepository
	tRefundRepo   repository.TRefundRepository
	rMerchantRepo repository.RMerchantRepository
	rProductRepo  repository.RProductRepository
	shippingRepo  repository.ShippingRepository
	searchRepo    repository.SearchRepository
	cityRepo      repository.CityRepository
	ongkirRepo    repository.OngkirRepository
}

func NewMongoRepo(
	db *mongo.Database,
	es *elasticsearch.Client,
	redis *redis.Client,
) RepoConfig {
	return &mongoRepoConfig{
		userRepo:      mongoRepo.NewCustomerRepository(db),
		adminRepo:     mongoRepo.NewAdminRepository(db),
		merchantRepo:  mongoRepo.NewMerchantRepository(db),
		productRepo:   mongoRepo.NewProductRepository(db),
		cartRepo:      mongoRepo.NewCartRepository(db),
		orderRepo:     mongoRepo.NewOrderRepository(db),
		returRepo:     mongoRepo.NewReturRepository(db),
		tBuyerRepo:    mongoRepo.NewTBuyerRepository(db),
		tSellerRepo:   mongoRepo.NewTSellerRepository(db),
		tRefundRepo:   mongoRepo.NewTRefundRepository(db),
		rMerchantRepo: mongoRepo.NewRMerchantRepository(db),
		rProductRepo:  mongoRepo.NewRProductRepository(db),
		shippingRepo:  mongoRepo.NewShippingRepository(db),
		searchRepo:    elasticRepo.NewElasticSearchRepository(es),
		cityRepo:      redisRepo.NewCityRepo(redis),
		ongkirRepo:    redisRepo.NewOngkirRepo(redis),
	}
}

func (mr *mongoRepoConfig) GetRepoCustomer() repository.CustomerRepository {
	return mr.userRepo
}

func (mr *mongoRepoConfig) GetRepoAdmin() repository.AdminRepository {
	return mr.adminRepo
}

func (mr *mongoRepoConfig) GetRepoMerchant() repository.MerchantRepository {
	return mr.merchantRepo
}

func (mr *mongoRepoConfig) GetRepoProduct() repository.ProductRepository {
	return mr.productRepo
}

func (mr *mongoRepoConfig) GetRepoCart() repository.CartRepository {
	return mr.cartRepo
}

func (mr *mongoRepoConfig) GetRepoOrder() repository.OrderRepository {
	return mr.orderRepo
}

func (mr *mongoRepoConfig) GetRepoRetur() repository.ReturRepository {
	return mr.returRepo
}

func (mr *mongoRepoConfig) GetRepoTBuyer() repository.TBuyerRepository {
	return mr.tBuyerRepo
}

func (mr *mongoRepoConfig) GetRepoTSeller() repository.TSellerRepository {
	return mr.tSellerRepo
}

func (mr *mongoRepoConfig) GetRepoTRefund() repository.TRefundRepository {
	return mr.tRefundRepo
}

func (mr *mongoRepoConfig) GetRepoRMerchant() repository.RMerchantRepository {
	return mr.rMerchantRepo
}

func (mr *mongoRepoConfig) GetRepoRProduct() repository.RProductRepository {
	return mr.rProductRepo
}

func (mr *mongoRepoConfig) GetRepoShipping() repository.ShippingRepository {
	return mr.shippingRepo
}

func (mr *mongoRepoConfig) GetRepoSearch() repository.SearchRepository {
	return mr.searchRepo
}

func (mr *mongoRepoConfig) GetRepoCity() repository.CityRepository {
	return mr.cityRepo
}

func (mr *mongoRepoConfig) GetRepoOngkir() repository.OngkirRepository {
	return mr.ongkirRepo
}
