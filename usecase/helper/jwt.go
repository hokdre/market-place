package helper

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/market-place/domain"
	"github.com/market-place/usecase/usecase_error"
)

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte(os.Getenv("JWT_SECRET"))

func EncodeToken(credential domain.Credential) (string, error) {
	var token string
	unSingnedToken := jwt.NewWithClaims(JWT_SIGNING_METHOD, credential)
	token, err := unSingnedToken.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		return token, usecase_error.ErrInternalServerError
	}
	return token, nil
}

func DecodeToken(token string) (domain.Credential, error) {
	var credential domain.Credential

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method != JWT_SIGNING_METHOD {
			err := errors.New("Token method is not match")
			return credential, err
		}

		return JWT_SIGNATURE_KEY, nil
	})

	if err != nil {
		return credential, usecase_error.ErrNotAuthentication
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return credential, usecase_error.ErrNotAuthentication
	}

	credential.Email = claims["email"].(string)
	credential.Issuer = claims["iss"].(string)
	credential.UserID = claims["user_id"].(string)
	credential.CartID = claims["cart_id"].(string)
	credential.LoginType = claims["login_type"].(string)

	return credential, nil
}
