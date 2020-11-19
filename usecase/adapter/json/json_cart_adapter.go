package adapterJSON

import (
	"encoding/json"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterCartJSON struct{}

func (a *AdapterCartJSON) DecodeAddItemInput(input []byte) (adapter.CartAddItemInput, error) {
	var item adapter.CartAddItemInput
	if err := json.Unmarshal(input, &item); err != nil {
		return item, usecase_error.ErrBadParamInput
	}
	return item, nil
}

func (a *AdapterCartJSON) DecodeUpdateInput(input []byte) (adapter.CartUpdateItemInput, error) {
	var item adapter.CartUpdateItemInput
	if err := json.Unmarshal(input, &item); err != nil {
		return item, usecase_error.ErrBadParamInput
	}
	return item, nil
}
