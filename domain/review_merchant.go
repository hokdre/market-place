package domain

import (
	"time"
)

type RMerchant struct {
	ID         string                  `json:"_id" bson:"_id" validate:"required"`
	MerchantID string                  `json:"merchant_id" bson:"merchant_id" validate:"required"`
	Customer   DenomarlizationCustomer `json:"customer" bson:"customer" validate:"required"`
	Rating     uint                    `json:"rating" bson:"rating" validate:"required,min=1,max=5"`
	Comment    string                  `json:"comment" bson:"comment" validate:"required"`
	CreatedAt  time.Time               `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt  time.Time               `json:"updated_at" bson:"updated_at" validate:"required"`
}

type RMerchantSearchOptions struct {
	MerchantID string `json:"merchant_id"`
	Last       string `json:"last"`
}
