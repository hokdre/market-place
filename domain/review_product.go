package domain

import "time"

type RProduct struct {
	ID        string                  `json:"_id" bson:"_id" validate:"required"`
	ProductID string                  `json:"product_id" bson:"product_id" validate:"required"`
	Customer  DenomarlizationCustomer `json:"customer" bson:"customer" validate:"required"`
	Rating    uint                    `json:"rating" bson:"rating" validate:"required,min=1,max=5"`
	Comment   string                  `json:"comment" bson:"comment" validate:"required"`
	CreatedAt time.Time               `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt time.Time               `json:"updated_at" bson:"updated_at" validate:"required"`
}

type RProductSearchOptions struct {
	ProductID string `json:"product_id"`
	Last      string `json:"last"`
}
