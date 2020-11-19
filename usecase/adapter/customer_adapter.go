package adapter

import (
	"time"

	"github.com/market-place/domain"
)

type CustomerCreateInput struct {
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	RePassword string    `json:"re_password"`
	Addresses  []Address `json:"addresses" `
	Born       string    `json:"born" `
	BirthDay   time.Time `json:"birth_day"`
	Phone      string    `json:"phone"`
	Gender     string    `json:"gender"`
}

type CustomerUpdateInput struct {
	Name     string    `json:"name" bson:"name" validate:"required"`
	Born     string    `json:"born" bson:"born" validate:"required"`
	BirthDay time.Time `json:"birth_day" bson:"birth_day" validate:"ltfield=CreatedAt"`
	Phone    string    `json:"phone" bson:"phone" validate:"required,phone"`
	Gender   string    `json:"gender" bson:"gender" validate:"required,gender"`
}

type CustomerAddressCreateInput struct {
	City       domain.City `json:"city"`
	Street     string      `json:"street"`
	Number     string      `json:"number"`
	PostalCode string      `json:"postal_code"`
}

type CustomerAddressUpdateInput struct {
	City       domain.City `json:"city"`
	Street     string      `json:"street"`
	Number     string      `json:"number"`
	PostalCode string      `json:"postal_code"`
}

type CustomerUpdatePasswordInput struct {
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}

type CustomerBankCreateInput struct {
	Number   string `json:"number"`
	BankCode string `json:"bank_code"`
}

type CustomerBankUpdateInput struct {
	Number   string `json:"number"`
	BankCode string `json:"bank_code"`
}

type CustomerSearchOptions struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CustomerAdapter interface {
	DecodeCreateInput([]byte) (CustomerCreateInput, error)
	DecodeUpdateInput([]byte) (CustomerUpdateInput, error)
	DecodeUpdatePassword([]byte) (CustomerUpdatePasswordInput, error)
	DecodeBankInput([]byte) (CustomerBankCreateInput, error)
	DecodeBankUpdate([]byte) (CustomerBankUpdateInput, error)
	DecodeAddressInput([]byte) (CustomerAddressCreateInput, error)
	DecodeAddressUpdate([]byte) (CustomerAddressUpdateInput, error)
}
