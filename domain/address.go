package domain

type Address struct {
	ID     string `json:"_id" bson:"_id" validate:"required"`
	City   City `json:"city" bson:"city" validate:"required"`
	Street string `json:"street" bson:"street" validate:"required"`
	Number string `json:"number" bson:"number" validate:"required"`
}
