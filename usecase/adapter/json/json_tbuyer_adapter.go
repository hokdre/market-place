package adapterJSON

import (
	"encoding/json"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterTBuyerJSON struct{}

func (a *AdapterTBuyerJSON) DecodeCreateInput(input []byte) (adapter.TBuyerCreateInput, error) {
	var tbuyer adapter.TBuyerCreateInput
	if err := json.Unmarshal(input, &tbuyer); err != nil {
		return tbuyer, usecase_error.ErrBadParamInput
	}
	return tbuyer, nil
}

func (a *AdapterTBuyerJSON) DecodeRejectInput(input []byte) (adapter.TbuyerRejectInput, error) {
	var reject adapter.TbuyerRejectInput
	if err := json.Unmarshal(input, &reject); err != nil {
		return reject, usecase_error.ErrBadParamInput
	}
	return reject, nil
}
