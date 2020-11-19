package domain

import (
	"time"
)

type Cart struct {
	ID        string    `json:"_id" bson:"_id" validate:"required"`
	Items     []Item    `json:"items" bson:"items" validate:"unique_items"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" validate:"required"`
}

type Item struct {
	Product  DenormalizationProduct  `json:"product" bson:"product" validate:"required"`
	Merchant DenormalizationMerchant `json:"merchant" bson:"merchant" validate:"required"`
	Quantity int64                   `json:"quantity" bson:"quantity" validate:"min=1"`
	Sizes    []string                `json:"sizes" bson:"sizes"`
	Colors   []string                `json:"colors" bson:"colors"`
	Note     string                  `json:"note" bson:"note"`
	Updated  bool                    `json:"updated" bson:"updated"`
	Message  string                  `json:"message" bson:"message"`
}

type CartSearchOptions struct {
	//cart's items have item with field product.id equals to search productID keyword
	ProductID string
	//cart's items have item with field merchant.id equals to search merchantID keyword
	MerchantID string
}
