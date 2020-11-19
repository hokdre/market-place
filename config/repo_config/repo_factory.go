package config

import (
	config "github.com/market-place/config/infrastructure_config"
	"github.com/market-place/usecase/repository"
)

const (
	MONGODB = "MONGODB"
)

type RepoConfig interface {
	GetRepoCustomer() repository.CustomerRepository
	GetRepoAdmin() repository.AdminRepository
	GetRepoMerchant() repository.MerchantRepository
	GetRepoProduct() repository.ProductRepository
	GetRepoShipping() repository.ShippingRepository
	GetRepoCart() repository.CartRepository
	GetRepoOrder() repository.OrderRepository
	GetRepoRetur() repository.ReturRepository
	GetRepoTBuyer() repository.TBuyerRepository
	GetRepoTSeller() repository.TSellerRepository
	GetRepoTRefund() repository.TRefundRepository
	GetRepoRMerchant() repository.RMerchantRepository
	GetRepoRProduct() repository.RProductRepository
	GetRepoSearch() repository.SearchRepository
	GetRepoCity() repository.CityRepository
	GetRepoOngkir() repository.OngkirRepository
}

func NewRepoConfig(
	infrastructureConfig config.InfrastructureConfig,
) RepoConfig {

	return NewMongoRepo(
		infrastructureConfig.GetMongoDBDatabase(),
		infrastructureConfig.GetElasticClient(),
		infrastructureConfig.GetRedisClient(),
	)

}
