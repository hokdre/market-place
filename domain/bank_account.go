package domain

type BankAccount struct {
	ID       string `json:"_id" bson:"_id" validate:"required"`
	Number   string `json:"number" bson:"number" validate:"required,bank_number"`
	BankCode string `json:"bank_code" bson:"bank_code" validate:"required,bank_provider"`
}
