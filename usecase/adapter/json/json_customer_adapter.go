package adapterJSON

import (
	"encoding/json"
	"fmt"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapaterCustomerJSON struct{}

func (a *AdapaterCustomerJSON) DecodeCreateInput(input []byte) (adapter.CustomerCreateInput, error) {
	var customer adapter.CustomerCreateInput
	if err := json.Unmarshal(input, &customer); err != nil {
		fmt.Printf("[CUSTOMER JSON] : DECODE CREATE INPUT : %#v \n", err)
		return customer, usecase_error.ErrBadParamInput
	}
	return customer, nil
}

func (a *AdapaterCustomerJSON) DecodeUpdateInput(input []byte) (adapter.CustomerUpdateInput, error) {
	var customer adapter.CustomerUpdateInput
	if err := json.Unmarshal(input, &customer); err != nil {
		return customer, usecase_error.ErrBadParamInput
	}
	return customer, nil
}

func (a *AdapaterCustomerJSON) DecodeUpdatePasswordInput(input []byte) (adapter.CustomerUpdatePasswordInput, error) {
	var password adapter.CustomerUpdatePasswordInput
	if err := json.Unmarshal(input, &password); err != nil {
		return password, usecase_error.ErrBadParamInput
	}
	return password, nil
}

func (a *AdapaterCustomerJSON) DecodeBankInput(input []byte) (adapter.CustomerBankCreateInput, error) {
	var bankAccount adapter.CustomerBankCreateInput
	if err := json.Unmarshal(input, &bankAccount); err != nil {
		return bankAccount, usecase_error.ErrBadParamInput
	}
	return bankAccount, nil
}

func (a *AdapaterCustomerJSON) DecodeBankUpdate(input []byte) (adapter.CustomerBankUpdateInput, error) {
	var bankAccount adapter.CustomerBankUpdateInput
	if err := json.Unmarshal(input, &bankAccount); err != nil {
		return bankAccount, usecase_error.ErrBadParamInput
	}
	return bankAccount, nil
}

func (a *AdapaterCustomerJSON) DecodeAddressInput(input []byte) (adapter.CustomerAddressCreateInput, error) {
	var address adapter.CustomerAddressCreateInput
	if err := json.Unmarshal(input, &address); err != nil {
		return address, usecase_error.ErrBadParamInput
	}
	return address, nil
}

func (a *AdapaterCustomerJSON) DecodeAddressUpdate(input []byte) (adapter.CustomerAddressUpdateInput, error) {
	var address adapter.CustomerAddressUpdateInput
	if err := json.Unmarshal(input, &address); err != nil {
		return address, usecase_error.ErrBadParamInput
	}
	return address, nil
}
