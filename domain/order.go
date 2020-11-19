package domain

import (
	"time"
)

const (
	STATUS_ORDER_MENUNGGU_PEMBAYARAN = "STATUS_ORDER_MENUNGGU_PEMBAYARAN"
	STATUS_ORDER_SEDANG_DIPROSES     = "STATUS_ORDER_SEDANG_DIPROSES"
	STATUS_ORDER_SEDANG_DIKIRIM      = "STATUS_ORDER_SEDANG_DIKIRIM"
	STATUS_ORDER_DI_CANCEL           = "STATUS_ORDER_DI_CANCEL"
	/*
	* Next Rilis Features :)
	*
	 */
	// STATUS_ORDER_DI_AJUKAN_RETUR     = "STATUS_ORDER_DI_AJUKAN_RETUR"
	// STATUS_ORDER_RETUR_DITERIMA      = "STATUS_ORDER_RETUR_DITERIMA"
	// STATUS_ORDER_RETUR_DITOLAK       = "STATUS_ORDER_RETUR_DITOLAK"
	// STATUS_ORDER_DIPUTUSKAN          = "STATUS_ORDER_DIPUTUSKAN"
	STATUS_ORDER_SELESAI = "STATUS_ORDER_SELESAI"
)

type OrderItems struct {
	Product   DenormalizationProduct `json:"product"`
	Quantity  int64                  `json:"quantity" bson:"quantity" validate:"min=1"`
	BuyerNote string                 `json:"buyer_note" bson:"buyer_note"`
	Colors    []string               `json:"colors" bson:"colors"`
	Sizes     []string               `json:"sizes" bson:"sizes"`
	Price     int64                  `json:"price" bson:"price" validate:"min=0"`
}

type Order struct {
	ID               string                  `json:"_id" bson:"_id" validate:"required"`
	TransactionsID   string                  `json:"transaction_id" bson:"transaction_id"`
	OrderItems       []OrderItems            `json:"order_items" bson:"order_items" validate:"required"`
	Merchant         DenormalizationMerchant `json:"merchant" bson:"merchant" validate:"required"`
	Customer         DenomarlizationCustomer `json:"customer" bson:"customer" validate:"required"`
	ReceiverName     string                  `json:"receiver_name" bson:"receiver_name" validate:"required"`
	ReceiverPhone    string                  `json:"receiver_phone" bson:"receiver_phone"`
	ReceiverAddress  Address                 `json:"receiver_address" bson:"receiver_address"`
	Shipping         ShippingProvider        `json:"shipping" bson:"shipping" validate:"required"`
	ShippingCost     int64                   `json:"shipping_cost" bson:"shipping_cost" validate:"min=0"`
	ServiceName      string                  `json:"service_name" bson:"service_name"`
	StatusOrder      string                  `json:"status_order" bson:"status_order"`
	ResiNumber       string                  `json:"resi_number" bson:"resi_number"`
	ShippingPhoto    string                  `json:"shipping_photo" bson:"shipping_photo"`
	ReviewedMerchant bool                    `json:"reviewed_merchant" bson:"reviewed_merchant"`
	ReviewedProduct  bool                    `json:"reviewed_product" bson:"reviewed_product"`
	Delivered        bool                    `json:"delivered" bson:"delivered"`
	CreatedAt        time.Time               `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt        time.Time               `json:"updated_at" bson:"updated_at" validate:"required"`
}

type OrderSearchOptions struct {
	//order's customer.id equals to search customerID keyword
	CustomerID string
	//order's merchant.id equals to search merchantID keyword
	MerchantID string
	//order's shipping.id equals to search shippingID keyword
	ShippingID string
	//order's product.id equals to search productID keyword
	ProductID     string
	TransactionID string
	Status        string
}
