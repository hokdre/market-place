package repository

import (
	"context"

	"github.com/market-place/domain"
)

type SearchRepository interface {
	SuggestionSearch(ctx context.Context, keyword string) (domain.Search, error)
	ProductSearch(ctx context.Context, category, secondCategory, thirdCategory string, city string, min, max int64, keyword string, lastDate string) (domain.SearchProduct, error)
	ProductTerlarisSearch(ctx context.Context, page int64)
	ProductTerpopulerSearch(ctx context.Context, page int64)
	MerchantProductSearch(ctx context.Context, merchantId string, etalase string, productName string, lastDate string, number int64) ([]domain.Product, error)
}
