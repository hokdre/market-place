package adapterJSON

import (
	"encoding/json"
	"fmt"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterProductJSON struct{}

func (a *AdapterProductJSON) DecodeCreateInput(input []byte) (adapter.ProductCreateInput, error) {
	var product adapter.ProductCreateInput
	if err := json.Unmarshal(input, &product); err != nil {
		fmt.Printf("[JSON-PRODUCT-ADAPTER] : DECODE CREATE INPUT %#v \n", err)
		return product, usecase_error.ErrBadParamInput
	}
	return product, nil
}

func (a *AdapterProductJSON) DecodeUpdateInput(input []byte) (adapter.ProductUpdateInput, error) {
	var product adapter.ProductUpdateInput
	if err := json.Unmarshal(input, &product); err != nil {
		fmt.Printf("[JSON-PRODUCT-ADAPTER] : DECODE UPDATE INPUT %#v \n", err)

		return product, usecase_error.ErrBadParamInput
	}
	return product, nil
}
