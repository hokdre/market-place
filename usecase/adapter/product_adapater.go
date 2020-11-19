package adapter

import "github.com/market-place/domain"

type ProductUpdateInput struct {
	Name        string          `json:"name"`
	Category    domain.Category `json:"category"`
	Tags        []string        `json:"tags"`
	Etalase     string          `json:"etalase"`
	Colors      []string        `json:"colors"`
	Sizes       []string        `json:"sizes"`
	Weight      float64         `json:"weight"`
	Width       float64         `json:"width"`
	Height      float64         `json:"height"`
	Long        float64         `json:"long"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Stock       float64         `json:"stock"`
}

type ProductCreateInput struct {
	Name        string          `json:"name"`
	Category    domain.Category `json:"category"`
	Etalase     string          `json:"etalase"`
	Tags        []string        `json:"tags"`
	Colors      []string        `json:"colors"`
	Sizes       []string        `json:"sizes"`
	Weight      float64         `json:"weight"`
	Width       float64         `json:"width"`
	Height      float64         `json:"height"`
	Long        float64         `json:"long"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Stock       float64         `json:"stock"`
}

type ProductSearchOptions struct {
	//product's name contains regex search name keyword
	Name string
	//product's categories have item with category name contains regex search category keyword
	Category   string
	MerchantID string
	Etalase    string
	Tags       string
	//product's description contains regex search keyword
	Description string
	//product's price in range between search keyword price
	Price int64
	//product's merchant city equals to search city keyword
	City string
}

type ProductAdapter interface {
	DecodeCreateInput([]byte) (ProductCreateInput, error)
	DecodeUpdateInput([]byte) (ProductUpdateInput, error)
}
