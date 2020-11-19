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

type ProductUsecase interface {
	Create(ctx context.Context, input adapter.ProductCreateInput) (domain.Product, error)
	GetByID(ctx context.Context, productID string) (domain.Product, error)
	Fetch(ctx context.Context, cursor string, num int64, input adapter.ProductSearchOptions) ([]domain.Product, error)
	UpdateData(ctx context.Context, input adapter.ProductUpdateInput, productID string) (domain.Product, error)
	UploadPhotos(ctx context.Context, fileNames []string, productID string) (domain.Product, error)
	DeleteOne(ctx context.Context, productID string) (domain.Product, error)
	ProductTerlaris(ctx context.Context) ([]map[string]interface{}, error)
}

type productUsecase struct {
	productRepo    repository.ProductRepository
	merchantRepo   repository.MerchantRepository
	cartRepo       repository.CartRepository
	orderRepo      repository.OrderRepository
	contextTimeOut time.Duration
}

func NewProductUsecase(
	productRepo repository.ProductRepository,
	merchantRepo repository.MerchantRepository,
	cartRepo repository.CartRepository,
	orderRepo repository.OrderRepository,
	contextTimeOut time.Duration,
) ProductUsecase {
	return &productUsecase{
		productRepo:    productRepo,
		merchantRepo:   merchantRepo,
		cartRepo:       cartRepo,
		orderRepo:      orderRepo,
		contextTimeOut: contextTimeOut,
	}
}

func (p *productUsecase) validate(value interface{}) error {
	if entityErr := helper.NewValidationEntity().Validate(value); entityErr != nil {
		return entityErr
	}

	return nil
}

func (p *productUsecase) Create(ctx context.Context, input adapter.ProductCreateInput) (domain.Product, error) {
	var product domain.Product

	credential := ctx.Value("credential")
	if credential == nil {
		return product, usecase_error.ErrNotAuthentication
	}
	userInfo, ok := credential.(domain.Credential)
	if !ok {
		return product, usecase_error.ErrNotAuthentication
	}

	ctx, cancel := context.WithTimeout(ctx, p.contextTimeOut)
	defer cancel()

	merchant, err := p.merchantRepo.GetByID(ctx, userInfo.MerchantID)
	if err != nil {
		return product, err
	}

	product.ID = guuid.New().String()
	product.Merchant = merchant.DenomarlizationData()
	product.Name = input.Name
	product.Weight = input.Weight
	product.Width = input.Width
	product.Height = input.Height
	product.Long = input.Long
	product.Description = input.Description
	product.Category = input.Category
	product.Etalase = input.Etalase
	product.Tags = input.Tags
	product.Colors = input.Colors
	product.Sizes = input.Sizes
	product.Photos = []string{"https://storage.googleapis.com/ecommerce_s2l_assets/default-product.jpg"}
	product.Price = input.Price
	product.Stock = input.Stock
	product.CreatedAt = time.Now().Truncate(time.Millisecond)
	product.UpdatedAt = time.Now().Truncate(time.Millisecond)

	if entityErr := p.validate(product); entityErr != nil {
		return product, entityErr
	}

	product, err = p.productRepo.Create(ctx, product)
	if err != nil {
		return product, err
	}

	return product, err
}

func (p *productUsecase) GetByID(ctx context.Context, productID string) (domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeOut)
	defer cancel()
	return p.productRepo.GetByID(ctx, productID)
}

func (p *productUsecase) Fetch(ctx context.Context, cursor string, num int64, input adapter.ProductSearchOptions) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeOut)
	defer cancel()

	search := domain.ProductSearchOptions{
		Name:        input.Name,
		Category:    input.Category,
		Description: input.Description,
		MerchantID:  input.MerchantID,
		Price:       input.Price,
		City:        input.City,
	}

	products, err := p.productRepo.Fetch(ctx, cursor, num, search)
	if err != nil {
		return nil, err
	}

	return products, err
}

func (p *productUsecase) UpdateData(ctx context.Context, input adapter.ProductUpdateInput, productID string) (domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeOut)
	defer cancel()
	product, err := p.productRepo.GetByID(ctx, productID)
	if err != nil {
		return product, err
	}
	merchant, err := p.merchantRepo.GetByID(ctx, product.Merchant.ID)
	if err != nil {
		return product, err
	}

	product.Name = input.Name
	product.Weight = input.Weight
	product.Width = input.Width
	product.Height = input.Height
	product.Long = input.Long
	product.Description = input.Description
	product.Category = input.Category
	product.Tags = input.Tags
	product.Etalase = input.Etalase
	product.Colors = input.Colors
	product.Sizes = input.Sizes
	product.Price = input.Price
	product.Stock = input.Stock
	product.UpdatedAt = time.Now().Truncate(time.Millisecond)

	index := -1
	for i, mProduct := range merchant.Products {
		if mProduct.ID == product.ID {
			index = i
			break
		}
	}
	if index != -1 {
		merchant.Products[index] = product
	}

	if entityErr := p.validate(product); entityErr != nil {
		return product, entityErr
	}
	if merchant, err = p.merchantRepo.UpdateOne(ctx, merchant); err != nil {
		return product, err
	}
	return p.productRepo.UpdateOne(ctx, product)
}

func (p *productUsecase) UploadPhotos(ctx context.Context, fileNames []string, productID string) (domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeOut)
	defer cancel()

	product, err := p.productRepo.GetByID(ctx, productID)
	if err != nil {
		return product, err
	}

	product.Photos = fileNames
	return p.productRepo.UpdateOne(ctx, product)
}

func (p *productUsecase) DeleteOne(ctx context.Context, productID string) (domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	product, err := p.productRepo.GetByID(ctx, productID)
	if err != nil {
		return product, err
	}

	merchant, err := p.merchantRepo.GetByID(ctx, product.Merchant.ID)
	if err != nil {
		return product, err
	}

	index := -1
	for i, mProduct := range merchant.Products {
		if mProduct.ID == product.ID {
			index = i
			break
		}
	}

	if index != -1 {
		if index < len(merchant.Products)-1 {
			merchant.Products = append(merchant.Products[:index], merchant.Products[index+1:]...)
		} else {
			merchant.Products = merchant.Products[:index]
		}
	}

	_, err = p.merchantRepo.UpdateOne(ctx, merchant)
	if err != nil {
		return product, err
	}

	return p.productRepo.DeleteOne(ctx, product)
}

func (p *productUsecase) ProductTerlaris(ctx context.Context) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	productsData, err := p.orderRepo.ProductTerlaris(ctx)
	if err != nil {
		return productsData, err
	}

	var wgFetchUpdatedProductData sync.WaitGroup
	wgFetchUpdatedProductData.Add(len(productsData))

	for index, productData := range productsData {
		go func(index int, productData map[string]interface{}) {
			defer wgFetchUpdatedProductData.Done()
			productMap := productData["product"].(map[string]interface{})
			productID := productMap["_id"].(string)
			product, err := p.productRepo.GetByID(ctx, productID)
			if err != nil {
				return
			}
			productsData[index]["product"] = product.DenormalizationData()
		}(index, productData)
	}

	wgFetchUpdatedProductData.Wait()

	return productsData, nil
}
