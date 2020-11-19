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

type AuthenticationUsecase interface {
	LoginCustomer(context.Context, adapter.LoginInput) (domain.Credential, error)
	LoginAdmin(context.Context, adapter.LoginInput) (domain.Credential, error)
	ValidateLogin(token string) (domain.Credential, error)
	VerifiedAsCustomer(domain.Credential) error
	VerifiedAsAdmin(domain.Credential) error
	VerifiedCustomerAuthor(domain.Credential, string) error
	VerifiedAdminAuthor(domain.Credential, string) error
	VerifiedMerchantOwner(domain.Credential, string) error
	VerifiedProductOwner(context.Context, domain.Credential, string) error
	VerifiedCartOwner(domain.Credential, string) error
	VerifiedCustomerTransactionOwner(context.Context, domain.Credential, string) error
	VerifiedMerchantOrderOwner(context.Context, domain.Credential, string) error
	VerifiedCustomerOrderOwner(context.Context, domain.Credential, string) error
}

type authenticationUseCase struct {
	customerRepo   repository.CustomerRepository
	adminRepo      repository.AdminRepository
	productRepo    repository.ProductRepository
	tbuyerRepo     repository.TBuyerRepository
	orderRepo      repository.OrderRepository
	contextTimeout time.Duration
}

func NewAuthenticationUseCase(
	customerRepo repository.CustomerRepository,
	adminRepo repository.AdminRepository,
	productRepo repository.ProductRepository,
	tbuyerRepo repository.TBuyerRepository,
	orderRepo repository.OrderRepository,
	contextTimeout time.Duration,
) AuthenticationUsecase {
	return &authenticationUseCase{
		customerRepo:   customerRepo,
		adminRepo:      adminRepo,
		productRepo:    productRepo,
		tbuyerRepo:     tbuyerRepo,
		orderRepo:      orderRepo,
		contextTimeout: contextTimeout,
	}
}

func (a *authenticationUseCase) ValidateLogin(token string) (domain.Credential, error) {
	credential, err := helper.DecodeToken(token)
	if err != nil {
		fmt.Printf("[AUTHENTICATION] : VALIDATE LOGIN %#v \n", err)
		return credential, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	switch credential.LoginType {
	case domain.LOGIN_AS_ADMIN:
		if err := a.validateAdminActive(ctx, &credential); err != nil {
			fmt.Printf("[AUTHENTICATION] : VALIDATE ADMIN ACTIVE %#v \n", err)
			return credential, err
		}
	case domain.LOGIN_AS_CUSTOMER:
		if err := a.validateCustomerActive(ctx, &credential); err != nil {
			fmt.Printf("[AUTHENTICATION] : VALIDATE CUSTOMER ACTIVE %#v \n", err)
			return credential, err
		}
	}

	return credential, err
}

func (a *authenticationUseCase) validateCustomerActive(ctx context.Context, credential *domain.Credential) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	customer, err := a.customerRepo.GetByID(ctx, credential.UserID)
	if err != nil {
		return err
	}

	credential.UserID = customer.ID
	credential.MerchantID = customer.MerchantID
	credential.CartID = customer.CartID
	credential.Email = customer.Email
	credential.LoginType = domain.LOGIN_AS_CUSTOMER

	return nil
}

func (a *authenticationUseCase) validateAdminActive(ctx context.Context, credential *domain.Credential) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	admin, err := a.adminRepo.GetByID(ctx, credential.UserID)
	if err != nil {
		return usecase_error.ErrNotAuthorization
	}

	credential.UserID = admin.ID
	credential.MerchantID = ""
	credential.CartID = ""
	credential.Email = admin.Email
	credential.LoginType = domain.LOGIN_AS_ADMIN

	return nil
}

