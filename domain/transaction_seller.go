package domain

import "time"

type TSeller struct {
	ID            string    `json:"_id" bson:"_id" validate:"required"`
	OrderID       string    `json:"order_id" bson:"order_id" validate:"required"`
	MerchantID    string    `json:"merchant_id" bson:"merchant_id" validate:"required"`
	AdminID       string    `json:"admin_id" bson:"admin_id" validate:"required"`
	TotalTransfer uint      `json:"total_transfer" bson:"total_transfer" validate:"min=0"`
	Message       string    `json:"message" validate:"required"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at" validate:"required"`
}

type TSellerSearchOptions struct {
	OrderID    string `json:"order_id"`
	MerchantID string `json:"merchant_id"`
	AdminID    string `json:"admin_id"`
}
