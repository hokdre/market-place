package repository

import (
	"context"

	"github.com/market-place/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer domain.Customer) (domain.Customer, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.CustomerSearchOptions) ([]domain.Customer, error)
	GetByID(ctx context.Context, id string) (domain.Customer, error)
	UpdateOne(ctx context.Context, customer domain.Customer) (domain.Customer, error)
	DeleteOne(ctx context.Context, customer domain.Customer) (domain.Customer, error)
	DeleteAll(ctx context.Context) error
}
