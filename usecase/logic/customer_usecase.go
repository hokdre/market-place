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
)

type CustomerUsecase interface {
	Create(ctx context.Context, input adapter.CustomerCreateInput) (domain.Customer, error)
	GetByID(ctx context.Context, customerID string) (domain.Customer, error)
	Fetch(ctx context.Context, cursor string, num int64, options adapter.CustomerSearchOptions) ([]domain.Customer, error)
	UpdateBiodata(ctx context.Context, customerBiodata adapter.CustomerUpdateInput, customerID string) (domain.Customer, error)
	UploadAvatar(ctx context.Context, fileName, customerID string) (domain.Customer, error)
	UpdatePassword(ctx context.Context, input adapter.CustomerUpdatePasswordInput, customerID string) (domain.Customer, error)
	AddBankAccount(ctx context.Context, input adapter.CustomerBankCreateInput, customerID string) (domain.Customer, error)
	UpdateBankAccount(ctx context.Context, input adapter.CustomerBankUpdateInput, bankID, customerID string) (domain.Customer, error)
	AddAddress(ctx context.Context, input adapter.CustomerAddressCreateInput, customerID string) (domain.Customer, error)
	UpdateAddress(ctx context.Context, input adapter.CustomerAddressUpdateInput, addID, customerID string) (domain.Customer, error)
	RemoveAddress(ctx context.Context, addID, customerID string) (domain.Customer, error)
	DeleteOne(ctx context.Context, customer domain.Customer) (domain.Customer, error)
}

type customerUsecase struct {
	customerRepo   repository.CustomerRepository
	rmerchantRepo  repository.RMerchantRepository
	rproductRepo   repository.RProductRepository
	cartRepo       repository.CartRepository
	orderRepo      repository.OrderRepository
	contextTimeout time.Duration
}

func NewCustomerUsecase(
	customerRepo repository.CustomerRepository,
	rmerchantRepo repository.RMerchantRepository,
	rproductRepo repository.RProductRepository,
	cartRepo repository.CartRepository,
	orderRepo repository.OrderRepository,
	contextTimeout time.Duration,
) CustomerUsecase {
	return &customerUsecase{
		customerRepo:   customerRepo,
		rmerchantRepo:  rmerchantRepo,
		rproductRepo:   rproductRepo,
		cartRepo:       cartRepo,
		orderRepo:      orderRepo,
		contextTimeout: contextTimeout,
	}
}

func (c *customerUsecase) validate(value interface{}) error {
	if entityErr := helper.NewValidationEntity().Validate(value); entityErr != nil {
		return entityErr
	}

	return nil
}

func (c *customerUsecase) isEmailRegistered(ctx context.Context, customer domain.Customer) (bool, error) {
	var noCursor string = ""
	var numCustomer int64 = 1
	search := domain.CustomerSearchOptions{
		Email: customer.Email,
	}
	customers, err := c.customerRepo.Fetch(ctx, noCursor, numCustomer, search)
	if err != nil {
		return false, err
	}

	return len(customers) == 1, nil
}

func (c *customerUsecase) isCustomerExist(ctx context.Context, id string) (bool, error) {
	var registered bool

	_, err := c.customerRepo.GetByID(ctx, id)
	if err != nil {
		if err == usecase_error.ErrNotFound {
			return registered, nil
		}
		return registered, err
	}

	registered = true
	return registered, nil
}

