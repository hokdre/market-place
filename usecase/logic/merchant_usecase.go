package logic

import (
	"context"
	"fmt"
	"sync"
	"time"

	guuid "github.com/google/uuid"
	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/helper"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
	"golang.org/x/sync/errgroup"
)

type MerchantUsecase interface {
	Create(ctx context.Context, input adapter.MerchantCreateInput) (domain.Merchant, error)
	GetByID(ctx context.Context, merchantID string) (domain.Merchant, error)
	Fetch(ctx context.Context, cursor string, num int64, input adapter.MerchantSearchOptions) ([]domain.Merchant, error)
	UpdateData(ctx context.Context, merchantData adapter.MerchantUpdateInput, merchantID string) (domain.Merchant, error)
	UploadAvatar(ctx context.Context, fileName, merchantID string) (domain.Merchant, error)
	AddShipping(ctx context.Context, shippingID string, merchantID string) (domain.Merchant, error)
	RemoveShipping(ctx context.Context, shippingID string, merchantID string) (domain.Merchant, error)
	AddBankAccount(ctx context.Context, input adapter.MerchantBankCreateInput, merchantID string) (domain.Merchant, error)
	UpdateBankAccount(ctx context.Context, input adapter.MerchantBankUpdateInput, accountID, merchantID string) (domain.Merchant, error)
	AddEtalase(ctx context.Context, input adapter.MerchantEtalaseCreateInput, merchantID string) (domain.Merchant, error)
	DeleteEtalase(ctx context.Context, etalaseName string, merchantID string) (domain.Merchant, error)
	UpdateAddress(ctx context.Context, address domain.Address, merchantID string) (domain.Merchant, error)
}

type merchantUsecase struct {
	merchantRepo   repository.MerchantRepository
	customerRepo   repository.CustomerRepository
	shippingRepo   repository.ShippingRepository
	productRepo    repository.ProductRepository
	orderRepo      repository.OrderRepository
	cartRepo       repository.CartRepository
	contextTimeout time.Duration
}

func NewMerchantUsecase(
	merchantRepo repository.MerchantRepository,
	customerRepo repository.CustomerRepository,
	shippingRepo repository.ShippingRepository,
	productRepo repository.ProductRepository,
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	contextTimeout time.Duration,
) MerchantUsecase {
	return &merchantUsecase{
		merchantRepo:   merchantRepo,
		customerRepo:   customerRepo,
		shippingRepo:   shippingRepo,
		productRepo:    productRepo,
		orderRepo:      orderRepo,
		cartRepo:       cartRepo,
		contextTimeout: contextTimeout,
	}

}

func (m *merchantUsecase) syncDenormalizationData(ctx context.Context, merchant domain.Merchant) error {
	var wgFind sync.WaitGroup
	var carts []domain.Cart
	var products []domain.Product

	var errs []error
	wgFind.Add(2)

	go func() {
		defer wgFind.Done()
		searchCart := domain.CartSearchOptions{
			MerchantID: merchant.ID,
		}
		c, err := m.cartRepo.Fetch(ctx, "", 0, searchCart)
		if err != nil {
			errs = append(errs, err)
			return
		}
		carts = c
	}()

	go func() {
		defer wgFind.Done()
		searchProduct := domain.ProductSearchOptions{
			MerchantID: merchant.ID,
		}
		p, err := m.productRepo.Fetch(ctx, "", 0, searchProduct)
		if err != nil {
			errs = append(errs, err)
			return
		}
		products = p
	}()

	wgFind.Wait()
	if len(errs) != 0 {
		return errs[0]
	}

	var wgSave sync.WaitGroup
	wgSave.Add(2)
	go func() {
		defer wgSave.Done()
		var wgSaveCart sync.WaitGroup
		for _, cart := range carts {

			index := -1
			for i, item := range cart.Items {
				if item.Merchant.ID == merchant.ID {
					index = i
					break
				}
			}
			if index != -1 {
				cart.Items[index].Merchant = merchant.DenomarlizationData()
				wgSaveCart.Add(1)
				go func(cart domain.Cart) {
					defer wgSaveCart.Done()
					m.cartRepo.UpdateOne(ctx, cart)
				}(cart)
			}

		}
		wgSaveCart.Wait()
	}()
	go func() {
		defer wgSave.Done()
		var wgSaveProduct sync.WaitGroup
		for _, product := range products {
			wgSaveProduct.Add(1)
			go func(product domain.Product) {
				defer wgSaveProduct.Done()
				m.productRepo.UpdateOne(ctx, product)
			}(product)
		}
		wgSaveProduct.Wait()
	}()
	wgSave.Wait()

	return nil
}

