package adapter

type ReturCreateInput struct {
	OrderID        string `json:"order_id"`
	CustomerReason string `json:"customer_reason"`
}

type ReturRejectInput struct {
	MerchantReason string `json:"merchant_reason"`
}

type ReturShippingInput struct {
	ShippingID string `json:"shipping_id"`
	ResiNumber string `json:"resi_number"`
}

type ReturAdapter interface {
	DecodeCreateInput([]byte) (ReturCreateInput, error)
	DecodeRejectInput([]byte) (ReturRejectInput, error)
	DecodeShippingInput([]byte) (ReturShippingInput, error)
}
