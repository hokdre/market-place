package http_api

import (
	"net/http"

	"github.com/market-place/usecase/logic"
)

type TRefundAPI interface {
	Create(http.ResponseWriter, *http.Request)
	GetByID(http.ResponseWriter, *http.Request)
	Fetch(http.ResponseWriter, *http.Request)
	UpdateOne(http.ResponseWriter, *http.Request)
}

type trefundAPI struct {
	trefundUsecase logic.TRefundUsecase
}

func NewTRefundAPI(trefundUsecase logic.TRefundUsecase) TRefundAPI {
	return &trefundAPI{
		trefundUsecase: trefundUsecase,
	}
}

func (t *trefundAPI) Create(http.ResponseWriter, *http.Request) {

}

func (t *trefundAPI) GetByID(http.ResponseWriter, *http.Request) {

}

func (t *trefundAPI) Fetch(http.ResponseWriter, *http.Request) {

}

func (t *trefundAPI) UpdateOne(http.ResponseWriter, *http.Request) {

}
