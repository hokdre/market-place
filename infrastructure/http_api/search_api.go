package http_api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/market-place/infrastructure/http_api/http_response"
	"github.com/market-place/usecase/logic"
)

type SearchAPI interface {
	SuggestionSearch(w http.ResponseWriter, r *http.Request)
	ProductSearch(w http.ResponseWriter, r *http.Request)
	ProductTerlarisSearch(w http.ResponseWriter, r *http.Request)
	ProductTerpopulerSearch(w http.ResponseWriter, r *http.Request)
	MerchantProductSearch(w http.ResponseWriter, r *http.Request)
}

type searchAPI struct {
	searchUsecase logic.SearchUsecase
}

func NewSearchAPI(
	searchUsecase logic.SearchUsecase,
) SearchAPI {
	return &searchAPI{
		searchUsecase: searchUsecase,
	}
}

func (s *searchAPI) SuggestionSearch(w http.ResponseWriter, r *http.Request) {
	keyword := r.FormValue("keyword")
	suggestions, err := s.searchUsecase.SuggestionSearch(r.Context(), keyword)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, suggestions)
}

func (s *searchAPI) ProductSearch(w http.ResponseWriter, r *http.Request) {
	topCategory := r.FormValue("c")
	secondCategory := r.FormValue("sc")
	thirdCategory := r.FormValue("tc")
	city := r.FormValue("city")
	min := 0
	if r.FormValue("min") != "" {
		min, _ = strconv.Atoi(r.FormValue("min"))
	}
	max := 0
	if r.FormValue("max") != "" {
		max, _ = strconv.Atoi(r.FormValue("max"))
	}
	keyword := r.FormValue("keyword")
	lastDate := r.FormValue("last")

	products, err := s.searchUsecase.ProductSearch(r.Context(), topCategory, secondCategory, thirdCategory, city, int64(min), int64(max), keyword, lastDate)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, products)
}

func (s *searchAPI) ProductTerlarisSearch(w http.ResponseWriter, r *http.Request) {

}

func (s *searchAPI) ProductTerpopulerSearch(w http.ResponseWriter, r *http.Request) {

}

func (s *searchAPI) MerchantProductSearch(w http.ResponseWriter, r *http.Request) {
	merchantId := mux.Vars(r)["id"]

	etalase := r.FormValue("etalase")
	sizeStr := r.FormValue("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		size = 0
	}
	lastItemDate := r.FormValue("last")
	productName := r.FormValue("name")

	products, err := s.searchUsecase.MerchantProductSearch(r.Context(), merchantId, etalase, productName, lastItemDate, int64(size))
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, products)
}