func (c *customerUsecase) Create(ctx context.Context, input adapter.CustomerCreateInput) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	var customer domain.Customer
	var cart domain.Cart
	cart = domain.Cart{
		ID:        guuid.New().String(),
		Items:     []domain.Item{},
		CreatedAt: time.Now().Truncate(time.Millisecond),
		UpdatedAt: time.Now().Truncate(time.Millisecond),
	}

	if input.Password != input.RePassword {
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "RePassword",
				Message: "RePassword must be equal to Password",
			},
		}
		fmt.Printf("[CUSTOMER USECASE] : PASSWORD NOT EQUAL : %#v \n", err)
		return customer, err
	}
	var addresses []domain.Address
	for _, add := range input.Addresses {
		var address domain.Address
		address.ID = guuid.New().String()
		address.Street = add.Street
		address.City = add.City
		address.Number = add.Number
		addresses = append(addresses, address)
	}

	customer.ID = guuid.New().String()
	customer.CartID = cart.ID
	customer.Name = input.Name
	customer.Email = input.Email
	customer.Password = input.Password
	customer.Addresses = addresses
	customer.Born = input.Born
	customer.BirthDay = input.BirthDay
	customer.Phone = input.Phone
	customer.Gender = input.Gender
	customer.Avatar = "https://storage.googleapis.com/ecommerce_s2l_assets/default-user.png"
	customer.CreatedAt = time.Now().Truncate(time.Millisecond)
	customer.UpdatedAt = time.Now().Truncate(time.Millisecond)

	if entityErr := c.validate(customer); entityErr != nil {
		fmt.Printf("[CUSTOMER USECASE] : VALIDATE CUSTOMER ENTITY : %#v \n", entityErr)
		return customer, entityErr
	}
	for _, address := range customer.Addresses {
		if entityErr := c.validate(address); entityErr != nil {
			fmt.Printf("[CUSTOMER USECASE] : VALIDATE ADDRESS : %#v \n", entityErr)
			return customer, entityErr
		}
	}
	if isEmailHasReg, err := c.isEmailRegistered(ctx, customer); err != nil || isEmailHasReg {
		if isEmailHasReg {
			err := usecase_error.ErrBadEntityInput{
				usecase_error.ErrEntityField{
					Field:   "Email",
					Message: "Email is not unique",
				},
			}
			fmt.Printf("[CUSTOMER USECASE] : VALIDATE EMAIL NOT UNIQUE : %#v \n", err)

			return customer, err
		}
		return customer, err
	}

	hashedPassword, err := helper.NewEncription().Encrypt([]byte(customer.Password))
	if err != nil {
		return customer, err
	}
	customer.Password = hashedPassword

	var wg sync.WaitGroup
	wg.Add(2)
	var errs []error
	go func() {
		defer wg.Done()
		cart, err = c.cartRepo.Create(ctx, cart)
		if err != nil {
			errs = append(errs, err)
		}
	}()
	go func() {
		defer wg.Done()
		customer, err = c.customerRepo.Create(ctx, customer)
		if err != nil {
			errs = append(errs, err)
		}
	}()

	wg.Wait()
	return customer, nil
}

func (c *customerUsecase) Fetch(ctx context.Context, cursor string, num int64, options adapter.CustomerSearchOptions) ([]domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	search := domain.CustomerSearchOptions{
		Name:  options.Name,
		Email: options.Email,
	}
	customers, err := c.customerRepo.Fetch(ctx, cursor, num, search)
	if err != nil {
		return nil, err
	}

	return customers, err
}

func (c *customerUsecase) GetByID(ctx context.Context, id string) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	return c.customerRepo.GetByID(ctx, id)
}

func (c *customerUsecase) UpdateBiodata(ctx context.Context, customerBiodata adapter.CustomerUpdateInput, customerID string) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	customer, err := c.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return customer, err
	}

	customer.Name = customerBiodata.Name
	customer.Born = customerBiodata.Born
	customer.BirthDay = customerBiodata.BirthDay
	customer.Gender = customerBiodata.Gender
	customer.Phone = customerBiodata.Phone
	if err := c.validate(customer); err != nil {
		return customer, err
	}

	return c.customerRepo.UpdateOne(ctx, customer)
}

func (c *customerUsecase) UpdatePassword(ctx context.Context, input adapter.CustomerUpdatePasswordInput, customerID string) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	customer, err := c.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return customer, err
	}

	if input.Password != input.RePassword {
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "RePassword",
				Message: "RePassword must be equal to Password",
			},
		}
		return customer, err
	}

	customer.Password = input.Password
	if err := c.validate(customer); err != nil {
		return customer, err
	}

	hashedPassword, err := helper.NewEncription().Encrypt([]byte(customer.Password))
	if err != nil {
		return customer, err
	}
	customer.Password = hashedPassword

	return c.customerRepo.UpdateOne(ctx, customer)
}

