package adapterJSON

import (
	"encoding/json"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterTRefundJSON struct{}

func (a *AdapterTRefundJSON) DecodeCreateInput(input []byte) (adapter.TRefundInput, error) {
	var trefund adapter.TRefundInput
	if err := json.Unmarshal(input, &trefund); err != nil {
		return trefund, usecase_error.ErrBadParamInput
	}
	return trefund, nil
}

func (a *AdapterTRefundJSON) DecodeUpdateInput(input []byte) (adapter.TRefundUpdateInput, error) {
	var trefund adapter.TRefundUpdateInput
	if err := json.Unmarshal(input, &trefund); err != nil {
		return trefund, usecase_error.ErrBadParamInput
	}
	return trefund, nil
}
