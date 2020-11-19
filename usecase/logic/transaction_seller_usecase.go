package logic

import (
	"context"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/repository"
)

type TSellerUsecase interface {
	Create(ctx context.Context, input adapter.TSellerCreateInput) (domain.TSeller, error)
	GetByID(ctx context.Context, transactionID string) (domain.TSeller, error)
	Fetch(ctx context.Context, search domain.TSellerSearchOptions) ([]domain.TSeller, error)
	UpdateOne(ctx context.Context, input adapter.TSellerUpdateInput, transactionID string) (domain.TSeller, error)
}

type tsellerUsecase struct {
	orderRepo      repository.OrderRepository
	customerRepo   repository.CustomerRepository
	merchantRepo   repository.MerchantRepository
	contextTimeOut time.Duration
}

func NewTSellerUsecase(
	orderRepo repository.OrderRepository,
	customerRepo repository.CustomerRepository,
	merchantRepo repository.MerchantRepository,
	contextTimeOut time.Duration,
) TSellerUsecase {
	return &tsellerUsecase{
		orderRepo:      orderRepo,
		customerRepo:   customerRepo,
		merchantRepo:   merchantRepo,
		contextTimeOut: contextTimeOut,
	}
}

func (t *tsellerUsecase) Create(ctx context.Context, input adapter.TSellerCreateInput) (domain.TSeller, error) {
	var tseller domain.TSeller
	return tseller, nil
}

func (t *tsellerUsecase) GetByID(ctx context.Context, transactionID string) (domain.TSeller, error) {
	var tseller domain.TSeller
	return tseller, nil
}

func (t *tsellerUsecase) Fetch(ctx context.Context, search domain.TSellerSearchOptions) ([]domain.TSeller, error) {
	var tsellers []domain.TSeller
	return tsellers, nil
}

func (t *tsellerUsecase) UpdateOne(ctx context.Context, input adapter.TSellerUpdateInput, transactionID string) (domain.TSeller, error) {
	var tseller domain.TSeller
	return tseller, nil
}
