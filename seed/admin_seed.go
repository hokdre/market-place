package seed

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/logic"
)

func SeedSuperAdmin(adminUsecase logic.AdminUsecase) (domain.Admin, error) {
	log.SetOutput(os.Stdout)
	log.Println("Seed admin: starting!!")
	birthDay, err := time.Parse(time.RFC3339, "1998-09-11T00:00:00Z")
	if err != nil {
		return domain.Admin{}, err
	}

	input := adapter.AdminCreateInput{
		Email:      "superadmin@gmail.com",
		Name:       "superadmin",
		Password:   "password!23Z",
		RePassword: "password!23Z",
		Addresses: []adapter.Address{
			adapter.Address{
				City: domain.City{
					CityID:       "327",
					CityName:     "Palembang",
					ProvinceID:   "33",
					ProvinceName: "Sumatera Selatan",
					PostalCode:   "30111",
				},
				Street: "jl tanah kusir II ",
				Number: "52",
			},
		},
		Born:     "jakarta",
		BirthDay: birthDay,
		Phone:    "085273989895",
		Gender:   "M",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	admin, err := adminUsecase.Create(ctx, input)
	if err != nil {
		log.Printf("Seed admin: failed cause, %s \n", err)
	}
	return admin, err
}
