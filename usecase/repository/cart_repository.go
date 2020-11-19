package repository

import (
	"context"

	"github.com/market-place/domain"
)

type CartRepository interface {
	Create(ctx context.Context, cart domain.Cart) (domain.Cart, error)
	Fetch(ctx context.Context, cursor string, num int64, options domain.CartSearchOptions) ([]domain.Cart, error)
	GetByID(ctx context.Context, id string) (domain.Cart, error)
	UpdateOne(ctx context.Context, cart domain.Cart) (domain.Cart, error)
	DeleteOne(ctx context.Context, cart domain.Cart) (domain.Cart, error)
	DeleteAll(ctx context.Context) error
}
