package repository

import (
	"context"

	"github.com/market-place/domain"
)

type TRefundRepository interface {
	Create(ctx context.Context, transaction domain.TRefund) (domain.TRefund, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.TRefundSearchOptions) ([]domain.TRefund, error)
	GetByID(ctx context.Context, id string) (domain.TRefund, error)
	UpdateOne(ctx context.Context, transaction domain.TRefund) (domain.TRefund, error)
	DeleteOne(ctx context.Context, transaction domain.TRefund) (domain.TRefund, error)
	DeleteAll(ctx context.Context) error
}
