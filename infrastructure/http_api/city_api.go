package http_api

import (
	"net/http"

	"github.com/market-place/infrastructure/http_api/http_response"
	"github.com/market-place/usecase/logic"
)

type CityAPI interface {
	GetCity(w http.ResponseWriter, r *http.Request)
}

type cityAPI struct {
	cityUsecase logic.CityUsecase
}

func NewCityAPI(cityUsecase logic.CityUsecase) CityAPI {
	return &cityAPI{
		cityUsecase: cityUsecase,
	}
}

func (c *cityAPI) GetCity(w http.ResponseWriter, r *http.Request) {
	keyword := r.FormValue("keyword")
	cities, err := c.cityUsecase.GetCity(r.Context(), keyword)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, cities)
}
