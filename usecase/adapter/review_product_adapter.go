package adapter

type RProductCreateInput struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Rating    uint   `json:"rating"`
	Comment   string `json:"comment"`
}

type RProductUpdateInput struct {
	Rating  uint   `json:"rating"`
	Comment string `json:"comment"`
}

type RProductAdapter interface {
	DecodeCreateInput([]byte) (RProductCreateInput, error)
	DecodeUpdateInput([]byte) (RProductUpdateInput, error)
}
