package repository

import (
	"context"

	"github.com/market-place/domain"
)

type CityRepository interface {
	GetCity(ctx context.Context, keyword string) ([]domain.City, error)
}
