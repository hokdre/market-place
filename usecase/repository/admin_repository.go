package repository

import (
	"context"

	"github.com/market-place/domain"
)

type AdminRepository interface {
	Create(ctx context.Context, admin domain.Admin) (domain.Admin, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.AdminSearchOptions) ([]domain.Admin, error)
	GetByID(ctx context.Context, id string) (domain.Admin, error)
	UpdateOne(ctx context.Context, admin domain.Admin) (domain.Admin, error)
	DeleteOne(ctx context.Context, admin domain.Admin) (domain.Admin, error)
	DeleteAll(ctx context.Context) error
}
