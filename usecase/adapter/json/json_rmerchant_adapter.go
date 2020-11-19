package adapterJSON

import (
	"encoding/json"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterRMerchantJSON struct{}

func (a *AdapterRMerchantJSON) DecodeCreateInput(input []byte) (adapter.RMerchantCreateInput, error) {
	var rmerchant adapter.RMerchantCreateInput
	if err := json.Unmarshal(input, &rmerchant); err != nil {
		return rmerchant, usecase_error.ErrBadParamInput
	}
	return rmerchant, nil
}

func (a *AdapterRMerchantJSON) DecodeUpdateInput(input []byte) (adapter.RMerchantUpdateInput, error) {
	var rmerchant adapter.RMerchantUpdateInput
	if err := json.Unmarshal(input, &rmerchant); err != nil {
		return rmerchant, usecase_error.ErrBadParamInput
	}
	return rmerchant, nil
}
