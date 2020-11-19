package usecase_error

import (
	"errors"
)

var (
	// ErrConfilt will throw if user can do a procedure just one but try more than one
	//example : user try create merchant twice
	ErrConflict = errors.New("Item Is Already Exist")
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Item Is Not Exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Param Is Not Valid")
	// ErrNotAuthorization will throw if the given credential data user from decoded JWT token is not valid
	ErrNotAuthorization = errors.New("Credential Is Not Valid")
	// ErrNotAuthentication will throw if then given jwt token is not found in request header or token cannot be decoded
	ErrNotAuthentication = errors.New("Token Is Not Valid")
	// ErrBadEntityInput will throw if the given input data is not pass in entity validation and pass ErrEntityField for client response
	ErrEntityValidation = "ENTITY VALIDATION"
	// ErrLogin will throw if the given username and password is not match or not find, and pass ErrLoginField for client response
	ErrLogin = "LOGIN FAILED"
)

//response error ErrLogin for client
type ErrLoginField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ErrLoginField) Error() string {
	return ErrLogin
}

//response err  ErrEntityValidation for client
type ErrEntityField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrBadEntityInput []ErrEntityField

func (e ErrBadEntityInput) Error() string {
	return ErrEntityValidation
}
