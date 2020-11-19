package adapter

import "github.com/market-place/domain"

type Address struct {
	City   domain.City `json:"city"`
	Street string      `json:"street"`
	Number string      `json:"number"`
}
