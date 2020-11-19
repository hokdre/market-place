package adapter

import (
	"time"

	"github.com/market-place/domain"
)

type AdminCreateInput struct {
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	RePassword string    `json:"re_password"`
	Addresses  []Address `json:"addresses"`
	Born       string    `json:"born"`
	BirthDay   time.Time `json:"birth_day"`
	Phone      string    `json:"phone"`
	Gender     string    `json:"gender"`
}

type AdminUpdateInput struct {
	Name     string    `json:"name"`
	Born     string    `json:"born"`
	BirthDay time.Time `json:"birth_day"`
	Phone    string    `json:"phone"`
	Gender   string    `json:"gender"`
}

type AdminSearchOptions struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AdminUpdatePasswordInput struct {
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}

type AdminAddressCreateInput struct {
	City       domain.City `json:"city"`
	Street     string      `json:"street"`
	Number     string      `json:"number"`
	PostalCode string      `json:"postal_code"`
}

type AdminAddressUpdateInput struct {
	City       domain.City `json:"city"`
	Street     string      `json:"street"`
	Number     string      `json:"number"`
	PostalCode string      `json:"postal_code"`
}

type AdminAdapter interface {
	DecodeCreateInput([]byte) (AdminCreateInput, error)
	DecodeUpdateInput([]byte) (AdminUpdateInput, error)
	DecodeUpdatePasswordInput([]byte) (AdminUpdatePasswordInput, error)
	DecodeAddressInput([]byte) (AdminAddressCreateInput, error)
	DecodeAddressUpdate([]byte) (AdminAddressUpdateInput, error)
}
