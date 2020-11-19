package seed

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/logic"
)

func SeedShipping(shippingUsecase logic.ShippingUsecase) ([]domain.ShippingProvider, error) {
	log.SetOutput(os.Stdout)
	log.Println("Seed Shipping : starting!")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var shippings []domain.ShippingProvider
	var errs []error
	var wgSeedShipping sync.WaitGroup
	wgSeedShipping.Add(3)

	go func() {
		defer wgSeedShipping.Done()
		jne := adapter.ShippingCreateInput{
			ID:   "jne",
			Name: "JNE",
		}
		shipping, err := shippingUsecase.Create(ctx, jne)
		if err != nil {
			log.Printf("Seed Shipping : failed cause, %s \n", err)
			errs = append(errs, err)
			return
		}

		shippings = append(shippings, shipping)
	}()

	go func() {
		defer wgSeedShipping.Done()
		jne := adapter.ShippingCreateInput{
			ID:   "tiki",
			Name: "TIKI",
		}
		shipping, err := shippingUsecase.Create(ctx, jne)
		if err != nil {
			log.Printf("Seed Shipping : failed cause, %s \n", err)
			errs = append(errs, err)
			return
		}

		shippings = append(shippings, shipping)
	}()

	go func() {
		defer wgSeedShipping.Done()
		jne := adapter.ShippingCreateInput{
			ID:   "pos",
			Name: "Pos Indonesia",
		}
		shipping, err := shippingUsecase.Create(ctx, jne)
		if err != nil {
			log.Printf("Seed Shipping : failed cause, %s \n", err)
			errs = append(errs, err)
		}
		shippings = append(shippings, shipping)
	}()

	wgSeedShipping.Wait()
	if len(errs) != 0 {
		return shippings, errs[0]
	}

	log.Println("Seed Shipping : success!")
	return shippings, nil
}
