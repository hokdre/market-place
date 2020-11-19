package domain

type Category struct {
	Top       string `json:"top" bson:"top" validate:"required"`
	SecondSub string `json:"second_sub" bson:"second_sub" validate:"required"`
	ThirdSub  string `json:"third_sub" bson:"third_sub" validate:"required"`
}
