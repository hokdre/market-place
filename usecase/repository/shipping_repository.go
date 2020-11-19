package repository

import (
	"context"

	"github.com/market-place/domain"
)

type ShippingRepository interface {
	Create(ctx context.Context, transaction domain.ShippingProvider) (domain.ShippingProvider, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.ShippingProviderSearchOptions) ([]domain.ShippingProvider, error)
	GetByID(ctx context.Context, id string) (domain.ShippingProvider, error)
	UpdateOne(ctx context.Context, transaction domain.ShippingProvider) (domain.ShippingProvider, error)
	DeleteOne(ctx context.Context, transaction domain.ShippingProvider) (domain.ShippingProvider, error)
	DeleteAll(ctx context.Context) error
}
