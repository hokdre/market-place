package logic

import (
	"context"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/repository"
)

type SearchUsecase interface {
	SuggestionSearch(ctx context.Context, keyword string) (domain.Search, error)
	ProductSearch(ctx context.Context, category, secondCategory, thirdCategory string, city string, min, max int64, keyword string, lastDate string) (domain.SearchProduct, error)
	ProductTerlarisSearch(ctx context.Context, page int64)
	ProductTerpopulerSearch(ctx context.Context, page int64)
	MerchantProductSearch(ctx context.Context, merchantId string, etalase string, productName string, lastItemDate string, size int64) ([]domain.Product, error)
}

type searchUsecase struct {
	contextTimeOut time.Duration
	searchRepo     repository.SearchRepository
}

func NewSearchUsecase(
	searchRepo repository.SearchRepository,
	contextTimeOut time.Duration,
) SearchUsecase {
	return &searchUsecase{
		contextTimeOut: contextTimeOut,
		searchRepo:     searchRepo,
	}
}

func (s *searchUsecase) SuggestionSearch(ctx context.Context, keyword string) (domain.Search, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeOut)
	defer cancel()
	return s.searchRepo.SuggestionSearch(ctx, keyword)
}

func (s *searchUsecase) ProductSearch(ctx context.Context, category, secondCategory, thirdCategory string, city string, min, max int64, keyword string, lastDate string) (domain.SearchProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeOut)
	defer cancel()
	return s.searchRepo.ProductSearch(ctx, category, secondCategory, thirdCategory, city, min, max, keyword, lastDate)
}

func (s *searchUsecase) ProductTerlarisSearch(ctx context.Context, page int64) {

}

func (s *searchUsecase) ProductTerpopulerSearch(ctx context.Context, page int64) {

}

func (s *searchUsecase) MerchantProductSearch(ctx context.Context, merchantId string, etalase string, productName string, lastItemDate string, size int64) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeOut)
	defer cancel()

	return s.searchRepo.MerchantProductSearch(ctx, merchantId, etalase, productName, lastItemDate, size)
}
