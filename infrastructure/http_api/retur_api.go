package http_api

import (
	"net/http"

	"github.com/market-place/usecase/logic"
)

type ReturAPI interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	AcceptRetur(w http.ResponseWriter, r *http.Request)
	RejectRetur(w http.ResponseWriter, r *http.Request)
	InputShipping(w http.ResponseWriter, r *http.Request)
}

type returAPI struct {
	returUsecase logic.ReturUseCase
}

func NewReturAPI(returUsecase logic.ReturUseCase) ReturAPI {
	return &returAPI{
		returUsecase: returUsecase,
	}
}

func (re *returAPI) Create(w http.ResponseWriter, r *http.Request) {

}

func (re *returAPI) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (re *returAPI) Fetch(w http.ResponseWriter, r *http.Request) {

}

func (re *returAPI) AcceptRetur(w http.ResponseWriter, r *http.Request) {

}

func (re *returAPI) RejectRetur(w http.ResponseWriter, r *http.Request) {

}

func (re *returAPI) InputShipping(w http.ResponseWriter, r *http.Request) {

}
