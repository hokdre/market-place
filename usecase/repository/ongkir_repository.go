package repository

import (
	"context"

	"github.com/market-place/domain"
)

type OngkirRepository interface {
	GetOngkir(ctx context.Context, origin, destination string, providers []string) ([]domain.Ongkir, error)
}
