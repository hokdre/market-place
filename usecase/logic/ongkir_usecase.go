package logic

import (
	"context"
	"math"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/repository"
)

type OngkirUsecase interface {
	GetOngkir(ctx context.Context, origin, destination string, weight float64, providers []string) ([]domain.Ongkir, error)
}

type ongkirUsecase struct {
	ongkirRepo     repository.OngkirRepository
	contextTimeout time.Duration
}

func NewOngkirUsecase(
	ongkirRepo repository.OngkirRepository,
	contextTimeout time.Duration,
) OngkirUsecase {
	return &ongkirUsecase{
		ongkirRepo:     ongkirRepo,
		contextTimeout: contextTimeout,
	}
}

func (o *ongkirUsecase) GetOngkir(ctx context.Context, origin, destination string, weight float64, providers []string) ([]domain.Ongkir, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	ongkirs, err := o.ongkirRepo.GetOngkir(ctx, origin, destination, providers)
	if err != nil {
		return ongkirs, err
	}

	minimumWeight := int64(1)
	overWeightTolerance := float64(300)

	overWeight := math.Mod(weight, 1000)
	kg := int64(weight-overWeight) / 1000

	isWeightLessAKg := kg < 1
	isOverWeightRoundUp := overWeight > overWeightTolerance

	if isWeightLessAKg {
		kg = minimumWeight
	}
	if !isWeightLessAKg && isOverWeightRoundUp {
		kg++
	}

	for ongkirIndex, ongkir := range ongkirs {
		for serviceIndex, service := range ongkir.Services {
			service.Cost = kg * service.Cost
			ongkirs[ongkirIndex].Services[serviceIndex] = service
		}
	}

	return ongkirs, nil
}