func (c *customerUsecase) UploadAvatar(ctx context.Context, avatar, id string) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	customer, err := c.customerRepo.GetByID(ctx, id)
	if err != nil {
		return customer, err
	}

	customer.Avatar = avatar
	return c.customerRepo.UpdateOne(ctx, customer)
}

func (c *customerUsecase) AddBankAccount(ctx context.Context, input adapter.CustomerBankCreateInput, customerID string) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	customer, err := c.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return customer, err
	}

	bankAccount := domain.BankAccount{}
	bankAccount.ID = guuid.New().String()
	bankAccount.Number = input.Number
	bankAccount.BankCode = input.BankCode
	customer.BankAccounts = append(customer.BankAccounts, bankAccount)
	if err := c.validate(bankAccount); err != nil {
		return customer, err
	}
	if err := c.validate(customer); err != nil {
		return customer, err
	}
	return c.customerRepo.UpdateOne(ctx, customer)
}

func (c *customerUsecase) UpdateBankAccount(ctx context.Context, input adapter.CustomerBankUpdateInput, accountID, customerID string) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	customer, err := c.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return customer, err
	}

	found := false
	index := 0
	for i, account := range customer.BankAccounts {
		if account.ID == accountID {
			found = true
			index = i
		}
	}
	if !found {
		return customer, usecase_error.ErrNotFound
	}
	customer.BankAccounts[index].BankCode = input.BankCode
	customer.BankAccounts[index].Number = input.Number
	if err := c.validate(customer.BankAccounts[index]); err != nil {
		return customer, err
	}
	if err := c.validate(customer); err != nil {
		return customer, err
	}

	return c.customerRepo.UpdateOne(ctx, customer)
}

func (c *customerUsecase) AddAddress(ctx context.Context, input adapter.CustomerAddressCreateInput, id string) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	customer, err := c.customerRepo.GetByID(ctx, id)
	if err != nil {
		return customer, err
	}

	address := domain.Address{}
	address.ID = guuid.New().String()
	address.City = input.City
	address.Street = input.Street
	address.Number = input.Number
	customer.Addresses = append(customer.Addresses, address)
	if err := c.validate(address); err != nil {
		return customer, err
	}
	if err := c.validate(customer); err != nil {
		return customer, err
	}

	return c.customerRepo.UpdateOne(ctx, customer)
}

func (c *customerUsecase) UpdateAddress(ctx context.Context, input adapter.CustomerAddressUpdateInput, addId, customerID string) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	customer, err := c.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return customer, err
	}

	found := false
	index := 0
	for i, addr := range customer.Addresses {
		if addr.ID == addId {
			found = true
			index = i
		}
	}
	if !found {
		return customer, usecase_error.ErrNotFound
	}
	customer.Addresses[index].City = input.City
	customer.Addresses[index].Street = input.Street
	customer.Addresses[index].Number = input.Number
	if err := c.validate(customer.Addresses[index]); err != nil {
		return customer, err
	}
	if err := c.validate(customer); err != nil {
		return customer, err
	}

	return c.customerRepo.UpdateOne(ctx, customer)
}

func (c *customerUsecase) RemoveAddress(ctx context.Context, addID, customerID string) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	customer, err := c.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return customer, err
	}

	found := false
	index := 0
	for i, addr := range customer.Addresses {
		if addr.ID == addID {
			found = true
			index = i
		}
	}
	if !found {
		return customer, usecase_error.ErrNotFound
	}
	customer.Addresses = append(customer.Addresses[:index], customer.Addresses[index+1:]...)
	if err := c.validate(customer); err != nil {
		return customer, err
	}

	return c.customerRepo.UpdateOne(ctx, customer)
}

func (c *customerUsecase) DeleteOne(ctx context.Context, customer domain.Customer) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	if isExist, err := c.isCustomerExist(ctx, customer.ID); err != nil || !isExist {
		if !isExist {
			return customer, usecase_error.ErrNotFound
		}

		return customer, err
	}

	return c.customerRepo.DeleteOne(ctx, customer)
}
