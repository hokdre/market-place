package adapterJSON

import (
	"encoding/json"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterReturJSON struct{}

func (a *AdapterReturJSON) DecodeCreateInput(input []byte) (adapter.ReturCreateInput, error) {
	var retur adapter.ReturCreateInput
	if err := json.Unmarshal(input, &retur); err != nil {
		return retur, usecase_error.ErrBadParamInput
	}
	return retur, nil
}

func (a *AdapterReturJSON) DecodeRejectInput(input []byte) (adapter.ReturRejectInput, error) {
	var reject adapter.ReturRejectInput
	if err := json.Unmarshal(input, &reject); err != nil {
		return reject, usecase_error.ErrBadParamInput
	}
	return reject, nil
}

func (a *AdapterReturJSON) DecodeShippingInput(input []byte) (adapter.ReturShippingInput, error) {
	var shipping adapter.ReturShippingInput
	if err := json.Unmarshal(input, &shipping); err != nil {
		return shipping, usecase_error.ErrBadParamInput
	}
	return shipping, nil
}
