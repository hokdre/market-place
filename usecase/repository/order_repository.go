package repository

import (
	"context"

	"github.com/market-place/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order domain.Order) (domain.Order, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.OrderSearchOptions) ([]domain.Order, error)
	GetByID(ctx context.Context, id string) (domain.Order, error)
	UpdateOne(ctx context.Context, order domain.Order) (domain.Order, error)
	DeleteOne(ctx context.Context, order domain.Order) (domain.Order, error)
	DeleteAll(ctx context.Context) error
	EstimasiPendapatan(ctx context.Context, merchantID string, startDay string, endDay string) ([]map[string]interface{}, error)
	OrderSummary(ctx context.Context, merchantID string, startDay string, endDay string) (map[string]int64, error)
	ProductTerlaris(ctx context.Context) ([]map[string]interface{}, error)
}
