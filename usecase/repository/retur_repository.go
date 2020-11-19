package repository

import (
	"context"

	"github.com/market-place/domain"
)

type ReturRepository interface {
	Create(ctx context.Context, transaction domain.Retur) (domain.Retur, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.ReturSearchOptions) ([]domain.Retur, error)
	GetByID(ctx context.Context, id string) (domain.Retur, error)
	UpdateOne(ctx context.Context, transaction domain.Retur) (domain.Retur, error)
	DeleteOne(ctx context.Context, transaction domain.Retur) (domain.Retur, error)
	DeleteAll(ctx context.Context) error
}
