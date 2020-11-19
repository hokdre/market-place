package logic

import (
	"context"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/repository"
)

type ReturUseCase interface {
	Create(ctx context.Context, input adapter.ReturCreateInput) (domain.Retur, error)
	GetByID(ctx context.Context, returID string) (domain.Retur, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.ReturSearchOptions) ([]domain.Retur, error)
	AcceptRetur(ctx context.Context, returID string) (domain.Retur, error)
	RejectRetur(ctx context.Context, input adapter.ReturRejectInput, returID string) (domain.Retur, error)
	InputShipping(ctx context.Context, input adapter.ReturShippingInput, returID string) (domain.Retur, error)
	UploadShippingPhoto(ctx context.Context, fileName string, returID string) (domain.Retur, error)
}

type returUsecase struct {
	returRepo      repository.ReturRepository
	shippingRepo   repository.ShippingRepository
	contextTimeOut time.Duration
}

func NewReturUsecase(
	returRepo repository.ReturRepository,
	shippingRepo repository.ShippingRepository,
	contextTimeOut time.Duration,
) ReturUseCase {
	return &returUsecase{
		returRepo:      returRepo,
		shippingRepo:   shippingRepo,
		contextTimeOut: contextTimeOut,
	}
}

func (r *returUsecase) Create(ctx context.Context, input adapter.ReturCreateInput) (domain.Retur, error) {
	var retur domain.Retur
	return retur, nil
}

func (r *returUsecase) GetByID(ctx context.Context, returID string) (domain.Retur, error) {
	var retur domain.Retur
	return retur, nil
}

func (r *returUsecase) Fetch(ctx context.Context, cursor string, num int64, options domain.ReturSearchOptions) ([]domain.Retur, error) {
	var returs []domain.Retur
	return returs, nil
}

func (r *returUsecase) AcceptRetur(ctx context.Context, returID string) (domain.Retur, error) {
	var retur domain.Retur
	return retur, nil
}

func (r *returUsecase) RejectRetur(ctx context.Context, input adapter.ReturRejectInput, returID string) (domain.Retur, error) {
	var retur domain.Retur
	return retur, nil
}

func (r *returUsecase) InputShipping(ctx context.Context, input adapter.ReturShippingInput, returID string) (domain.Retur, error) {
	var retur domain.Retur
	return retur, nil
}

func (r *returUsecase) UploadShippingPhoto(ctx context.Context, fileName string, returID string) (domain.Retur, error) {
	var retur domain.Retur
	return retur, nil
}
