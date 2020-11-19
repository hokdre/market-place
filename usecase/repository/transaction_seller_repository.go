package repository

import (
	"context"

	"github.com/market-place/domain"
)

type TSellerRepository interface {
	Create(ctx context.Context, transaction domain.TSeller) (domain.TSeller, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.TSellerSearchOptions) ([]domain.TSeller, error)
	GetByID(ctx context.Context, id string) (domain.TSeller, error)
	UpdateOne(ctx context.Context, transaction domain.TSeller) (domain.TSeller, error)
	DeleteOne(ctx context.Context, transaction domain.TSeller) (domain.TSeller, error)
	DeleteAll(ctx context.Context) error
}
