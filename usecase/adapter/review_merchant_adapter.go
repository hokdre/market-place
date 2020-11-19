package adapter

type RMerchantCreateInput struct {
	MerchantID string `json:"merchant_id"`
	OrderID    string `json:"order_id"`
	Rating     uint   `json:"rating"`
	Comment    string `json:"comment"`
}

type RMerchantUpdateInput struct {
	Rating  uint   `json:"rating"`
	Comment string `json:"comment"`
}

type RMerchantAdapter interface {
	DecodeCreateInput([]byte) (RMerchantCreateInput, error)
	DecodeUpdateInput([]byte) (RMerchantUpdateInput, error)
}
