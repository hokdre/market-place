package http_api

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/market-place/infrastructure/http_api/http_response"
	adapterJSON "github.com/market-place/usecase/adapter/json"
	"github.com/market-place/usecase/logic"
)

type CartAPI interface {
	GetByID(w http.ResponseWriter, r *http.Request)
	AddProduct(w http.ResponseWriter, r *http.Request)
	UpdateItemInCart(w http.ResponseWriter, r *http.Request)
	RemoveProduct(w http.ResponseWriter, r *http.Request)
	ClearProduct(w http.ResponseWriter, r *http.Request)
}

type cartAPI struct {
	cartUsecase logic.CartUsecase
	authUsecase logic.AuthenticationUsecase
	serialize   adapterJSON.AdapterCartJSON
}

func NewCartAPI(
	cartUsecase logic.CartUsecase,
	authUsecase logic.AuthenticationUsecase,
) CartAPI {
	return &cartAPI{
		cartUsecase: cartUsecase,
		authUsecase: authUsecase,
		serialize:   adapterJSON.AdapterCartJSON{},
	}
}

func (c *cartAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	cartID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCartOwner(credential, cartID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	cart, err := c.cartUsecase.GetByID(r.Context(), cartID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, cart)
}

func (c *cartAPI) AddProduct(w http.ResponseWriter, r *http.Request) {
	cartID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCartOwner(credential, cartID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := c.serialize.DecodeAddItemInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	ctx := context.WithValue(r.Context(), "credential", credential)
	cart, err := c.cartUsecase.AddProduct(ctx, input, cartID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, cart)
}

func (c *cartAPI) UpdateItemInCart(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["pID"]
	cartID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCartOwner(credential, cartID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := c.serialize.DecodeUpdateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	ctx := context.WithValue(r.Context(), "credential", credential)
	cart, err := c.cartUsecase.UpdateItemInCart(ctx, input, productID, cartID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, cart)
}

func (c *cartAPI) RemoveProduct(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["pID"]
	cartID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCartOwner(credential, cartID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	cart, err := c.cartUsecase.RemoveProduct(r.Context(), productID, cartID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, cart)
}

func (c *cartAPI) ClearProduct(w http.ResponseWriter, r *http.Request) {
	cartID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCartOwner(credential, cartID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	cart, err := c.cartUsecase.ClearProduct(r.Context(), cartID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, cart)
}