func (m *merchantUsecase) isNameUnique(ctx context.Context, merchant domain.Merchant) error {
	_, err := m.merchantRepo.GetByName(ctx, merchant.Name)
	// if err nil there's customer
	if err == nil {
		fmt.Printf("[USECASE-VALIDATION] : MERCHANT, %#v \n", " Name is Not Unique")
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "Name",
				Message: "Name is not unique",
			},
		}
		return err
	}

	if err != nil && err == usecase_error.ErrNotFound {
		return nil
	}

	return err
}

func (m *merchantUsecase) validate(value interface{}) error {
	if entityErr := helper.NewValidationEntity().Validate(value); entityErr != nil {
		return entityErr
	}

	return nil
}

func (m *merchantUsecase) Create(ctx context.Context, input adapter.MerchantCreateInput) (domain.Merchant, error) {
	var merchant domain.Merchant
	merchant.ID = guuid.New().String()

	userInfo := ctx.Value("credential")
	if userInfo == nil {
		return merchant, usecase_error.ErrNotAuthentication
	}
	credential := userInfo.(domain.Credential)

	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	owner, err := m.customerRepo.GetByID(ctx, credential.UserID)
	if err != nil {
		return merchant, err
	}
	if owner.MerchantID != "" {
		return merchant, usecase_error.ErrConflict
	}
	owner.MerchantID = merchant.ID

	shipping, err := m.shippingRepo.GetByID(ctx, input.ShippingID)
	if err != nil {
		return merchant, err
	}
	merchant.Shippings = []domain.ShippingProvider{shipping}

	merchant.Avatar = "https://storage.googleapis.com/ecommerce_s2l_assets/default-merchant.jpeg"
	merchant.CreatedAt = time.Now().Truncate(time.Millisecond)
	merchant.UpdatedAt = time.Now().Truncate(time.Millisecond)
	merchant.Name = input.Name
	merchant.Description = input.Description
	merchant.Phone = input.Phone
	merchant.Address = domain.Address{}
	merchant.Address.ID = guuid.New().String()
	merchant.Address.City = input.Address.City
	merchant.Address.Street = input.Address.Street
	merchant.Address.Number = input.Address.Number
	merchant.LocationPoint = input.LocationPoint
	if err := m.validate(merchant); err != nil {
		return merchant, err
	}
	if err := m.isNameUnique(ctx, merchant); err != nil {
		errVal, ok := err.(usecase_error.ErrBadEntityInput)
		if ok {
			return merchant, errVal
		}
		return merchant, err
	}
	if err := m.validate(merchant.Address); err != nil {
		return merchant, err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		owner, err = m.customerRepo.UpdateOne(ctx, owner)
		if err != nil {
			return err
		}

		return nil
	})
	g.Go(func() error {
		merchant, err = m.merchantRepo.Create(ctx, merchant)
		if err != nil {
			return err
		}

		return nil
	})
	if err := g.Wait(); err != nil {
		return merchant, nil
	}

	return merchant, nil
}

func (m *merchantUsecase) GetByID(ctx context.Context, merchantID string) (domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	return m.merchantRepo.GetByID(ctx, merchantID)
}

