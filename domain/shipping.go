package domain

import (
	"time"
)

type ShippingProvider struct {
	ID        string    `json:"_id" bson:"_id" validate:"required"`
	Name      string    `json:"name" bson:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" validate:"required"`
}

type ShippingProviderSearchOptions struct {
	Name string `json:"name" bson:"name"`
}
