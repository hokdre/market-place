package domain

import (
	"time"
)

type Admin struct {
	ID        string    `json:"_id" bson:"_id" validate:"required"`
	Email     string    `json:"email" bson:"email" validate:"required,email"`
	Name      string    `json:"name" bson:"name" validate:"required"`
	Password  string    `json:"-" bson:"password" validate:"min=8,clower,cupper,cnumeric,csymbol"`
	Addresses []Address `json:"addresses" bson:"addresses" validate:"required,min=1,unique_addresses"`
	Born      string    `json:"born" bson:"born" validate:"required"`
	BirthDay  time.Time `json:"birth_day" bson:"birth_day" validate:"required,ltfield=CreatedAt"`
	Phone     string    `json:"phone" bson:"phone" validate:"required,phone"`
	Avatar    string    `json:"avatar" bson:"avatar" validate:"required"`
	Gender    string    `json:"gender" bson:"gender" validate:"required,gender"`
	Confrimed bool      `json:"confrimed" bson:"confrimed"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type AdminSearchOptions struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
