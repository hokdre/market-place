package adapterJSON

import (
	"encoding/json"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterAuthJSON struct{}

func (a *AdapterAuthJSON) DecodeLoginInput(input []byte) (adapter.LoginInput, error) {
	var login adapter.LoginInput
	if err := json.Unmarshal(input, &login); err != nil {
		return login, usecase_error.ErrBadParamInput
	}
	return login, nil
}
