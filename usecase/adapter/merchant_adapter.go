package adapter

import "github.com/market-place/domain"

type MerchantCreateInput struct {
	Name          string               `json:"name"`
	Address       Address              `json:"address"`
	Phone         string               `json:"phone"`
	LocationPoint domain.LocationPoint `json:"location_point"`
	Description   string               `json:"description"`
	ShippingID    string               `json:"shipping_id"`
}

type MerchantUpdateInput struct {
	Phone         string               `json:"phone"`
	Description   string               `json:"description"`
	Address       Address              `json:"address"`
	LocationPoint domain.LocationPoint `json:"location_point"`
}

type MerchantBankCreateInput struct {
	Number   string `json:"number"`
	BankCode string `json:"bank_code"`
}

type MerchantBankUpdateInput struct {
	Number   string `json:"number"`
	BankCode string `json:"bank_code"`
}

type MerchantEtalaseCreateInput struct {
	Name string `json:"name"`
}

type MerchantSearchOptions struct {
	//Merchant's name contains regex search keyword
	Name string
	//Merchant's address city equals to with search keyword
	City string
	//Merchant's description regex search keyword
	Description string
}

type MerchantAdapter interface {
	DecodeCreateInput([]byte) (MerchantCreateInput, error)
	DecodeUpdateInput([]byte) (MerchantUpdateInput, error)
	DecodeBankInput([]byte) (MerchantBankCreateInput, error)
	DecodeBankUpdate([]byte) (MerchantBankUpdateInput, error)
}
