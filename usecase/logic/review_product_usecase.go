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

type ReviewProductUsecase interface {
	Create(ctx context.Context, input adapter.RProductCreateInput) (domain.RProduct, error)
	Fetch(ctx context.Context, search domain.RProductSearchOptions) ([]domain.RProduct, error)
	UpdateOne(ctx context.Context, input adapter.RProductUpdateInput, reviewID string) (domain.RProduct, error)
}

type reviewProductUsecase struct {
	reviewProductRepo repository.RProductRepository
	orderRepo         repository.OrderRepository
	customerRepo      repository.CustomerRepository
	productRepo       repository.ProductRepository
	contextTimeOut    time.Duration
}

func NewRProductUsecase(
	reviewProductRepo repository.RProductRepository,
	orderRepo repository.OrderRepository,
	customerRepo repository.CustomerRepository,
	productRepo repository.ProductRepository,
	contextTimeOut time.Duration,
) ReviewProductUsecase {
	return &reviewProductUsecase{
		reviewProductRepo: reviewProductRepo,
		orderRepo:         orderRepo,
		customerRepo:      customerRepo,
		productRepo:       productRepo,
		contextTimeOut:    contextTimeOut,
	}
}

func (r *reviewProductUsecase) validate(value interface{}) error {
	if entityErr := helper.NewValidationEntity().Validate(value); entityErr != nil {
		return entityErr
	}

	return nil
}

func (r *reviewProductUsecase) Create(ctx context.Context, input adapter.RProductCreateInput) (domain.RProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeOut)
	defer cancel()

	var review domain.RProduct
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
	var product domain.Product

	var errGroup []error

	var wgRead sync.WaitGroup
	wgRead.Add(3)
	go func() {
		defer wgRead.Done()
		c, e := r.customerRepo.GetByID(ctx, userInfo.UserID)
		if e != nil {
			errGroup = append(errGroup, e)
			return
		}
		customer = c
	}()
	go func() {
		defer wgRead.Done()
		o, e := r.orderRepo.GetByID(ctx, input.OrderID)
		if e != nil {
			errGroup = append(errGroup, e)
			return
		}
		order = o
	}()
	go func() {
		defer wgRead.Done()
		p, e := r.productRepo.GetByID(ctx, input.ProductID)
		if e != nil {
			errGroup = append(errGroup, e)
			return
		}
		product = p
	}()

	wgRead.Wait()
	if len(errGroup) != 0 {
		return review, errGroup[0]
	}

	product.NumReview++
	product.Rating = (product.Rating + float64(input.Rating)) / float64(product.NumReview)

	review.ID = guuid.New().String()
	createdAt := time.Now().Truncate(time.Millisecond)
	updatedAt := time.Now().Truncate(time.Millisecond)
	review.CreatedAt = createdAt
	review.UpdatedAt = updatedAt

	review.Customer = customer.DenomalizationCustomer()
	review.ProductID = input.ProductID
	review.Comment = input.Comment
	review.Rating = input.Rating

	if err := r.validate(review); err != nil {
		return review, err
	}

	var wgSave sync.WaitGroup
	wgSave.Add(3)
	var err error
	go func() {
		defer wgSave.Done()
		r.orderRepo.UpdateOne(ctx, order)
	}()
	go func() {
		defer wgSave.Done()
		r.productRepo.UpdateOne(ctx, product)
	}()
	go func() {
		defer wgSave.Done()
		review, err = r.reviewProductRepo.Create(ctx, review)
	}()
	if !order.ReviewedProduct {
		order.ReviewedProduct = true
		wgSave.Add(1)
		go func() {
			defer wgSave.Done()
			_, err = r.orderRepo.UpdateOne(ctx, order)
		}()
	}

	wgSave.Wait()
	return review, err
}

func (r *reviewProductUsecase) Fetch(ctx context.Context, search domain.RProductSearchOptions) ([]domain.RProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeOut)
	defer cancel()
	return r.reviewProductRepo.Fetch(ctx, search.Last, 10, search)
}

func (r *reviewProductUsecase) UpdateOne(ctx context.Context, input adapter.RProductUpdateInput, reviewID string) (domain.RProduct, error) {
	var rproduct domain.RProduct
	return rproduct, nil
}
