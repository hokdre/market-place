package domain

import "time"

type Retur struct {
	ID                string           `json:"_id" bson:"_id" validate:"required"`
	OrderID           string           `json:"order_id" bson:"order_id" validate:"required"`
	CustomerReason    string           `json:"customer_reason" bson:"customer_reason" validate:"required"`
	MerchantReason    string           `json:"merchant_reason" bson:"merchant_reason"`
	MerchantAccepment string           `json:"merchant_accepment" bson:"merchant_accepment" validate:"required"`
	Shipping          ShippingProvider `json:"shipping" bson:"shipping" validate:"required"`
	ShippingCost      uint             `json:"shipping_cost" bson:"shipping_cost" validate:"min=0"`
	ShippingPhoto     string           `json:"shipping_photo" bson:"shipping_photo"`
	ResiNumber        string           `json:"resi_number" bson:"resi_number"`
	CreatedAt         time.Time        `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt         time.Time        `json:"updated_at" bson:"updated_at" validate:"required"`
}

type ReturSearchOptions struct {
	OrderID string
}
