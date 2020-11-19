package domain

import "time"

type LocationPoint struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lon float64 `json:"lon" bson:"lon"`
}

type Merchant struct {
	ID            string             `json:"_id" bson:"_id" validate:"required"`
	Name          string             `json:"name" bson:"name" validate:"required"`
	Address       Address            `json:"address" bson:"address" validate:"required"`
	Avatar        string             `json:"avatar" bson:"avatar" validate:"required"`
	Phone         string             `json:"phone" bson:"phone"  validate:"required,phone"`
	Description   string             `json:"description" bson:"description" validate:"required,max=250"`
	Etalase       []string           `json:"etalase" bson:"etalase" validate:"unique_etalase"`
	Products      []Product          `json:"products" bson:"-"`
	Reviews       []RMerchant        `json:"reviews" bson:"-"`
	Rating        float64            `json:"rating"  bson:"rating" validate:"max=5"`
	NumReview     float64            `json:"num_review" bson:"num_review"`
	Shippings     []ShippingProvider `json:"shippings" bson:"shippings" validate:"min=1,unique_shippings"`
	BankAccounts  []BankAccount      `json:"bank_accounts" bson:"bank_accounts" validate:"unique_bank_accounts"`
	LocationPoint LocationPoint      `json:"location_point" bson:"location_point"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at" validate:"required"`
}

func (m *Merchant) DenomarlizationData() DenormalizationMerchant {
	denom := DenormalizationMerchant{}
	denom.ID = m.ID
	denom.Name = m.Name
	denom.Avatar = m.Avatar
	denom.Phone = m.Phone
	denom.Address = m.Address
	denom.Shippings = m.Shippings
	denom.LocationPoint = m.LocationPoint
	denom.Rating = m.Rating
	denom.NumReview = m.NumReview
	return denom
}

type DenormalizationMerchant struct {
	ID            string             `json:"_id" bson:"_id" validate:"required"`
	Name          string             `json:"name" bson:"name" validate:"required"`
	Avatar        string             `json:"avatar" bson:"avatar" validate:"required"`
	Phone         string             `json:"phone" bson:"phone"  validate:"required,phone"`
	Address       Address            `json:"address" bson:"address" validate:"required,min=1"`
	Shippings     []ShippingProvider `json:"shippings" bson:"shippings"`
	LocationPoint LocationPoint      `json:"location_point" bson:"location_point"`
	Rating        float64            `json:"rating" bson:"rating"`
	NumReview     float64            `json:"num_review" bson:"num_review"`
}

type MerchantSearchOptions struct {
	//Merchant's name contains regex search keyword
	Name string
	//Merchant's address city equals to with search keyword
	City string
	//Merchant's description regex search keyword
	Description string
	//Merchant's shippings have item's id equals to search shippingID
	ShippingID string
	//Merchant's products have item's id equals to search productID
	ProductID string
	//Merchant's reviews have item's id equals to search reviewID
	ReviewID string
}
