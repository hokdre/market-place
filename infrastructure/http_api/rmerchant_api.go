package http_api

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/market-place/domain"
	"github.com/market-place/infrastructure/http_api/http_response"
	adapterJSON "github.com/market-place/usecase/adapter/json"
	"github.com/market-place/usecase/logic"
)

type RMerchantAPI interface {
	Create(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
}

type rmerchantAPI struct {
	rmerchantUsecase logic.ReviewMerchantUsecase
	authUsecase      logic.AuthenticationUsecase
	serialize        adapterJSON.AdapterRMerchantJSON
}

func NewRMerchantAPI(
	rmerchantUsecase logic.ReviewMerchantUsecase,
	authUsecase logic.AuthenticationUsecase,
) RMerchantAPI {
	return &rmerchantAPI{
		rmerchantUsecase: rmerchantUsecase,
		authUsecase:      authUsecase,
		serialize:        adapterJSON.AdapterRMerchantJSON{},
	}
}

func (rm *rmerchantAPI) Create(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	credential, err := rm.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	ctx := context.WithValue(r.Context(), "credential", credential)
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := rm.serialize.DecodeCreateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := rm.authUsecase.VerifiedCustomerOrderOwner(ctx, credential, input.OrderID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	review, err := rm.rmerchantUsecase.Create(ctx, input)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, review)
}

func (rm *rmerchantAPI) Fetch(w http.ResponseWriter, r *http.Request) {
	merchantID := r.FormValue("merchantID")
	last := r.FormValue("last")
	search := domain.RMerchantSearchOptions{
		MerchantID: merchantID,
		Last:       last,
	}
	reviews, err := rm.rmerchantUsecase.Fetch(r.Context(), search)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, reviews)
}
