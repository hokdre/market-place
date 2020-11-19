package repository

import (
	"context"

	"github.com/market-place/domain"
)

type TBuyerRepository interface {
	Create(ctx context.Context, transaction domain.TBuyer) (domain.TBuyer, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.TBuyerSearchOptions) ([]domain.TBuyer, error)
	GetByID(ctx context.Context, id string) (domain.TBuyer, error)
	UpdateOne(ctx context.Context, transaction domain.TBuyer) (domain.TBuyer, error)
	DeleteOne(ctx context.Context, transaction domain.TBuyer) (domain.TBuyer, error)
	DeleteAll(ctx context.Context) error
}