func (m *merchantUsecase) Fetch(ctx context.Context, cursor string, num int64, input adapter.MerchantSearchOptions) ([]domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	search := domain.MerchantSearchOptions{
		Name:        input.Name,
		City:        input.City,
		Description: input.Description,
	}
	customers, err := m.merchantRepo.Fetch(ctx, cursor, num, search)
	if err != nil {
		return nil, err
	}

	return customers, err
}

func (m *merchantUsecase) UpdateData(ctx context.Context, input adapter.MerchantUpdateInput, merchantID string) (domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	merchant, err := m.merchantRepo.GetByID(ctx, merchantID)
	if err != nil {
		return merchant, err
	}

	merchant.Phone = input.Phone
	merchant.Description = input.Description
	merchant.Address.City = input.Address.City
	merchant.Address.Street = input.Address.Street
	merchant.Address.Number = input.Address.Number
	merchant.LocationPoint = input.LocationPoint
	if err := m.validate(merchant.Address); err != nil {
		return merchant, err
	}
	if err := m.validate(merchant); err != nil {
		return merchant, err
	}

	var wgSave sync.WaitGroup
	wgSave.Add(2)
	go func() {
		defer wgSave.Done()
		m.syncDenormalizationData(ctx, merchant)
	}()

	var updatedMerchant domain.Merchant
	go func() {
		defer wgSave.Done()
		updatedMerchant, err = m.merchantRepo.UpdateOne(ctx, merchant)
	}()
	wgSave.Wait()

	return updatedMerchant, err
}

func (m *merchantUsecase) UploadAvatar(ctx context.Context, fileName, merchantID string) (domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	merchant, err := m.merchantRepo.GetByID(ctx, merchantID)
	if err != nil {
		return merchant, err
	}
	merchant.Avatar = fileName

	var wgSave sync.WaitGroup
	wgSave.Add(2)
	go func() {
		defer wgSave.Done()
		m.syncDenormalizationData(ctx, merchant)
	}()

	var updatedMerchant domain.Merchant
	go func() {
		defer wgSave.Done()
		updatedMerchant, err = m.merchantRepo.UpdateOne(ctx, merchant)
	}()
	wgSave.Wait()

	return updatedMerchant, err
}

func (m *merchantUsecase) AddShipping(ctx context.Context, shippingID string, merchantID string) (domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	merchant, err := m.merchantRepo.GetByID(ctx, merchantID)
	if err != nil {
		return merchant, err
	}
	shipping, err := m.shippingRepo.GetByID(ctx, shippingID)
	if err != nil {
		return merchant, err
	}

	merchant.Shippings = append(merchant.Shippings, shipping)
	if err := m.validate(shipping); err != nil {
		return merchant, err
	}
	if err := m.validate(merchant); err != nil {
		return merchant, err
	}

	var wgSave sync.WaitGroup
	wgSave.Add(2)
	go func() {
		defer wgSave.Done()
		m.syncDenormalizationData(ctx, merchant)
	}()

	var updatedMerchant domain.Merchant
	go func() {
		defer wgSave.Done()
		updatedMerchant, err = m.merchantRepo.UpdateOne(ctx, merchant)
	}()
	wgSave.Wait()

	return updatedMerchant, err
}

func (m *merchantUsecase) RemoveShipping(ctx context.Context, shippingID string, merchantID string) (domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	merchant, err := m.merchantRepo.GetByID(ctx, merchantID)
	if err != nil {
		return merchant, err
	}

	found := false
	index := 0
	for i, shipping := range merchant.Shippings {
		if shipping.ID == shippingID {
			found = true
			index = i
		}
	}
	if !found {
		return merchant, usecase_error.ErrNotFound
	}
	if len(merchant.Shippings) == 1 {
		fmt.Printf("[USECASE-VALIDATION] : MERCHANT, %#v \n", " Shipping Is Empty")
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "Shippings",
				Message: "Shippings is must be filled",
			},
		}
		return merchant, err
	}

	merchant.Shippings = append(merchant.Shippings[:index], merchant.Shippings[index+1:]...)

	var wgSave sync.WaitGroup
	wgSave.Add(2)
	go func() {
		defer wgSave.Done()
		m.syncDenormalizationData(ctx, merchant)
	}()

	var updatedMerchant domain.Merchant
	go func() {
		defer wgSave.Done()
		updatedMerchant, err = m.merchantRepo.UpdateOne(ctx, merchant)
	}()
	wgSave.Wait()

	return updatedMerchant, err
}

