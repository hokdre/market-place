package domain

import (
	"time"
)

// user as customer
type Customer struct {
	ID           string        `json:"_id" bson:"_id"`
	CartID       string        `json:"cart_id" bson:"cart_id" validate:"required"`
	MerchantID   string        `json:"merchant_id" bson:"merchant_id"`
	Email        string        `json:"email" bson:"email" validate:"required,email"`
	Name         string        `json:"name" bson:"name" validate:"required"`
	Password     string        `json:"-" bson:"password" validate:"min=8,clower,cupper,cnumeric,csymbol"`
	Addresses    []Address     `json:"addresses" bson:"addresses" validate:"required,min=1,unique_addresses"`
	Born         string        `json:"born" bson:"born" validate:"required"`
	BirthDay     time.Time     `json:"birth_day" bson:"birth_day" validate:"required,ltfield=CreatedAt"`
	Phone        string        `json:"phone" bson:"phone" validate:"required,phone"`
	Avatar       string        `json:"avatar" bson:"avatar" validate:"required"`
	Gender       string        `json:"gender" bson:"gender" validate:"required,gender"`
	BankAccounts []BankAccount `json:"bank_accounts" bson:"bank_accounts" validate:"unique_bank_accounts"`
	Confrimed    bool          `json:"confrimed" bson:"confrimed"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
}

func (c *Customer) DenomalizationCustomer() DenomarlizationCustomer {
	return DenomarlizationCustomer{
		ID:        c.ID,
		Email:     c.Email,
		Name:      c.Name,
		Addresses: c.Addresses,
		Phone:     c.Phone,
		Avatar:    c.Avatar,
	}
}

type DenomarlizationCustomer struct {
	ID        string    `json:"_id" bson:"_id"`
	Email     string    `json:"email" bson:"email"`
	Name      string    `json:"name" bson:"name"`
	Addresses []Address `json:"addresses" bson:"addresses"`
	Phone     string    `json:"phone" bson:"phone"`
	Avatar    string    `json:"avatar" bson:"avatar"`
}

type CustomerSearchOptions struct {
	Name       string
	Email      string
	CartID     string
	MerchantID string
}
