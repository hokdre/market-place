package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/helper"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
)

type CartUsecase interface {
	GetByID(ctx context.Context, cartID string) (domain.Cart, error)
	AddProduct(ctx context.Context, input adapter.CartAddItemInput, cartID string) (domain.Cart, error)
	UpdateItemInCart(ctx context.Context, itemData adapter.CartUpdateItemInput, productID string, cartID string) (domain.Cart, error)
	RemoveProduct(ctx context.Context, productID string, cartID string) (domain.Cart, error)
	ClearProduct(ctx context.Context, cartID string) (domain.Cart, error)
}

type cartUsecase struct {
	cartRepo       repository.CartRepository
	productRepo    repository.ProductRepository
	contextTimeout time.Duration
}

func NewCartUsecase(
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	contextTimeout time.Duration,
) CartUsecase {
	return &cartUsecase{
		cartRepo:       cartRepo,
		productRepo:    productRepo,
		contextTimeout: contextTimeout,
	}
}

func (c *cartUsecase) validate(value interface{}) error {
	if entityErr := helper.NewValidationEntity().Validate(value); entityErr != nil {
		return entityErr
	}

	return nil
}

func (c *cartUsecase) GetByID(ctx context.Context, cartID string) (domain.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	return c.cartRepo.GetByID(ctx, cartID)
}

func (c *cartUsecase) AddProduct(ctx context.Context, input adapter.CartAddItemInput, cartID string) (domain.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	cart, err := c.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		return cart, err
	}
	product, err := c.productRepo.GetByID(ctx, input.ProductID)
	if err != nil {
		return cart, err
	}

	credential := ctx.Value("credential")
	if credential == nil {
		return cart, usecase_error.ErrNotAuthorization
	}
	userInfo := credential.(domain.Credential)
	if userInfo.MerchantID == product.Merchant.ID {
		fmt.Printf("[USECASE-VALIDATION] : CART  %#v \n", "ADD OWN PRODUCT")
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "Items",
				Message: "Cannot Add product that is owned by buyer",
			},
		}
		return cart, err
	}

	item := domain.Item{}
	item.Product = product.DenormalizationData()
	item.Merchant = product.Merchant
	item.Note = input.Note
	item.Quantity = input.Quantity
	item.Colors = input.Colors
	item.Sizes = input.Sizes

	indexItemInCart := -1
	for index, itemInCart := range cart.Items {
		if item.Product.ID == itemInCart.Product.ID {
			indexItemInCart = index
			break
		}
	}
	if indexItemInCart != -1 {
		cart.Items[indexItemInCart].Quantity++
	} else {
		cart.Items = append(cart.Items, item)
	}

	if err := c.validate(item); err != nil {
		return cart, err
	}
	if err := c.validate(cart); err != nil {
		return cart, err
	}

	return c.cartRepo.UpdateOne(ctx, cart)
}

func (c *cartUsecase) UpdateItemInCart(ctx context.Context, input adapter.CartUpdateItemInput, productID string, cartID string) (domain.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	cart, err := c.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		return cart, err
	}

	index := 0
	found := false
	for i, item := range cart.Items {
		if item.Product.ID == productID {
			index = i
			found = true
		}
	}
	if !found {
		return cart, usecase_error.ErrNotFound
	}

	item := cart.Items[index]
	product, err := c.productRepo.GetByID(ctx, item.Product.ID)
	if err != nil {
		return cart, err
	}
	if int64(product.Stock) < input.Quantity {
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "Quantity",
				Message: "Stock product tidak mencukupi",
			},
		}
		fmt.Printf("[CART USECASE] : product'stock less than quantity:  %#v \n", err)
		return cart, err

	}
	item.Quantity = input.Quantity
	item.Colors = input.Colors
	item.Sizes = input.Sizes
	item.Note = input.Note

	cart.Items[index] = item
	if err := c.validate(item); err != nil {
		return cart, err
	}
	if err := c.validate(cart); err != nil {
		return cart, err
	}

	return c.cartRepo.UpdateOne(ctx, cart)
}

func (c *cartUsecase) RemoveProduct(ctx context.Context, productID string, cartID string) (domain.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	cart, err := c.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		return cart, err
	}

	index := 0
	found := false
	for i, item := range cart.Items {
		if item.Product.ID == productID {
			index = i
			found = true
		}
	}
	if !found {
		return cart, usecase_error.ErrNotFound
	}

	if len(cart.Items) == 1 {
		cart.Items = []domain.Item{}
	} else {
		cart.Items = append(cart.Items[:index], cart.Items[index+1:]...)
	}

	if err := c.validate(cart); err != nil {
		return cart, err
	}

	return c.cartRepo.UpdateOne(ctx, cart)
}

func (c *cartUsecase) ClearProduct(ctx context.Context, cartID string) (domain.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	cart, err := c.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		return cart, err
	}

	cart.Items = []domain.Item{}
	return c.cartRepo.UpdateOne(ctx, cart)
}
