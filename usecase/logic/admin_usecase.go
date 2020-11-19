package logic

import (
	"context"
	"log"
	"os"
	"time"

	guuid "github.com/google/uuid"
	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/helper"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
)

type AdminUsecase interface {
	Create(ctx context.Context, input adapter.AdminCreateInput) (domain.Admin, error)
	GetByID(ctx context.Context, adminID string) (domain.Admin, error)
	Fetch(ctx context.Context, cursor string, num int64, options adapter.AdminSearchOptions) ([]domain.Admin, error)
	UpdateBiodata(ctx context.Context, adminBiodata adapter.AdminUpdateInput, adminID string) (domain.Admin, error)
	UploadAvatar(ctx context.Context, fileName, adminID string) (domain.Admin, error)
	UpdatePassword(ctx context.Context, input adapter.AdminUpdatePasswordInput, adminID string) (domain.Admin, error)
	AddAddress(ctx context.Context, address adapter.AdminAddressCreateInput, adminID string) (domain.Admin, error)
	UpdateAddress(ctx context.Context, input adapter.AdminAddressUpdateInput, addresID, adminID string) (domain.Admin, error)
	RemoveAddress(ctx context.Context, addressID, adminID string) (domain.Admin, error)
	DeleteOne(ctx context.Context, admin domain.Admin) (domain.Admin, error)
}

type adminUsecase struct {
	adminRepo      repository.AdminRepository
	contextTimeout time.Duration
}

func NewAdminUsecase(
	adminRepo repository.AdminRepository,
	contextTimeout time.Duration,
) AdminUsecase {
	return &adminUsecase{
		adminRepo:      adminRepo,
		contextTimeout: contextTimeout,
	}
}

func (a *adminUsecase) validate(value interface{}) error {
	log.SetOutput(os.Stdout)
	if entityErr := helper.NewValidationEntity().Validate(value); entityErr != nil {
		log.Printf("Admin validate : %s \n", entityErr)
		return entityErr
	}

	return nil
}

func (a *adminUsecase) isEmailRegistered(ctx context.Context, admin domain.Admin) (bool, error) {
	log.SetOutput(os.Stdout)

	var noCursor string = ""
	var numAdmin int64 = 1
	search := domain.AdminSearchOptions{
		Email: admin.Email,
	}
	admins, err := a.adminRepo.Fetch(ctx, noCursor, numAdmin, search)
	if err != nil {
		log.Printf("Admin validate email : %s \n", err)
		return false, err
	}

	return len(admins) == 1, nil
}

func (a *adminUsecase) Create(ctx context.Context, input adapter.AdminCreateInput) (domain.Admin, error) {
	log.SetOutput(os.Stdout)
	log.Println("Admin Create : starting!")

	var admin domain.Admin

	if input.Password != input.RePassword {
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "RePassword",
				Message: "RePassword must be equal to Password",
			},
		}
		log.Println("Admin Create :failed cause, password and repassword not equal")
		return admin, err
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

	admin.ID = guuid.New().String()
	admin.Name = input.Name
	admin.Email = input.Email
	admin.Password = input.Password
	admin.Addresses = addresses
	admin.Born = input.Born
	admin.BirthDay = input.BirthDay
	admin.Phone = input.Phone
	admin.Gender = input.Gender
	admin.Avatar = "https://storage.googleapis.com/ecommerce_s2l_assets/default-user.png"
	admin.CreatedAt = time.Now().Truncate(time.Millisecond)
	admin.UpdatedAt = time.Now().Truncate(time.Millisecond)
	if entityErr := a.validate(admin); entityErr != nil {
		log.Printf("Admin Create : failed cause, %s \n", entityErr)
		return admin, entityErr
	}

	for _, add := range addresses {
		if entityErr := a.validate(add); entityErr != nil {
			return admin, entityErr
		}
	}

	if isEmailHasReg, err := a.isEmailRegistered(ctx, admin); err != nil || isEmailHasReg {
		if isEmailHasReg {
			err := usecase_error.ErrBadEntityInput{
				usecase_error.ErrEntityField{
					Field:   "Email",
					Message: "Email is not unique",
				},
			}
			log.Println("Admin Create: failed cause, email not unique")
			return admin, err
		}
		return admin, err
	}

	hashedPassword, err := helper.NewEncription().Encrypt([]byte(admin.Password))
	if err != nil {
		log.Printf("Admin Create: failed when ecrypt pass cause, %s", err)
		return admin, err
	}
	admin.Password = hashedPassword

	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	return a.adminRepo.Create(ctx, admin)
}

