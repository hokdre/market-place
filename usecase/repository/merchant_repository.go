package repository

import (
	"context"

	"github.com/market-place/domain"
)

type MerchantRepository interface {
	Create(ctx context.Context, merchant domain.Merchant) (domain.Merchant, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.MerchantSearchOptions) ([]domain.Merchant, error)
	GetByID(ctx context.Context, id string) (domain.Merchant, error)
	//name is unique
	GetByName(ctx context.Context, name string) (domain.Merchant, error)
	UpdateOne(ctx context.Context, merchant domain.Merchant) (domain.Merchant, error)
	DeleteOne(ctx context.Context, merchant domain.Merchant) (domain.Merchant, error)
	DeleteAll(ctx context.Context) error
}
