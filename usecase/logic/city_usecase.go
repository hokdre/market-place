package logic

import (
	"context"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/repository"
)

type CityUsecase interface {
	GetCity(ctx context.Context, keyword string) ([]domain.City, error)
}

type cityUsecase struct {
	cityRepo       repository.CityRepository
	contextTimeout time.Duration
}

func NewCityUsecase(
	cityRepo repository.CityRepository,
	contextTimeout time.Duration,
) CityUsecase {
	return &cityUsecase{
		cityRepo:       cityRepo,
		contextTimeout: contextTimeout,
	}
}

func (c *cityUsecase) GetCity(ctx context.Context, keyword string) ([]domain.City, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return c.cityRepo.GetCity(ctx, keyword)
}
