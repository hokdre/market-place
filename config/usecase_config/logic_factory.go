package config

import (
	"time"

	repoConfig "github.com/market-place/config/repo_config"
	"github.com/market-place/usecase/logic"
)

const (
	contextTimeOut = 2 * time.Second
)

type UsecaseConfig interface {
	GetCustomersUseCase() logic.CustomerUsecase
	GetAdminsUseCase() logic.AdminUsecase
	GetMerchantUseCase() logic.MerchantUsecase
	GetProductUseCase() logic.ProductUsecase
	GetCartUseCase() logic.CartUsecase
	GetOrderUseCase() logic.OrderUsecase
	GetReturUseCase() logic.ReturUseCase
	GetTBuyerUseCase() logic.TBuyerUsecase
	GetTSellerUseCase() logic.TSellerUsecase
	GetTRefundUseCase() logic.TRefundUsecase
	GetRMerchantUseCase() logic.ReviewMerchantUsecase
	GetRProductUseCase() logic.ReviewProductUsecase
	GetAuthUsecase() logic.AuthenticationUsecase
	GetShippingUsecase() logic.ShippingUsecase
	GetSearchUsecase() logic.SearchUsecase
	GetCityUsecase() logic.CityUsecase
	GetOngkirUsecase() logic.OngkirUsecase
}

type usecaseConfig struct {
	repoConfig repoConfig.RepoConfig
}

func NewUsecaseConfig(repoConfig repoConfig.RepoConfig) UsecaseConfig {
	return &usecaseConfig{
		repoConfig: repoConfig,
	}
}

func (l *usecaseConfig) GetCustomersUseCase() logic.CustomerUsecase {
	return logic.NewCustomerUsecase(
		l.repoConfig.GetRepoCustomer(),
		l.repoConfig.GetRepoRMerchant(),
		l.repoConfig.GetRepoRProduct(),
		l.repoConfig.GetRepoCart(),
		l.repoConfig.GetRepoOrder(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetAdminsUseCase() logic.AdminUsecase {
	return logic.NewAdminUsecase(
		l.repoConfig.GetRepoAdmin(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetMerchantUseCase() logic.MerchantUsecase {
	return logic.NewMerchantUsecase(
		l.repoConfig.GetRepoMerchant(),
		l.repoConfig.GetRepoCustomer(),
		l.repoConfig.GetRepoShipping(),
		l.repoConfig.GetRepoProduct(),
		l.repoConfig.GetRepoOrder(),
		l.repoConfig.GetRepoCart(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetProductUseCase() logic.ProductUsecase {
	return logic.NewProductUsecase(
		l.repoConfig.GetRepoProduct(),
		l.repoConfig.GetRepoMerchant(),
		l.repoConfig.GetRepoCart(),
		l.repoConfig.GetRepoOrder(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetCartUseCase() logic.CartUsecase {
	return logic.NewCartUsecase(
		l.repoConfig.GetRepoCart(),
		l.repoConfig.GetRepoProduct(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetOrderUseCase() logic.OrderUsecase {
	return logic.NewOrderUsecase(
		l.repoConfig.GetRepoOrder(),
		l.repoConfig.GetRepoProduct(),
		l.repoConfig.GetRepoCustomer(),
		l.repoConfig.GetRepoMerchant(),
		l.repoConfig.GetRepoCart(),
		l.repoConfig.GetRepoTBuyer(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetReturUseCase() logic.ReturUseCase {
	return logic.NewReturUsecase(
		l.repoConfig.GetRepoRetur(),
		l.repoConfig.GetRepoShipping(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetTBuyerUseCase() logic.TBuyerUsecase {
	return logic.NewTBuyerUsecase(
		l.repoConfig.GetRepoTBuyer(),
		l.repoConfig.GetRepoOrder(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetTSellerUseCase() logic.TSellerUsecase {
	return logic.NewTSellerUsecase(
		l.repoConfig.GetRepoOrder(),
		l.repoConfig.GetRepoCustomer(),
		l.repoConfig.GetRepoMerchant(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetTRefundUseCase() logic.TRefundUsecase {
	return logic.NewTRefundUsecase(
		l.repoConfig.GetRepoTRefund(),
		l.repoConfig.GetRepoOrder(),
		l.repoConfig.GetRepoCustomer(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetRMerchantUseCase() logic.ReviewMerchantUsecase {
	return logic.NewRMerchantUsecase(
		l.repoConfig.GetRepoRMerchant(),
		l.repoConfig.GetRepoCustomer(),
		l.repoConfig.GetRepoOrder(),
		l.repoConfig.GetRepoMerchant(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetRProductUseCase() logic.ReviewProductUsecase {
	return logic.NewRProductUsecase(
		l.repoConfig.GetRepoRProduct(),
		l.repoConfig.GetRepoOrder(),
		l.repoConfig.GetRepoCustomer(),
		l.repoConfig.GetRepoProduct(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetAuthUsecase() logic.AuthenticationUsecase {
	return logic.NewAuthenticationUseCase(
		l.repoConfig.GetRepoCustomer(),
		l.repoConfig.GetRepoAdmin(),
		l.repoConfig.GetRepoProduct(),
		l.repoConfig.GetRepoTBuyer(),
		l.repoConfig.GetRepoOrder(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetShippingUsecase() logic.ShippingUsecase {
	return logic.NewShippingUsecase(
		l.repoConfig.GetRepoShipping(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetSearchUsecase() logic.SearchUsecase {
	return logic.NewSearchUsecase(
		l.repoConfig.GetRepoSearch(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetCityUsecase() logic.CityUsecase {
	return logic.NewCityUsecase(
		l.repoConfig.GetRepoCity(),
		contextTimeOut,
	)
}

func (l *usecaseConfig) GetOngkirUsecase() logic.OngkirUsecase {
	return logic.NewOngkirUsecase(
		l.repoConfig.GetRepoOngkir(),
		contextTimeOut,
	)
}
