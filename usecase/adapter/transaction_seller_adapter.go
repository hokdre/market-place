package adapter

type TSellerCreateInput struct {
	OrderID    string `json:"order_id"`
	MerchantID string `json:"merchant_id"`
	AdminID    string `json:"admin_id"`
	Message    string `json:"message"`
}

type TSellerUpdateInput struct {
	Message string `json:"message"`
}

type TSellerAdapter interface {
	DecodeCreateInput([]byte) (TSellerCreateInput, error)
	DecodeUpdateInput([]byte) (TSellerUpdateInput, error)
}