func (a *authenticationUseCase) LoginCustomer(ctx context.Context, input adapter.LoginInput) (domain.Credential, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	if input.Email == "" {
		fmt.Printf("[AUTHENTICATION] : LOGIN CUSTOMER %#v \n", "EMAIL EMPTY")
		return domain.Credential{}, usecase_error.ErrLoginField{
			Field:   "Email",
			Message: "Email is not registered",
		}
	}

	search := domain.CustomerSearchOptions{
		Email: input.Email,
	}
	var noCursor string
	var numReturned int64 = 1
	customers, err := a.customerRepo.Fetch(ctx, noCursor, numReturned, search)
	if err != nil {
		fmt.Printf("[AUTHENTICATION] : LOGIN CUSTOMER %#v \n", err)
		return domain.Credential{}, usecase_error.ErrInternalServerError
	}
	if len(customers) != 1 {
		fmt.Printf("[AUTHENTICATION] : LOGIN CUSTOMER %#v \n", "EMAIL NOT REGISTERED")
		return domain.Credential{}, usecase_error.ErrLoginField{
			Field:   "Email",
			Message: "Email is not registered",
		}
	}

	customer := customers[0]
	isMatch := helper.NewEncription().Compare(
		[]byte(customer.Password),
		[]byte(input.Password),
	)
	if !isMatch {
		fmt.Printf("[AUTHENTICATION] : LOGIN CUSTOMER %#v \n", "PASSWORD NOT MATCH")
		return domain.Credential{}, usecase_error.ErrLoginField{
			Field:   "Password",
			Message: "Password is wrong",
		}
	}

	credential := domain.NewCredential(
		customer.ID,
		customer.CartID,
		customer.MerchantID,
		customer.Email,
		domain.LOGIN_AS_CUSTOMER,
	)

	return credential, nil
}

func (a *authenticationUseCase) LoginAdmin(ctx context.Context, input adapter.LoginInput) (domain.Credential, error) {

	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	if input.Email == "" {
		fmt.Printf("[AUTHENTICATION] : LOGIN ADMIN %#v \n", "EMAIL EMPTY")

		return domain.Credential{}, usecase_error.ErrLoginField{
			Field:   "Email",
			Message: "Email is not registered",
		}
	}

	search := domain.AdminSearchOptions{
		Email: input.Email,
	}
	var noCursor string
	var numReturned int64 = 1
	admins, err := a.adminRepo.Fetch(ctx, noCursor, numReturned, search)
	if err != nil {
		fmt.Printf("[AUTHENTICATION] : LOGIN ADMIN FETCH %#v \n", err)
		return domain.Credential{}, usecase_error.ErrInternalServerError
	}
	if len(admins) != 1 {
		fmt.Printf("[AUTHENTICATION] : LOGIN ADMIN %#v \n", "EMAIL NOT REGISTERED")
		return domain.Credential{}, usecase_error.ErrLoginField{
			Field:   "Email",
			Message: "Email is not registered",
		}
	}

	customer := admins[0]
	isMatch := helper.NewEncription().Compare(
		[]byte(customer.Password),
		[]byte(input.Password),
	)
	if !isMatch {
		fmt.Printf("[AUTHENTICATION] : LOGIN ADMIN %#v \n", "PASSWORD NOT MATCH")
		return domain.Credential{}, usecase_error.ErrLoginField{
			Field:   "Password",
			Message: "Password is wrong",
		}
	}

	credential := domain.NewCredential(
		customer.ID,
		"",
		"",
		customer.Email,
		domain.LOGIN_AS_ADMIN,
	)
	return credential, nil
}

func (a *authenticationUseCase) VerifiedAsCustomer(credential domain.Credential) error {
	if credential.LoginType != domain.LOGIN_AS_CUSTOMER {
		fmt.Printf("[AUTHENTICATION] :VALIDATE AS CUSTOMER %#v \n", "LOGIN TYPE NOT AS CUSTOMER")
		return usecase_error.ErrNotAuthorization
	}
	return nil
}

func (a *authenticationUseCase) VerifiedAsAdmin(credential domain.Credential) error {
	if credential.LoginType != domain.LOGIN_AS_ADMIN {
		fmt.Printf("[AUTHENTICATION] :VALIDATE AS ADMIN %#v \n", "LOGIN TYPE NOT AS ADMIN")
		return usecase_error.ErrNotAuthorization
	}
	return nil
}

func (a *authenticationUseCase) VerifiedCustomerAuthor(credential domain.Credential, customerID string) error {
	if credential.UserID != customerID {
		fmt.Printf("[AUTHENTICATION] :VALIDATE CUSTOMER AUTHOR %#v \n", "CREDENTIAL NOT VALID")
		return usecase_error.ErrNotAuthorization
	}
	return nil
}

