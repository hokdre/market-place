package http_api

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/market-place/infrastructure/http_api/http_response"
	"github.com/market-place/usecase/adapter"
	adapterJSON "github.com/market-place/usecase/adapter/json"
	"github.com/market-place/usecase/logic"
)

type ShippingAPI interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	UpdateOne(w http.ResponseWriter, r *http.Request)
	DeleteOne(w http.ResponseWriter, r *http.Request)
}

type shippingAPI struct {
	shippingUsecase logic.ShippingUsecase
	authUsecase     logic.AuthenticationUsecase
	serialize       adapter.ShippingAdapter
}

func NewShippingAPI(
	shippingUsecase logic.ShippingUsecase,
	authUsecase logic.AuthenticationUsecase,
) ShippingAPI {
	return &shippingAPI{
		shippingUsecase: shippingUsecase,
		authUsecase:     authUsecase,
		serialize:       &adapterJSON.AdapterShippingJSON{},
	}
}

func (s *shippingAPI) Create(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	credential, err := s.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := s.authUsecase.VerifiedAsAdmin(credential); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := s.serialize.DecodeCreateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	shipping, err := s.shippingUsecase.Create(r.Context(), input)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, shipping)
}

func (s *shippingAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	shippingID := mux.Vars(r)["id"]
	shipping, err := s.shippingUsecase.GetByID(r.Context(), shippingID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, shipping)
}

func (s *shippingAPI) Fetch(w http.ResponseWriter, r *http.Request) {
	search := adapter.ShippingProviderSearchOptions{
		Name: r.FormValue("name"),
	}
	var defaultNum int64 = 10
	if num, err := strconv.Atoi(r.FormValue("num")); err == nil {
		defaultNum = int64(num)
	}

	cursor := r.FormValue("cursor")
	shippings, err := s.shippingUsecase.Fetch(r.Context(), cursor, defaultNum, search)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, shippings)
}

func (s *shippingAPI) UpdateOne(w http.ResponseWriter, r *http.Request) {
	shippingID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := s.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := s.authUsecase.VerifiedAsAdmin(credential); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := s.serialize.DecodeUpdateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	shipping, err := s.shippingUsecase.UpdateOne(r.Context(), input, shippingID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, shipping)
}

func (s *shippingAPI) DeleteOne(w http.ResponseWriter, r *http.Request) {
	shippingID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := s.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := s.authUsecase.VerifiedAsAdmin(credential); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	shipping, err := s.shippingUsecase.DeleteOne(r.Context(), shippingID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, shipping)
}
