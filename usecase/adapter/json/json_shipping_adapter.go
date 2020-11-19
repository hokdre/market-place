package adapterJSON

import (
	"encoding/json"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterShippingJSON struct{}

func (a *AdapterShippingJSON) DecodeCreateInput(input []byte) (adapter.ShippingCreateInput, error) {
	var shipping adapter.ShippingCreateInput
	if err := json.Unmarshal(input, &shipping); err != nil {
		return shipping, usecase_error.ErrBadParamInput
	}
	return shipping, nil
}

func (a *AdapterShippingJSON) DecodeUpdateInput(input []byte) (adapter.ShippingUpdateInput, error) {
	var shipping adapter.ShippingUpdateInput
	if err := json.Unmarshal(input, &shipping); err != nil {
		return shipping, usecase_error.ErrBadParamInput
	}
	return shipping, nil
}
