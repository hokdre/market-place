package http_api

import (
	"net/http"

	"github.com/market-place/usecase/logic"
)

type TSellerAPI interface {
	Create(http.ResponseWriter, *http.Request)
	GetByID(http.ResponseWriter, *http.Request)
	Fetch(http.ResponseWriter, *http.Request)
	UpdateOne(http.ResponseWriter, *http.Request)
}

type tsellerAPI struct {
	tsellerUsecase logic.TSellerUsecase
}

func NewTSellerAPI(tsellerUsecase logic.TSellerUsecase) TSellerAPI {
	return &tsellerAPI{
		tsellerUsecase: tsellerUsecase,
	}
}

func (t *tsellerAPI) Create(http.ResponseWriter, *http.Request) {

}

func (t *tsellerAPI) GetByID(http.ResponseWriter, *http.Request) {

}

func (t *tsellerAPI) Fetch(http.ResponseWriter, *http.Request) {

}

func (t *tsellerAPI) UpdateOne(http.ResponseWriter, *http.Request) {

}
