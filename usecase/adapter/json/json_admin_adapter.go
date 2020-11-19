package adapterJSON

import (
	"encoding/json"
	"fmt"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterAdminJSON struct{}

func (a *AdapterAdminJSON) DecodeCreateInput(input []byte) (adapter.AdminCreateInput, error) {
	var admin adapter.AdminCreateInput
	if err := json.Unmarshal(input, &admin); err != nil {
		fmt.Printf("[JSON-ADMIN-ADAPTER] : DECODE CREATE INPUT %#v \n", err)
		return admin, usecase_error.ErrBadParamInput
	}
	return admin, nil
}

func (a *AdapterAdminJSON) DecodeUpdateInput(input []byte) (adapter.AdminUpdateInput, error) {
	var admin adapter.AdminUpdateInput
	if err := json.Unmarshal(input, &admin); err != nil {
		fmt.Printf("[JSON-ADMIN-ADAPTER] : DECODE UPDATE INPUT %#v \n", err)
		return admin, usecase_error.ErrBadParamInput
	}
	return admin, nil
}

func (a *AdapterAdminJSON) DecodeUpdatePasswordInput(input []byte) (adapter.AdminUpdatePasswordInput, error) {
	var password adapter.AdminUpdatePasswordInput
	if err := json.Unmarshal(input, &password); err != nil {
		fmt.Printf("[JSON-ADMIN-ADAPTER] : DECODE UPDATE PASSWORD %#v \n", err)
		return password, usecase_error.ErrBadParamInput
	}
	return password, nil
}

func (a *AdapterAdminJSON) DecodeAddressInput(input []byte) (adapter.AdminAddressCreateInput, error) {
	var address adapter.AdminAddressCreateInput
	if err := json.Unmarshal(input, &address); err != nil {
		fmt.Printf("[JSON-ADMIN-ADAPTER] : DECODE ADDRESS INPUT %#v \n", err)
		return address, usecase_error.ErrBadParamInput
	}
	return address, nil
}

func (a *AdapterAdminJSON) DecodeAddressUpdate(input []byte) (adapter.AdminAddressUpdateInput, error) {
	var address adapter.AdminAddressUpdateInput
	if err := json.Unmarshal(input, &address); err != nil {
		fmt.Printf("[JSON-ADMIN-ADAPTER] : DECODE UPDATE ADDRESS %#v \n", err)
		return address, usecase_error.ErrBadParamInput
	}
	return address, nil
}
