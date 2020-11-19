package domain

import "time"

const (
	PEMBAYARAN_MENUNGGU_VERIFIKASI string = "PEMBAYARAN_MENUNGGU_VERIFIKASI"
	PEMBAYARAN_SEDANG_DIVERIFIKASI string = "PEMBAYARAN_SEDANG_DIVERIFIKASI"
	PEMBAYARAN_SUCCESS             string = "PEMBAYARAN_SUCCESS"
	PEMBAYARAN_GAGAL               string = "PEMBAYARAN_GAGAL"
)

type TBuyer struct {
	ID            string    `json:"_id" bson:"_id" validate:"required"`
	CustomerID    string    `json:"customer_id" bson:"customer_id" validate:"required"`
	AdminID       string    `json:"admin_id" bson:"admin_id" validate:"required"`
	TotalTransfer int64     `json:"total_transfer" bson:"total_transfer" validate:"min=0"`
	PaymentStatus string    `json:"payment_status" bson:"payment_status" validate:"required"`
	TransferPhoto string    `json:"transfer_photo" bson:"transfer_photo"`
	Message       string    `json:"message" bson:"message"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at" validate:"required"`
}

type TBuyerSearchOptions struct {
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
	AdminID    string `json:"admin_id"`
}