func (m *merchantUsecase) AddBankAccount(ctx context.Context, input adapter.MerchantBankCreateInput, merchantID string) (domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	merchant, err := m.merchantRepo.GetByID(ctx, merchantID)
	if err != nil {
		return merchant, err
	}

	bankAccount := domain.BankAccount{}
	bankAccount.ID = guuid.New().String()
	bankAccount.Number = input.Number
	bankAccount.BankCode = input.BankCode
	merchant.BankAccounts = append(merchant.BankAccounts, bankAccount)
	if err := m.validate(bankAccount); err != nil {
		return merchant, err
	}
	if err := m.validate(merchant); err != nil {
		return merchant, err
	}
	return m.merchantRepo.UpdateOne(ctx, merchant)
}

func (m *merchantUsecase) UpdateBankAccount(ctx context.Context, input adapter.MerchantBankUpdateInput, accountID, merchantID string) (domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	merchant, err := m.merchantRepo.GetByID(ctx, merchantID)
	if err != nil {
		return merchant, err
	}

	found := false
	index := 0
	for i, account := range merchant.BankAccounts {
		if account.ID == accountID {
			found = true
			index = i
		}
	}
	if !found {
		return merchant, usecase_error.ErrNotFound
	}
	merchant.BankAccounts[index].BankCode = input.BankCode
	merchant.BankAccounts[index].Number = input.Number
	if err := m.validate(merchant.BankAccounts[index]); err != nil {
		return merchant, err
	}
	if err := m.validate(merchant); err != nil {
		return merchant, err
	}

	return m.merchantRepo.UpdateOne(ctx, merchant)
}

func (m *merchantUsecase) AddEtalase(ctx context.Context, input adapter.MerchantEtalaseCreateInput, merchantID string) (domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	merchant, err := m.merchantRepo.GetByID(ctx, merchantID)
	if err != nil {
		return merchant, err
	}
	if input.Name == "" {
		return merchant, usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "Etalase",
				Message: "Etalase must be filled",
			},
		}
	}

	merchant.Etalase = append(merchant.Etalase, input.Name)
	if err := m.validate(merchant); err != nil {
		return merchant, err
	}
	return m.merchantRepo.UpdateOne(ctx, merchant)
}

func (m *merchantUsecase) DeleteEtalase(ctx context.Context, etalaseName string, merchantID string) (domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	merchant, err := m.merchantRepo.GetByID(ctx, merchantID)
	if err != nil {
		return merchant, err
	}

	found := false
	index := 0
	for i, etalase := range merchant.Etalase {
		if etalase == etalaseName {
			found = true
			index = i
		}
	}
	if !found {
		return merchant, usecase_error.ErrNotFound
	}
	if index < len(merchant.Etalase) {
		merchant.Etalase = append(merchant.Etalase[:index], merchant.Etalase[index+1:]...)
	} else {
		merchant.Etalase = merchant.Etalase[:index]
	}

	noCursor := ""
	num := int64(1)
	products, err := m.productRepo.Fetch(ctx, noCursor, num, domain.ProductSearchOptions{
		MerchantID: merchantID,
		Etalase:    etalaseName,
	})
	if len(products) != 0 {
		return merchant, usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "Etalase",
				Message: "Etalase still have products",
			},
		}
	}

	return m.merchantRepo.UpdateOne(ctx, merchant)
}

func (m *merchantUsecase) UpdateAddress(ctx context.Context, address domain.Address, merchantID string) (domain.Merchant, error) {
	var merchant domain.Merchant
	return merchant, nil
}
