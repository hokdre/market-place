package repository

import (
	"context"

	"github.com/market-place/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) (domain.Product, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.ProductSearchOptions) ([]domain.Product, error)
	GetByID(ctx context.Context, id string) (domain.Product, error)
	UpdateOne(ctx context.Context, product domain.Product) (domain.Product, error)
	DeleteOne(ctx context.Context, product domain.Product) (domain.Product, error)
	DeleteAll(ctx context.Context) error
}
