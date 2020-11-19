package adapter

type ProductOrder struct {
	ProductID string   `json:"product_id"`
	Quantity  int64    `json:"quantity"`
	BuyerNote string   `json:"buyer_note"`
	Colors    []string `json:"colors"`
	Sizes     []string `json:"sizes"`
}

type Order struct {
	MerchantID      string         `json:"merchant_id"`
	ReceiverName    string         `json:"receiver_name"`
	ReceiverPhone   string         `json:"receiver_phone"`
	ReceiverAddress Address        `json:"receiver_address"`
	ShippingID      string         `json:"shipping_id"`
	ShippingCost    int64          `json:"shipping_cost"`
	ServiceName     string         `json:"service_name"`
	Products        []ProductOrder `json:"products"`
}

type OrderCreateInput struct {
	Orders []Order `json:"orders"`
}

type OrderResiInput struct {
	ResiNumber string `json:"resi_number"`
}

type OrderSearchOptions struct {
	//order's customer.id equals to search customerID keyword
	CustomerID string `json:"customer_id"`
	//order's merchant.id equals to search merchantID keyword
	MerchantID string `json:"merchant_id"`

	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

type OrderAdapter interface {
	DecodeCreateInput([]byte) (OrderCreateInput, error)
	DecodeResiInput([]byte) (OrderResiInput, error)
}
