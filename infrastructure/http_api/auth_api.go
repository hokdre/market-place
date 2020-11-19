package http_api

import (
	"io/ioutil"
	"net/http"

	"github.com/market-place/infrastructure/http_api/http_response"
	adapterJSON "github.com/market-place/usecase/adapter/json"
	"github.com/market-place/usecase/helper"
	"github.com/market-place/usecase/logic"
)

type AuthAPI interface {
	CustomerLogin(w http.ResponseWriter, r *http.Request)
	AdminLogin(w http.ResponseWriter, r *http.Request)
}

type authAPI struct {
	authUsecase logic.AuthenticationUsecase
	serialize   adapterJSON.AdapterAuthJSON
}

func NewAuthAPI(authUsecase logic.AuthenticationUsecase) AuthAPI {
	return &authAPI{
		authUsecase: authUsecase,
		serialize:   adapterJSON.AdapterAuthJSON{},
	}
}

func (a *authAPI) CustomerLogin(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	input, err := a.serialize.DecodeLoginInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	credential, err := a.authUsecase.LoginCustomer(r.Context(), input)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	token, err := helper.EncodeToken(credential)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	res := map[string]interface{}{
		"token":      token,
		"credential": credential,
	}
	http_response.SendOkJSON(w, http.StatusOK, res)
}
func (a *authAPI) AdminLogin(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	input, err := a.serialize.DecodeLoginInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	credential, err := a.authUsecase.LoginAdmin(r.Context(), input)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	token, err := helper.EncodeToken(credential)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	res := map[string]interface{}{
		"token":      token,
		"credential": credential,
	}
	http_response.SendOkJSON(w, http.StatusOK, res)
}
