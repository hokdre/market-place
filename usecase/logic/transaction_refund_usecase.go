package logic

import (
	"context"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/repository"
)

type TRefundUsecase interface {
	Create(ctx context.Context, input adapter.TRefundInput) (domain.TRefund, error)
	GetByID(ctx context.Context, transactionID string) (domain.TRefund, error)
	Fetch(ctx context.Context, search domain.TRefundSearchOptions) ([]domain.TRefund, error)
	UpdateOne(ctx context.Context, input adapter.TRefundUpdateInput, transactionID string) (domain.TRefund, error)
}

type trefundUsecase struct {
	trefundRepo    repository.TRefundRepository
	orderRepo      repository.OrderRepository
	customerRepo   repository.CustomerRepository
	contextTimeOut time.Duration
}

func NewTRefundUsecase(
	trefundRepo repository.TRefundRepository,
	orderRepo repository.OrderRepository,
	customerRepo repository.CustomerRepository,
	contextTimeOut time.Duration,
) TRefundUsecase {
	return &trefundUsecase{
		trefundRepo:    trefundRepo,
		orderRepo:      orderRepo,
		customerRepo:   customerRepo,
		contextTimeOut: contextTimeOut,
	}
}

func (t *trefundUsecase) Create(ctx context.Context, input adapter.TRefundInput) (domain.TRefund, error) {
	var trefund domain.TRefund
	return trefund, nil
}

func (t *trefundUsecase) GetByID(ctx context.Context, transactionID string) (domain.TRefund, error) {
	var trefund domain.TRefund
	return trefund, nil
}

func (t *trefundUsecase) Fetch(ctx context.Context, search domain.TRefundSearchOptions) ([]domain.TRefund, error) {
	var trefunds []domain.TRefund
	return trefunds, nil
}

func (t *trefundUsecase) UpdateOne(ctx context.Context, input adapter.TRefundUpdateInput, transactionID string) (domain.TRefund, error) {
	var trefund domain.TRefund
	return trefund, nil
}
