package adapterJSON

import (
	"encoding/json"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterTSellerJSON struct{}

func (a *AdapterTSellerJSON) DecodeCreateInput(input []byte) (adapter.TSellerCreateInput, error) {
	var tseller adapter.TSellerCreateInput
	if err := json.Unmarshal(input, &tseller); err != nil {
		return tseller, usecase_error.ErrBadParamInput
	}
	return tseller, nil
}

func (a *AdapterTSellerJSON) DecodeUpdateInput(input []byte) (adapter.TSellerUpdateInput, error) {
	var tseller adapter.TSellerUpdateInput
	if err := json.Unmarshal(input, &tseller); err != nil {
		return tseller, usecase_error.ErrBadParamInput
	}
	return tseller, nil
}
