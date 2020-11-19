package http_api

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"github.com/market-place/domain"
	"github.com/market-place/infrastructure/http_api/helper"
	"github.com/market-place/infrastructure/http_api/http_response"
	adapterJSON "github.com/market-place/usecase/adapter/json"
	"github.com/market-place/usecase/logic"
	"github.com/market-place/usecase/usecase_error"
)

type TBuyerAPI interface {
	Create(http.ResponseWriter, *http.Request)
	GetByID(http.ResponseWriter, *http.Request)
	Fetch(http.ResponseWriter, *http.Request)
	UploadTransferPhoto(http.ResponseWriter, *http.Request)
	AcceptTransaction(http.ResponseWriter, *http.Request)
	RejectTransaction(http.ResponseWriter, *http.Request)
}

type tbuyerAPI struct {
	tbuyerUsecase logic.TBuyerUsecase
	authUsecase   logic.AuthenticationUsecase
	serialize     adapterJSON.AdapterTBuyerJSON
	gStorage      *storage.Client
}

func NewTBuyerAPI(
	tbuyerUsecase logic.TBuyerUsecase,
	authUsecase logic.AuthenticationUsecase,
	gStorage *storage.Client,
) TBuyerAPI {
	return &tbuyerAPI{
		tbuyerUsecase: tbuyerUsecase,
		authUsecase:   authUsecase,
		serialize:     adapterJSON.AdapterTBuyerJSON{},
		gStorage:      gStorage,
	}
}

func (t *tbuyerAPI) Create(http.ResponseWriter, *http.Request) {

}

func (t *tbuyerAPI) GetByID(http.ResponseWriter, *http.Request) {

}

func (t *tbuyerAPI) Fetch(w http.ResponseWriter, r *http.Request) {
	customerID := r.FormValue("customerID")
	status := r.FormValue("status")
	adminID := r.FormValue("adminID")

	token := r.Header.Get("token")
	credential, err := t.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	if customerID != "" && credential.LoginType == domain.LOGIN_AS_CUSTOMER {
		if err := t.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
			http_response.SendErrJSON(w, err)
			return
		}
	} else if customerID == "" {
		if err := t.authUsecase.VerifiedAsAdmin(credential); err != nil {
			http_response.SendErrJSON(w, err)
			return
		}
	}

	search := domain.TBuyerSearchOptions{
		CustomerID: customerID,
		AdminID:    adminID,
		Status:     status,
	}

	var defaultNum int64 = 10
	if num, err := strconv.Atoi(r.FormValue("num")); err == nil {
		defaultNum = int64(num)
	}

	cursor := r.FormValue("cursor")
	orders, err := t.tbuyerUsecase.Fetch(r.Context(), cursor, defaultNum, search)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, orders)
}

func (t *tbuyerAPI) UploadTransferPhoto(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "TransferPhoto",
				Message: "TransferPhoto is too large",
			},
		}
		http_response.SendErrJSON(w, err)
		return
	}

	transID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := t.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := t.authUsecase.VerifiedCustomerTransactionOwner(r.Context(), credential, transID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	multipart := helper.NewMultiPart(t.gStorage)
	err = multipart.ReadAvatar(r)
	if err != nil {
		err = usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "TransferPhoto",
				Message: "TransferPhoto is empty",
			},
		}
		http_response.SendErrJSON(w, err)
		return
	}
	if isEmage, err := multipart.IsImageAvatar(); !isEmage || err != nil {
		if err != nil {
			http_response.SendErrJSON(w, err)
			return
		}
		if !isEmage {
			err := usecase_error.ErrBadEntityInput{
				usecase_error.ErrEntityField{
					Field:   "TransferPhoto",
					Message: "TransferPhoto is not correct file type",
				},
			}
			http_response.SendErrJSON(w, err)
			return
		}
	}
	fileName, err := multipart.StorePhoto(r.Context())
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	tbuyer, err := t.tbuyerUsecase.UploadTransferPhoto(r.Context(), fileName, transID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, tbuyer)
}

func (t *tbuyerAPI) AcceptTransaction(w http.ResponseWriter, r *http.Request) {
	transID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := t.authUsecase.ValidateLogin(token)
	ctx := context.WithValue(r.Context(), "credential", credential)

	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := t.authUsecase.VerifiedAsAdmin(credential); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	tbuyer, err := t.tbuyerUsecase.AcceptTransaction(ctx, transID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, tbuyer)
}

func (t *tbuyerAPI) RejectTransaction(w http.ResponseWriter, r *http.Request) {
	transID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := t.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	ctx := context.WithValue(r.Context(), "credential", credential)
	if err := t.authUsecase.VerifiedAsAdmin(credential); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	input, err := t.serialize.DecodeRejectInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	tbuyer, err := t.tbuyerUsecase.RejectTransaction(ctx, input, transID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, tbuyer)
}
