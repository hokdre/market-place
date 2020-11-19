package http_response

import (
	"encoding/json"
	"net/http"

	"github.com/market-place/usecase/usecase_error"
)

type ErrResponse struct {
	Message string      `json:"error_message"`
	Errors  interface{} `json:"errors"`
}

type OkResponse struct {
	Data interface{} `json:"data"`
}

func SendErrJSON(w http.ResponseWriter, err error) {
	var errResponse ErrResponse
	errResponse.Message = err.Error()
	//validation entity error contains response feedback field error for client
	if err.Error() == usecase_error.ErrEntityValidation {
		errs, ok := err.(usecase_error.ErrBadEntityInput)
		if ok {
			errResponse.Errors = errs
		}
	}
	//validation login contains response feedback field error for client
	if err.Error() == usecase_error.ErrLogin {
		errs, ok := err.(usecase_error.ErrLoginField)
		if ok {
			errResponse.Errors = errs
		}
	}

	var HTTPCode int = http.StatusInternalServerError
	switch err {
	case usecase_error.ErrInternalServerError:
		HTTPCode = http.StatusInternalServerError
	case usecase_error.ErrBadParamInput:
		HTTPCode = http.StatusBadRequest
	case usecase_error.ErrConflict:
		HTTPCode = http.StatusBadRequest
	case usecase_error.ErrNotFound:
		HTTPCode = http.StatusNotFound
	case usecase_error.ErrNotAuthentication:
		HTTPCode = http.StatusUnauthorized
	case usecase_error.ErrNotAuthorization:
		HTTPCode = http.StatusForbidden
	default:
		if err.Error() == usecase_error.ErrEntityValidation {
			HTTPCode = http.StatusBadRequest
		}
		if err.Error() == usecase_error.ErrLogin {
			HTTPCode = http.StatusUnauthorized
		}
	}

	res, err := json.Marshal(errResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	send(w, HTTPCode, res)
}

func SendOkJSON(w http.ResponseWriter, HTTPCode int, data interface{}) {
	var okResponse OkResponse
	okResponse.Data = data

	res, err := json.Marshal(okResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	send(w, HTTPCode, res)
}

func send(w http.ResponseWriter, HTTPCode int, res []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(HTTPCode)
	if _, err := w.Write(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
