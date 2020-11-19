package adapterJSON

import (
	"encoding/json"
	"fmt"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterOrderJSON struct{}

func (a *AdapterOrderJSON) DecodeCreateInput(input []byte) (adapter.OrderCreateInput, error) {
	var order adapter.OrderCreateInput
	if err := json.Unmarshal(input, &order); err != nil {
		fmt.Printf("[DEBUG] : %#v \n", err)
		return order, usecase_error.ErrBadParamInput
	}
	return order, nil
}

func (a *AdapterOrderJSON) DecodeResiInput(input []byte) (adapter.OrderResiInput, error) {
	var resi adapter.OrderResiInput
	if err := json.Unmarshal(input, &resi); err != nil {
		return resi, usecase_error.ErrBadParamInput
	}
	return resi, nil
}
