package logic

import (
	"context"
	"sync"
	"time"

	guuid "github.com/google/uuid"
	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/helper"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
)

type ReviewMerchantUsecase interface {
	Create(ctx context.Context, input adapter.RMerchantCreateInput) (domain.RMerchant, error)
	Fetch(ctx context.Context, search domain.RMerchantSearchOptions) ([]domain.RMerchant, error)
	UpdateOne(ctx context.Context, input adapter.RMerchantUpdateInput, reviewID string) (domain.RMerchant, error)
}

type reviewMerchantUsecase struct {
	reviewMerchantRepo repository.RMerchantRepository
	customerRepo       repository.CustomerRepository
	orderRepo          repository.OrderRepository
	merchantRepo       repository.MerchantRepository
	contextTimeOut     time.Duration
}

func NewRMerchantUsecase(
	reviewMerchantRepo repository.RMerchantRepository,
	customerRepo repository.CustomerRepository,
	orderRepo repository.OrderRepository,
	merchantRepo repository.MerchantRepository,
	contextTimeOut time.Duration,
) ReviewMerchantUsecase {
	return &reviewMerchantUsecase{
		reviewMerchantRepo: reviewMerchantRepo,
		customerRepo:       customerRepo,
		orderRepo:          orderRepo,
		merchantRepo:       merchantRepo,
		contextTimeOut:     contextTimeOut,
	}
}

func (r *reviewMerchantUsecase) validate(value interface{}) error {
	if entityErr := helper.NewValidationEntity().Validate(value); entityErr != nil {
		return entityErr
	}

	return nil
}

func (r *reviewMerchantUsecase) Create(ctx context.Context, input adapter.RMerchantCreateInput) (domain.RMerchant, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeOut)
	defer cancel()

	var review domain.RMerchant
	credential := ctx.Value("credential")
	if credential == nil {
		return review, usecase_error.ErrNotAuthentication
	}
	userInfo, ok := credential.(domain.Credential)
	if !ok {
		return review, usecase_error.ErrNotAuthentication
	}

	var customer domain.Customer
	var order domain.Order
	var merchant domain.Merchant

	var err error

	var wgRead sync.WaitGroup
	wgRead.Add(3)
	go func() {
		defer wgRead.Done()
		customer, err = r.customerRepo.GetByID(ctx, userInfo.UserID)
	}()
	go func() {
		defer wgRead.Done()
		order, err = r.orderRepo.GetByID(ctx, input.OrderID)
	}()
	go func() {
		defer wgRead.Done()
		merchant, err = r.merchantRepo.GetByID(ctx, input.MerchantID)
	}()

	wgRead.Wait()
	if err != nil {
		return review, err
	}

	order.ReviewedMerchant = true
	merchant.NumReview++
	merchant.Rating = (merchant.Rating + float64(input.Rating)) / float64(merchant.NumReview)

	review.ID = guuid.New().String()
	createdAt := time.Now().Truncate(time.Millisecond)
	updatedAt := time.Now().Truncate(time.Millisecond)
	review.CreatedAt = createdAt
	review.UpdatedAt = updatedAt

	review.Customer = customer.DenomalizationCustomer()
	review.MerchantID = input.MerchantID
	review.Comment = input.Comment
	review.Rating = input.Rating

	if err := r.validate(review); err != nil {
		return review, err
	}

	var wgSave sync.WaitGroup
	wgSave.Add(3)
	go func() {
		defer wgSave.Done()
		_, err = r.orderRepo.UpdateOne(ctx, order)
	}()
	go func() {
		defer wgSave.Done()
		_, err = r.merchantRepo.UpdateOne(ctx, merchant)
	}()
	go func() {
		defer wgSave.Done()
		review, err = r.reviewMerchantRepo.Create(ctx, review)
	}()

	wgSave.Wait()
	return review, err
}

func (r *reviewMerchantUsecase) Fetch(ctx context.Context, search domain.RMerchantSearchOptions) ([]domain.RMerchant, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeOut)
	defer cancel()
	return r.reviewMerchantRepo.Fetch(ctx, search.Last, 10, search)
}

func (r *reviewMerchantUsecase) UpdateOne(ctx context.Context, input adapter.RMerchantUpdateInput, reviewID string) (domain.RMerchant, error) {
	var rmerchant domain.RMerchant
	return rmerchant, nil
}
