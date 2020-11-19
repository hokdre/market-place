package http_api

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"github.com/market-place/infrastructure/http_api/helper"
	"github.com/market-place/infrastructure/http_api/http_response"
	"github.com/market-place/usecase/adapter"
	adapterJSON "github.com/market-place/usecase/adapter/json"
	"github.com/market-place/usecase/logic"
	"github.com/market-place/usecase/usecase_error"
)

type MerchantAPI interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	UpdateData(w http.ResponseWriter, r *http.Request)
	UploadAvatar(w http.ResponseWriter, r *http.Request)
	AddShipping(w http.ResponseWriter, r *http.Request)
	RemoveShipping(w http.ResponseWriter, r *http.Request)
	AddBankAccounts(w http.ResponseWriter, r *http.Request)
	AddEtalase(w http.ResponseWriter, r *http.Request)
	DeleteEtalase(w http.ResponseWriter, r *http.Request)
	UpdateBankAccount(w http.ResponseWriter, r *http.Request)
}

type merchantAPI struct {
	merchantUsecase logic.MerchantUsecase
	authUsecase     logic.AuthenticationUsecase
	serialize       adapterJSON.AdapterMerchantJSON
	gStorage        *storage.Client
}

func NewMerchantAPI(
	merchantUsecase logic.MerchantUsecase,
	authUsecase logic.AuthenticationUsecase,
	gStorage *storage.Client,
) MerchantAPI {
	return &merchantAPI{
		merchantUsecase: merchantUsecase,
		authUsecase:     authUsecase,
		serialize:       adapterJSON.AdapterMerchantJSON{},
		gStorage:        gStorage,
	}
}

func (m *merchantAPI) Create(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	credential, err := m.authUsecase.ValidateLogin(token)
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
	input, err := m.serialize.DecodeCreateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	merchant, err := m.merchantUsecase.Create(ctx, input)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, merchant)
}

func (m *merchantAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	merchantID := mux.Vars(r)["id"]
	merchant, err := m.merchantUsecase.GetByID(r.Context(), merchantID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, merchant)
}

func (m *merchantAPI) Fetch(w http.ResponseWriter, r *http.Request) {
	search := adapter.MerchantSearchOptions{
		Name:        r.FormValue("name"),
		City:        r.FormValue("city"),
		Description: r.FormValue("description"),
	}

	var defaultNum int64 = 10
	if num, err := strconv.Atoi(r.FormValue("num")); err == nil {
		defaultNum = int64(num)
	}

	cursor := r.FormValue("cursor")
	merchants, err := m.merchantUsecase.Fetch(r.Context(), cursor, defaultNum, search)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, merchants)
}

func (m *merchantAPI) UpdateData(w http.ResponseWriter, r *http.Request) {
	merchantID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := m.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := m.authUsecase.VerifiedMerchantOwner(credential, merchantID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := m.serialize.DecodeUpdateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	merchant, err := m.merchantUsecase.UpdateData(r.Context(), input, merchantID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, merchant)
}

func (m *merchantAPI) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "Avatar",
				Message: "Avatar is too large",
			},
		}
		http_response.SendErrJSON(w, err)
		return
	}

	merchantID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := m.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := m.authUsecase.VerifiedMerchantOwner(credential, merchantID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	multipart := helper.NewMultiPart(m.gStorage)
	err = multipart.ReadAvatar(r)
	if err != nil {
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
					Field:   "Avatar",
					Message: "Avatar is not correct file type",
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

	merchant, err := m.merchantUsecase.UploadAvatar(r.Context(), fileName, merchantID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, merchant)
}

func (m *merchantAPI) AddShipping(w http.ResponseWriter, r *http.Request) {
	shippingID := mux.Vars(r)["sID"]
	merchantID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := m.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := m.authUsecase.VerifiedMerchantOwner(credential, merchantID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	shipping, err := m.merchantUsecase.AddShipping(r.Context(), shippingID, merchantID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusCreated, shipping)
}

func (m *merchantAPI) RemoveShipping(w http.ResponseWriter, r *http.Request) {
	shippingID := mux.Vars(r)["sID"]
	merchantID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := m.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := m.authUsecase.VerifiedMerchantOwner(credential, merchantID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	shipping, err := m.merchantUsecase.RemoveShipping(r.Context(), shippingID, merchantID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, shipping)
}

func (m *merchantAPI) AddBankAccounts(w http.ResponseWriter, r *http.Request) {
	merchantID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := m.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := m.authUsecase.VerifiedMerchantOwner(credential, merchantID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := m.serialize.DecodeBankInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	merchant, err := m.merchantUsecase.AddBankAccount(r.Context(), input, merchantID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, merchant)
}

func (m *merchantAPI) AddEtalase(w http.ResponseWriter, r *http.Request) {
	merchantID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := m.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := m.authUsecase.VerifiedMerchantOwner(credential, merchantID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := m.serialize.DecodeMerchantEtalaseInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	merchant, err := m.merchantUsecase.AddEtalase(r.Context(), input, merchantID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, merchant)
}

func (m *merchantAPI) DeleteEtalase(w http.ResponseWriter, r *http.Request) {
	merchantID := mux.Vars(r)["id"]
	etalase := mux.Vars(r)["etalase"]
	token := r.Header.Get("token")
	credential, err := m.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := m.authUsecase.VerifiedMerchantOwner(credential, merchantID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	merchant, err := m.merchantUsecase.DeleteEtalase(r.Context(), etalase, merchantID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, merchant)
}

func (m *merchantAPI) UpdateBankAccount(w http.ResponseWriter, r *http.Request) {
	merchantID := mux.Vars(r)["id"]
	accountBankID := mux.Vars(r)["bankID"]
	token := r.Header.Get("token")
	credential, err := m.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := m.authUsecase.VerifiedMerchantOwner(credential, merchantID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := m.serialize.DecodeBankUpdate(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	merchant, err := m.merchantUsecase.UpdateBankAccount(r.Context(), input, accountBankID, merchantID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, merchant)
}
