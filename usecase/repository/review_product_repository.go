package repository

import (
	"context"

	"github.com/market-place/domain"
)

type RProductRepository interface {
	Create(ctx context.Context, review domain.RProduct) (domain.RProduct, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.RProductSearchOptions) ([]domain.RProduct, error)
	GetByID(ctx context.Context, id string) (domain.RProduct, error)
	UpdateOne(ctx context.Context, review domain.RProduct) (domain.RProduct, error)
	DeleteOne(ctx context.Context, review domain.RProduct) (domain.RProduct, error)
	DeleteAll(ctx context.Context) error
}
