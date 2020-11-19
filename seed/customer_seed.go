package seed

import (
	"context"
	"sync"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/logic"
)

const (
	EmailHadi   string = "hadinw@gmail.com"
	EmailMike   string = "mukminmike@gmail.com"
	EmailJainal string = "jainal@gmail.com"
)

func SeedCustomer(customerUsecase logic.CustomerUsecase) ([]domain.Customer, error) {
	hadiBirthday, _ := time.Parse(time.RFC3339, "1997-01-11T00:00:00Z")
	customer1 := adapter.CustomerCreateInput{
		Email:      EmailHadi,
		Name:       "hadi nur wahid",
		Password:   "password!23Z",
		RePassword: "password!23Z",
		Addresses: []adapter.Address{
			adapter.Address{
				City: domain.City{
					CityID:       "255",
					ProvinceID:   "11",
					ProvinceName: "Jawa Timur",
					CityName:     "Malang",
					PostalCode:   "65163",
				},
				Street: "jl malang raya",
				Number: "10A",
			},
		},
		Gender:   "M",
		Born:     "Malang",
		BirthDay: hadiBirthday,
		Phone:    "099898989898",
	}

	mikeBirthDay, _ := time.Parse(time.RFC3339, "1997-01-11T00:00:00Z")
	customer2 := adapter.CustomerCreateInput{
		Email:      EmailMike,
		Name:       "Mukmin Mike",
		Password:   "password!23Z",
		RePassword: "password!23Z",
		Addresses: []adapter.Address{
			adapter.Address{
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
		},
		Gender:   "M",
		Born:     "Bogor",
		BirthDay: mikeBirthDay,
		Phone:    "099898989898",
	}

	jainalBirthDay, _ := time.Parse(time.RFC3339, "1997-01-11T00:00:00Z")
	customer3 := adapter.CustomerCreateInput{
		Email:      EmailJainal,
		Name:       "Jainal",
		Password:   "password!23Z",
		RePassword: "password!23Z",
		Addresses: []adapter.Address{
			adapter.Address{
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
		},
		Gender:   "M",
		Born:     "Lampung",
		BirthDay: jainalBirthDay,
		Phone:    "099898989898",
	}

	ctxSeedCustomer, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wgSeedCustomer sync.WaitGroup
	wgSeedCustomer.Add(3)

	var customers []domain.Customer
	var errs []error
	go func() {
		defer wgSeedCustomer.Done()

		hadi, err := customerUsecase.Create(ctxSeedCustomer, customer1)
		if err != nil {
			errs = append(errs, err)
		}

		customers = append(customers, hadi)
	}()

	go func() {
		defer wgSeedCustomer.Done()

		mike, err := customerUsecase.Create(ctxSeedCustomer, customer2)
		if err != nil {
			errs = append(errs, err)
		}

		customers = append(customers, mike)
	}()

	go func() {
		defer wgSeedCustomer.Done()

		jainal, err := customerUsecase.Create(ctxSeedCustomer, customer3)
		if err != nil {
			errs = append(errs, err)
		}

		customers = append(customers, jainal)
	}()

	wgSeedCustomer.Wait()
	if len(errs) != 0 {
		return customers, errs[0]
	}

	return customers, nil
}
