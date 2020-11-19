package domain

import "time"

type TRefund struct {
	ID            string    `json:"_id" bson:"_id" validate:"required"`
	OrderID       string    `json:"order_id" bson:"order_id" validate:"required"`
	CustomerID    string    `json:"customer_id" bson:"customer_id" validate:"required"`
	AdminID       string    `json:"admin_id" bson:"admin_id" validate:"required"`
	TotalTransfer uint      `json:"total_transfer" bson:"total_transfer" validate:"min=0"`
	Message       string    `json:"message" bson:"message" validate:"required"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at" validate:"required"`
}

type TRefundSearchOptions struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	AdminID    string `json:"admin_id"`
}
