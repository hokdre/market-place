package repository

import (
	"context"

	"github.com/market-place/domain"
)

type RMerchantRepository interface {
	Create(ctx context.Context, review domain.RMerchant) (domain.RMerchant, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.RMerchantSearchOptions) ([]domain.RMerchant, error)
	GetByID(ctx context.Context, id string) (domain.RMerchant, error)
	UpdateOne(ctx context.Context, review domain.RMerchant) (domain.RMerchant, error)
	DeleteOne(ctx context.Context, review domain.RMerchant) (domain.RMerchant, error)
	DeleteAll(ctx context.Context) error
}
