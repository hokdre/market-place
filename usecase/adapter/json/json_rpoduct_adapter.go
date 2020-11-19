package adapterJSON

import (
	"encoding/json"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterRProductJSON struct{}

func (a *AdapterRProductJSON) DecodeCreateInput(input []byte) (adapter.RProductCreateInput, error) {
	var rProduct adapter.RProductCreateInput
	if err := json.Unmarshal(input, &rProduct); err != nil {
		return rProduct, usecase_error.ErrBadParamInput
	}
	return rProduct, nil
}

func (a *AdapterRProductJSON) DecodeUpdateInput(input []byte) (adapter.RMerchantUpdateInput, error) {
	var rProduct adapter.RMerchantUpdateInput
	if err := json.Unmarshal(input, &rProduct); err != nil {
		return rProduct, usecase_error.ErrBadParamInput
	}
	return rProduct, nil
}
