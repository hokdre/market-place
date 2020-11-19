package logic

import (
	"context"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/helper"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
)

type ShippingUsecase interface {
	Create(ctx context.Context, input adapter.ShippingCreateInput) (domain.ShippingProvider, error)
	GetByID(ctx context.Context, shippingID string) (domain.ShippingProvider, error)
	Fetch(ctx context.Context, cursor string, num int64, options adapter.ShippingProviderSearchOptions) ([]domain.ShippingProvider, error)
	UpdateOne(ctx context.Context, input adapter.ShippingUpdateInput, shippingID string) (domain.ShippingProvider, error)
	DeleteOne(ctx context.Context, shippingID string) (domain.ShippingProvider, error)
}

type shippingUsecase struct {
	shippingRepo   repository.ShippingRepository
	contextTimeOut time.Duration
}

func NewShippingUsecase(
	shippingRepo repository.ShippingRepository,
	contextTimeOut time.Duration,
) ShippingUsecase {
	return &shippingUsecase{
		shippingRepo:   shippingRepo,
		contextTimeOut: contextTimeOut,
	}
}

func (c *shippingUsecase) validate(value interface{}) error {
	if entityErr := helper.NewValidationEntity().Validate(value); entityErr != nil {
		return entityErr
	}

	return nil
}

func (a *shippingUsecase) isNameRegistered(ctx context.Context, shipping domain.ShippingProvider) (bool, error) {
	var noCursor string = ""
	var numCustomer int64 = 1
	search := domain.ShippingProviderSearchOptions{
		Name: shipping.Name,
	}
	admins, err := a.shippingRepo.Fetch(ctx, noCursor, numCustomer, search)
	if err != nil {
		return false, err
	}

	return len(admins) == 1, nil
}

func (s *shippingUsecase) Create(ctx context.Context, input adapter.ShippingCreateInput) (domain.ShippingProvider, error) {
	var shipping domain.ShippingProvider
	shipping.ID = input.ID
	shipping.Name = input.Name
	shipping.CreatedAt = time.Now().Truncate(time.Millisecond)
	shipping.UpdatedAt = time.Now().Truncate(time.Millisecond)

	if err := s.validate(shipping); err != nil {
		return shipping, err
	}
	if isNameRegistered, err := s.isNameRegistered(ctx, shipping); err != nil || isNameRegistered {
		if isNameRegistered {
			err := usecase_error.ErrBadEntityInput{
				usecase_error.ErrEntityField{
					Field:   "Name",
					Message: "Name is not unique",
				},
			}
			return shipping, err
		}
		return shipping, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeOut)
	defer cancel()
	return s.shippingRepo.Create(ctx, shipping)
}

func (s *shippingUsecase) GetByID(ctx context.Context, shippingID string) (domain.ShippingProvider, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeOut)
	defer cancel()
	return s.shippingRepo.GetByID(ctx, shippingID)
}

func (s *shippingUsecase) Fetch(ctx context.Context, cursor string, num int64, options adapter.ShippingProviderSearchOptions) ([]domain.ShippingProvider, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeOut)
	defer cancel()

	search := domain.ShippingProviderSearchOptions{
		Name: options.Name,
	}
	shippings, err := s.shippingRepo.Fetch(ctx, cursor, num, search)
	if err != nil {
		return nil, err
	}

	return shippings, err
}

func (s *shippingUsecase) UpdateOne(ctx context.Context, input adapter.ShippingUpdateInput, shippingID string) (domain.ShippingProvider, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeOut)
	defer cancel()
	shipping, err := s.shippingRepo.GetByID(ctx, shippingID)
	if err != nil {
		return shipping, err
	}
	shipping.Name = input.Name

	if err := s.validate(shipping); err != nil {
		return shipping, err
	}
	if isNameRegistered, err := s.isNameRegistered(ctx, shipping); err != nil || isNameRegistered {
		if isNameRegistered {
			err := usecase_error.ErrBadEntityInput{
				usecase_error.ErrEntityField{
					Field:   "Name",
					Message: "Name is not unique",
				},
			}
			return shipping, err
		}
		return shipping, err
	}

	return s.shippingRepo.UpdateOne(ctx, shipping)
}

func (s *shippingUsecase) DeleteOne(ctx context.Context, shippingID string) (domain.ShippingProvider, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeOut)
	defer cancel()
	shipping, err := s.shippingRepo.GetByID(ctx, shippingID)
	if err != nil {
		return shipping, err
	}
	return s.shippingRepo.DeleteOne(ctx, shipping)
}