func (a *authenticationUseCase) VerifiedAdminAuthor(credential domain.Credential, adminID string) error {
	if credential.UserID != adminID {
		fmt.Printf("[AUTHENTICATION] :VALIDATE ADMIN AUTHOR %#v \n", "CREDENTIAL NOT VALID")
		return usecase_error.ErrNotAuthorization
	}
	return nil
}

func (a *authenticationUseCase) VerifiedMerchantOwner(credential domain.Credential, merchantID string) error {
	if credential.MerchantID != merchantID {
		fmt.Printf("[AUTHENTICATION] :VALIDATE MERCHANT OWNER %#v \n", "CREDENTIAL NOT VALID")
		return usecase_error.ErrNotAuthorization
	}
	return nil
}

func (a *authenticationUseCase) VerifiedProductOwner(ctx context.Context, credential domain.Credential, productID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if credential.MerchantID == "" || productID == "" {
		fmt.Printf("[AUTHENTICATION] :VALIDATE PRODUCT OWNER %#v \n", "MERCHANT ID OR PRODUCT ID EMPTY")
		return usecase_error.ErrNotAuthorization
	}

	product, err := a.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err != usecase_error.ErrNotFound {
			fmt.Printf("[AUTHENTICATION] :VALIDATE PRODUCT OWNER %#v \n", "PRODUCT NOT FOUND")
			return usecase_error.ErrNotAuthorization
		}

		return err
	}
	if product.Merchant.ID != credential.MerchantID {
		return usecase_error.ErrNotAuthorization
	}

	return nil
}

func (a *authenticationUseCase) VerifiedCartOwner(credential domain.Credential, cartID string) error {
	if credential.CartID != cartID {
		fmt.Printf("[AUTHENTICATION] :VALIDATE CART OWNER %#v \n", "CREDENTIAL NOT VALID")
		return usecase_error.ErrNotAuthorization
	}

	return nil

}

func (a *authenticationUseCase) VerifiedCustomerTransactionOwner(ctx context.Context, credential domain.Credential, transactionID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	tbuyer, err := a.tbuyerRepo.GetByID(ctx, transactionID)
	if err != nil {
		if err != usecase_error.ErrNotFound {
			fmt.Printf("[AUTHENTICATION] :VALIDATE TRANSACTION BUYER OWNER %#v \n", "TRANSACTION BUYER NOT FOUND")
			return usecase_error.ErrNotAuthorization
		}

		return err
	}
	if credential.UserID != tbuyer.CustomerID {
		fmt.Printf("[AUTHENTICATION] :VALIDATE TRANSACTION BUYER OWNER %#v \n", "CREDENTIAL NOT VALID")
		return usecase_error.ErrNotAuthorization
	}

	return nil
}

func (a *authenticationUseCase) VerifiedMerchantOrderOwner(ctx context.Context, credential domain.Credential, orderID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	order, err := a.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if err != usecase_error.ErrNotFound {
			fmt.Printf("[AUTHENTICATION] :VALIDATE MERCHANT ORDER OWNER %#v \n", "ORDER NOT FOUND")
			return usecase_error.ErrNotAuthorization
		}

		return err
	}

	if credential.MerchantID != order.Merchant.ID {
		fmt.Printf("[AUTHENTICATION] :VALIDATE MERCHANT ORDER OWNER %#v \n", "CREDENTIAL NOT VALID")
		return usecase_error.ErrNotAuthorization
	}

	return nil
}

func (a *authenticationUseCase) VerifiedCustomerOrderOwner(ctx context.Context, credential domain.Credential, orderID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	order, err := a.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if err != usecase_error.ErrNotFound {
			fmt.Printf("[AUTHENTICATION] :VALIDATE MERCHANT ORDER OWNER %#v \n", "ORDER NOT FOUND")
			return usecase_error.ErrNotAuthorization
		}

		return err
	}

	if credential.UserID != order.Customer.ID {
		fmt.Printf("[AUTHENTICATION] :VALIDATE MERCHANT ORDER OWNER %#v \n", "CREDENTIAL NOT VALID")
		return usecase_error.ErrNotAuthorization
	}

	return nil
}
