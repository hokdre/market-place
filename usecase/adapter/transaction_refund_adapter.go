package adapter

type TRefundInput struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Message    string `json:"message"`
}

type TRefundUpdateInput struct {
	Message string `json:"message"`
}

type TRefundAdapter interface {
	DecodeCreateInput([]byte) (TRefundInput, error)
	DecodeUpdateInput([]byte) (TRefundUpdateInput, error)
}
