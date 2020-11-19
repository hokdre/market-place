package adapter

type TBuyerCreateInput struct {
	TotalTransfer uint `json:"total_transfer"`
}

type TbuyerRejectInput struct {
	Message string `json:"message"`
}

type TBuyerAdapter interface {
	DecodeCreateInput([]byte) (TBuyerCreateInput, error)
	DecodeRejectInput([]byte) (TbuyerRejectInput, error)
}
