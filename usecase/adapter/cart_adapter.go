package adapter

type CartAddItemInput struct {
	ProductID string   `json:"product_id"`
	Quantity  int64    `json:"quantity"`
	Note      string   `json:"note"`
	Colors    []string `json:"colors"`
	Sizes     []string `json:"sizes"`
}

type CartUpdateItemInput struct {
	Quantity int64    `json:"quantity"`
	Colors   []string `json:"colors"`
	Sizes    []string `json:"sizes"`
	Note     string   `json:"note"`
}

type CartAdapter interface {
	DecodeAddItemInput(input []byte) (CartAddItemInput, error)
	DecodeUpdateItemInput(input []byte) (CartUpdateItemInput, error)
}
