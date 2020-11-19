package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var APPLICATION_NAME = "Market-Place"
var LOGIN_AS_CUSTOMER = "CUSTOMER"
var LOGIN_AS_ADMIN = "ADMIN"

type Credential struct {
	jwt.StandardClaims
	UserID     string `json:"user_id"`
	CartID     string `json:"cart_id"`
	Email      string `json:"email"`
	MerchantID string `json:"merchant_id"`
	LoginType  string `json:"login_type"`
}

func NewCredential(userID, cartID, merchantId, email, loginType string) Credential {
	return Credential{
		jwt.StandardClaims{
			Issuer:   APPLICATION_NAME,
			IssuedAt: time.Now().Unix(),
		},
		userID,
		cartID,
		email,
		merchantId,
		loginType,
	}
}
