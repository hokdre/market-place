package adapterJSON

import (
	"encoding/json"
	"fmt"

	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/usecase_error"
)

type AdapterMerchantJSON struct{}

func (a *AdapterMerchantJSON) DecodeCreateInput(input []byte) (adapter.MerchantCreateInput, error) {
	var merchant adapter.MerchantCreateInput
	if err := json.Unmarshal(input, &merchant); err != nil {
		fmt.Printf("[JSON-MERCHANT-ADAPTER] : DECODE CREATE  : %#v \n", err)
		return merchant, usecase_error.ErrBadParamInput
	}
	return merchant, nil
}

func (a *AdapterMerchantJSON) DecodeUpdateInput(input []byte) (adapter.MerchantUpdateInput, error) {
	var merchant adapter.MerchantUpdateInput
	if err := json.Unmarshal(input, &merchant); err != nil {
		fmt.Printf("[JSON-MERCHANT-ADAPTER] : DECODE UPDATE  : %#v \n", err)
		return merchant, usecase_error.ErrBadParamInput
	}
	return merchant, nil
}

func (a *AdapterMerchantJSON) DecodeBankInput(input []byte) (adapter.MerchantBankCreateInput, error) {
	var bankAccount adapter.MerchantBankCreateInput
	if err := json.Unmarshal(input, &bankAccount); err != nil {
		fmt.Printf("[JSON-MERCHANT-ADAPTER] : DECODE BANK INPUT  : %#v \n", err)
		return bankAccount, usecase_error.ErrBadParamInput
	}
	return bankAccount, nil
}

func (a *AdapterMerchantJSON) DecodeBankUpdate(input []byte) (adapter.MerchantBankUpdateInput, error) {
	var bankAccount adapter.MerchantBankUpdateInput
	if err := json.Unmarshal(input, &bankAccount); err != nil {
		fmt.Printf("[JSON-MERCHANT-ADAPTER] : DECODE BANK UPDATE  : %#v \n", err)
		return bankAccount, usecase_error.ErrBadParamInput
	}
	return bankAccount, nil
}

func (a *AdapterMerchantJSON) DecodeMerchantEtalaseInput(input []byte) (adapter.MerchantEtalaseCreateInput, error) {
	var etalase adapter.MerchantEtalaseCreateInput
	if err := json.Unmarshal(input, &etalase); err != nil {
		fmt.Printf("[JSON-MERCHANT-ADAPTER] : DECODE BANK UPDATE  : %#v \n", err)
		return etalase, usecase_error.ErrBadParamInput
	}
	return etalase, nil
}
