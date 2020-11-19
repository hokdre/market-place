package http_api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/market-place/infrastructure/http_api/http_response"
	"github.com/market-place/usecase/logic"
	"github.com/market-place/usecase/usecase_error"
)

type OngkirAPI interface {
	GetOgnkir(w http.ResponseWriter, r *http.Request)
}

type ongkirAPI struct {
	ongkirUsecase logic.OngkirUsecase
}

func NewOngkirAPI(
	ongkirUsecase logic.OngkirUsecase,
) OngkirAPI {
	return &ongkirAPI{
		ongkirUsecase: ongkirUsecase,
	}
}

func (o *ongkirAPI) GetOgnkir(w http.ResponseWriter, r *http.Request) {
	origin := r.FormValue("origin")
	destination := r.FormValue("destination")
	strProviders := r.FormValue("providers")
	providers := strings.Split(strProviders, ",")

	var weight float64 = 0
	strWeight := r.FormValue("weight")
	weight, err := strconv.ParseFloat(strWeight, 64)
	if err != nil {
		err = usecase_error.ErrBadParamInput
		http_response.SendErrJSON(w, err)
		return
	}

	ongkirs, err := o.ongkirUsecase.GetOngkir(r.Context(), origin, destination, weight, providers)

	http_response.SendOkJSON(w, http.StatusOK, ongkirs)
}
