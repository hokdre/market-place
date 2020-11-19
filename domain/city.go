package domain

import "fmt"

type City struct {
	CityID       string `json:"city_id" bson:"city_id" validate:"required"`
	CityName     string `json:"city_name" bson:"city_name" validate:"required"`
	ProvinceID   string `json:"province_id" bson:"province_id" validate:"required"`
	ProvinceName string `json:"province" bson:"province" validate:"required"`
	PostalCode   string `json:"postal_code" bson:"postal_code" validate:"required"`
}

func (c *City) EncodeToString() string {
	return fmt.Sprintf(
		`%s:%s:%s:%s:%s`,
		c.CityID,
		c.CityName,
		c.ProvinceID,
		c.ProvinceName,
		c.PostalCode,
	)
}
