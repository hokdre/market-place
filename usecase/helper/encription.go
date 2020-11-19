package helper

import (
	"log"
	"os"

	"github.com/market-place/usecase/usecase_error"
	"golang.org/x/crypto/bcrypt"
)

type Encription interface {
	Encrypt(pass []byte) (string, error)
	Compare(hashedPassword, password []byte) bool
}

func NewEncription() Encription {
	return &bcryptEncription{
		cost: bcrypt.MinCost,
	}
}

type bcryptEncription struct {
	cost int
}

func (b *bcryptEncription) Encrypt(pass []byte) (string, error) {
	log.SetOutput(os.Stdout)

	var passwordEncrypted string
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		log.Printf("Encryp password : %s \n", err)
		return passwordEncrypted, usecase_error.ErrInternalServerError
	}

	passwordEncrypted = string(hash)
	return passwordEncrypted, nil
}

func (b *bcryptEncription) Compare(hashedPassword, password []byte) bool {
	log.SetOutput(os.Stdout)
	var passwordCorrect bool = false
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		log.Printf("Decrypt password : %s \n", err)
		return passwordCorrect
	}

	passwordCorrect = true
	return passwordCorrect
}
