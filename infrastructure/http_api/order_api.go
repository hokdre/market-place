package http_api

import (
	"context"
	"fmt"
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

type OrderAPI interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	FetchOrderMerchant(w http.ResponseWriter, r *http.Request)
	FetchOrderCustomer(w http.ResponseWriter, r *http.Request)
	AddResiNumber(w http.ResponseWriter, r *http.Request)
	UploadShippingPhoto(w http.ResponseWriter, r *http.Request)
	FinishOrder(w http.ResponseWriter, r *http.Request)
	RejectOrder(w http.ResponseWriter, r *http.Request)
	AjukanPaketSampai(w http.ResponseWriter, r *http.Request)
	EstimasiPendapatan(w http.ResponseWriter, r *http.Request)
	OrderSummary(w http.ResponseWriter, r *http.Request)
}

type orderAPI struct {
	orderUsecase logic.OrderUsecase
	authUsecase  logic.AuthenticationUsecase
	serialize    adapterJSON.AdapterOrderJSON
	gStorage     *storage.Client
}

func NewOrderAPI(
	orderUsecase logic.OrderUsecase,
	authUsecase logic.AuthenticationUsecase,
	gStorage *storage.Client,
) OrderAPI {
	return &orderAPI{
		orderUsecase: orderUsecase,
		authUsecase:  authUsecase,
		serialize:    adapterJSON.AdapterOrderJSON{},
		gStorage:     gStorage,
	}
}

func (o *orderAPI) Create(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	credential, err := o.authUsecase.ValidateLogin(token)
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
	input, err := o.serialize.DecodeCreateInput(requestBody)
	if err != nil {
		fmt.Printf("[DEBUG] : %#v \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	orders, err := o.orderUsecase.Create(ctx, input)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, orders)
}

func (o *orderAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	orderId := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	_, err := o.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	order, err := o.orderUsecase.GetByID(r.Context(), orderId)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, order)
}

func (o *orderAPI) FetchOrderMerchant(w http.ResponseWriter, r *http.Request) {
	merchantID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := o.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := o.authUsecase.VerifiedMerchantOwner(credential, merchantID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	search := adapter.OrderSearchOptions{
		MerchantID: credential.MerchantID,
		Status:     r.FormValue("status"),
	}
	var defaultNum int64 = 10
	if num, err := strconv.Atoi(r.FormValue("num")); err == nil {
		defaultNum = int64(num)
	}

	cursor := r.FormValue("cursor")
	orders, err := o.orderUsecase.Fetch(r.Context(), cursor, defaultNum, search)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, orders)
}

func (o *orderAPI) FetchOrderCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := o.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := o.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	search := adapter.OrderSearchOptions{
		CustomerID: customerID,
		Status:     r.FormValue("status"),
	}

	var defaultNum int64 = 10
	if num, err := strconv.Atoi(r.FormValue("num")); err == nil {
		defaultNum = int64(num)
	}

	cursor := r.FormValue("cursor")
	orders, err := o.orderUsecase.Fetch(r.Context(), cursor, defaultNum, search)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, orders)
}

func (o *orderAPI) AddResiNumber(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := o.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := o.authUsecase.VerifiedMerchantOrderOwner(r.Context(), credential, orderID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := o.serialize.DecodeResiInput(body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	order, err := o.orderUsecase.InputResiNumber(r.Context(), input, orderID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, order)
}

func (o *orderAPI) UploadShippingPhoto(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := o.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := o.authUsecase.VerifiedMerchantOrderOwner(r.Context(), credential, orderID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "ShippingPhoto",
				Message: "ShippingPhoto is too large",
			},
		}
		http_response.SendErrJSON(w, err)
		return
	}

	multipart := helper.NewMultiPart(o.gStorage)
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
					Field:   "ShippingPhoto",
					Message: "ShippingPhoto is not correct file type",
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

	order, err := o.orderUsecase.UploadShippingPhoto(r.Context(), fileName, orderID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, order)
}

func (o *orderAPI) FinishOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := o.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := o.authUsecase.VerifiedCustomerOrderOwner(r.Context(), credential, orderID); err != nil {
		fmt.Printf("[DEBUG] : %#v \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	order, err := o.orderUsecase.FinishOrder(r.Context(), orderID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, order)
}

func (o *orderAPI) RejectOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := o.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := o.authUsecase.VerifiedMerchantOrderOwner(r.Context(), credential, orderID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	order, err := o.orderUsecase.RejectOrder(r.Context(), orderID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, order)
}

func (o *orderAPI) AjukanPaketSampai(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := o.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := o.authUsecase.VerifiedMerchantOrderOwner(r.Context(), credential, orderID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	order, err := o.orderUsecase.AjukanPaketSampai(r.Context(), orderID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, order)
}

func (o *orderAPI) EstimasiPendapatan(w http.ResponseWriter, r *http.Request) {
	startDate := r.FormValue("start")
	endDate := r.FormValue("end")
	token := r.Header.Get("token")
	credential, err := o.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	ctx := context.WithValue(r.Context(), "credential", credential)
	summary, err := o.orderUsecase.EstimasiPendapatan(ctx, startDate, endDate)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, summary)
}

func (o *orderAPI) OrderSummary(w http.ResponseWriter, r *http.Request) {
	startDate := r.FormValue("start")
	endDate := r.FormValue("end")
	token := r.Header.Get("token")
	credential, err := o.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	ctx := context.WithValue(r.Context(), "credential", credential)
	summary, err := o.orderUsecase.OrderSummary(ctx, startDate, endDate)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, summary)
}