func (a *adminUsecase) GetByID(ctx context.Context, adminID string) (domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	return a.adminRepo.GetByID(ctx, adminID)
}

func (a *adminUsecase) Fetch(ctx context.Context, cursor string, num int64, options adapter.AdminSearchOptions) ([]domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	search := domain.AdminSearchOptions{
		Name:  options.Name,
		Email: options.Email,
	}

	admins, err := a.adminRepo.Fetch(ctx, cursor, num, search)
	if err != nil {
		return nil, err
	}

	return admins, err
}

func (a *adminUsecase) UpdateBiodata(ctx context.Context, adminBiodata adapter.AdminUpdateInput, adminID string) (domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	admin, err := a.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return admin, err
	}

	admin.Name = adminBiodata.Name
	admin.Born = adminBiodata.Born
	admin.BirthDay = adminBiodata.BirthDay
	admin.Gender = adminBiodata.Gender
	admin.Phone = adminBiodata.Phone
	if err := a.validate(admin); err != nil {
		return admin, err
	}

	return a.adminRepo.UpdateOne(ctx, admin)
}

func (a *adminUsecase) UploadAvatar(ctx context.Context, fileName, adminID string) (domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	admin, err := a.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return admin, err
	}

	admin.Avatar = fileName
	return a.adminRepo.UpdateOne(ctx, admin)
}

func (a *adminUsecase) UpdatePassword(ctx context.Context, input adapter.AdminUpdatePasswordInput, adminID string) (domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	admin, err := a.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return admin, err
	}

	if input.Password != input.RePassword {
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "RePassword",
				Message: "RePassword must be equal to Password",
			},
		}
		return admin, err
	}

	admin.Password = input.Password
	if err := a.validate(admin); err != nil {
		return admin, err
	}

	hashedPassword, err := helper.NewEncription().Encrypt([]byte(admin.Password))
	if err != nil {
		return admin, err
	}
	admin.Password = hashedPassword

	return a.adminRepo.UpdateOne(ctx, admin)
}

func (a *adminUsecase) AddAddress(ctx context.Context, input adapter.AdminAddressCreateInput, adminID string) (domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	admin, err := a.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return admin, err
	}

	address := domain.Address{}
	address.ID = guuid.New().String()
	address.City = input.City
	address.Street = input.Street
	address.Number = input.Number
	admin.Addresses = append(admin.Addresses, address)
	if err := a.validate(address); err != nil {
		return admin, err
	}
	if err := a.validate(admin); err != nil {
		return admin, err
	}

	return a.adminRepo.UpdateOne(ctx, admin)
}

func (a *adminUsecase) UpdateAddress(ctx context.Context, input adapter.AdminAddressUpdateInput, addressID, adminID string) (domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	admin, err := a.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return admin, err
	}

	found := false
	index := 0
	for i, addr := range admin.Addresses {
		if addr.ID == addressID {
			found = true
			index = i
		}
	}
	if !found {
		return admin, usecase_error.ErrNotFound
	}
	admin.Addresses[index].City = input.City
	admin.Addresses[index].Street = input.Street
	admin.Addresses[index].Number = input.Number
	if err := a.validate(admin.Addresses[index]); err != nil {
		return admin, err
	}
	if err := a.validate(admin); err != nil {
		return admin, err
	}

	return a.adminRepo.UpdateOne(ctx, admin)
}

func (a *adminUsecase) RemoveAddress(ctx context.Context, addressID, adminID string) (domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	admin, err := a.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return admin, err
	}

	found := false
	index := 0
	for i, addr := range admin.Addresses {
		if addr.ID == addressID {
			found = true
			index = i
		}
	}
	if !found {
		return admin, usecase_error.ErrNotFound
	}

	admin.Addresses = append(admin.Addresses[:index], admin.Addresses[index+1:]...)
	return a.adminRepo.UpdateOne(ctx, admin)
}

func (a *adminUsecase) DeleteOne(ctx context.Context, admin domain.Admin) (domain.Admin, error) {
	return admin, nil
}
