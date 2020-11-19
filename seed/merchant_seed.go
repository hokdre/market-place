package seed

import (
	"context"
	"sync"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/logic"
)

var (
	MerchantMikeID   string
	MerchantHadiID   string
	MerchantJainalID string
)

func SeedMerchant(merchantUsecase logic.MerchantUsecase, customers []domain.Customer) ([]domain.Merchant, error) {
	var wgSeedMerchant sync.WaitGroup
	wgSeedMerchant.Add(2)

	var merchants []domain.Merchant
	var errs []error

	for _, customer := range customers {
		credential := domain.Credential{
			UserID: customer.ID,
		}
		ctx := context.WithValue(context.Background(), "credential", credential)

		switch customer.Email {
		case EmailJainal:
			go func(ctxJainal context.Context) {
				defer wgSeedMerchant.Done()
				merchantJainal, err := seedJainalMerchant(ctxJainal, merchantUsecase)
				if err != nil {
					errs = append(errs, err)
					return
				}

				MerchantJainalID = merchantJainal.ID
				merchants = append(merchants, merchantJainal)
			}(ctx)
		case EmailMike:
			go func(ctxMike context.Context) {
				defer wgSeedMerchant.Done()
				merchantMike, err := seedMikeMerchant(ctxMike, merchantUsecase)
				if err != nil {
					errs = append(errs, err)
					return
				}

				MerchantMikeID = merchantMike.ID
				merchants = append(merchants, merchantMike)
			}(ctx)
		}
	}

	wgSeedMerchant.Wait()
	if len(errs) != 0 {
		return []domain.Merchant{}, errs[0]
	}

	return merchants, nil
}

func seedMikeMerchant(ctx context.Context, merchantUsecase logic.MerchantUsecase) (domain.Merchant, error) {
	merchantMike := adapter.MerchantCreateInput{
		Name: "Merchant Mike",
		Address: adapter.Address{
			City: domain.City{
				CityID:       "79",
				ProvinceID:   "9",
				ProvinceName: "Jawa Barat",
				CityName:     "Bogor",
				PostalCode:   "16110",
			},
			Street: "jl bogor raya",
			Number: "10A",
		},
		Phone:       "089199889988",
		Description: "Merchant ini mejual berbagai macam kebutuhan laptop",
		ShippingID:  "pos",
	}

	mikeMerchant, err := merchantUsecase.Create(ctx, merchantMike)
	if err != nil {
		return domain.Merchant{}, err
	}

	ctxAddEtalase, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	//prepare time space insert and update in elasticsearch
	time.Sleep(30 * time.Millisecond)
	mikeMerchant, err = merchantUsecase.AddEtalase(
		ctxAddEtalase,
		adapter.MerchantEtalaseCreateInput{
			Name: "laptop",
		},
		mikeMerchant.ID,
	)
	if err != nil {
		return domain.Merchant{}, err
	}

	return mikeMerchant, nil
}

func seedJainalMerchant(ctx context.Context, merchantUsecase logic.MerchantUsecase) (domain.Merchant, error) {
	merchantJainal := adapter.MerchantCreateInput{
		Name: "Merchant Jainal",
		Address: adapter.Address{
			City: domain.City{
				CityID:       "223",
				ProvinceID:   "18",
				ProvinceName: "Lampung",
				CityName:     "Lampung Barat",
				PostalCode:   "34814",
			},
			Street: "jl lampung raya",
			Number: "10A",
		},
		Phone:       "089199889988",
		Description: "Merchant ini mejual berbagai macam kebutuhan handphone",
		ShippingID:  "pos",
	}

	jainalMerchant, err := merchantUsecase.Create(ctx, merchantJainal)
	if err != nil {
		return domain.Merchant{}, err
	}

	ctxAddEtalase, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	//prepare time space insert and update in elasticsearch
	time.Sleep(30 * time.Millisecond)
	jainalMerchant, err = merchantUsecase.AddEtalase(
		ctxAddEtalase,
		adapter.MerchantEtalaseCreateInput{
			Name: "handphone",
		},
		jainalMerchant.ID,
	)
	if err != nil {
		return domain.Merchant{}, err
	}

	return jainalMerchant, nil
}
