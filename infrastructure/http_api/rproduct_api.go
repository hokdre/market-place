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

type RProductAPI interface {
	Create(http.ResponseWriter, *http.Request)
	Fetch(http.ResponseWriter, *http.Request)
}

type rproductAPI struct {
	rproductUsecase logic.ReviewProductUsecase
	authUsecase     logic.AuthenticationUsecase
	serialize       adapterJSON.AdapterRProductJSON
}

func NewRProductAPI(
	rproductUsecase logic.ReviewProductUsecase,
	authUsecase logic.AuthenticationUsecase,
) RProductAPI {
	return &rproductAPI{
		rproductUsecase: rproductUsecase,
		authUsecase:     authUsecase,
		serialize:       adapterJSON.AdapterRProductJSON{},
	}
}

func (rp *rproductAPI) Create(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	credential, err := rp.authUsecase.ValidateLogin(token)
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
	input, err := rp.serialize.DecodeCreateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := rp.authUsecase.VerifiedCustomerOrderOwner(ctx, credential, input.OrderID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	review, err := rp.rproductUsecase.Create(ctx, input)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, review)
}

func (rp *rproductAPI) Fetch(w http.ResponseWriter, r *http.Request) {
	productID := r.FormValue("merchantID")
	last := r.FormValue("last")
	search := domain.RProductSearchOptions{
		ProductID: productID,
		Last:      last,
	}

	reviews, err := rp.rproductUsecase.Fetch(r.Context(), search)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, reviews)
}
