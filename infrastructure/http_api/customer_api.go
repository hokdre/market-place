package http_api

import (
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

type CustomerAPI interface {
	Create(http.ResponseWriter, *http.Request)
	Fetch(http.ResponseWriter, *http.Request)
	GetByID(http.ResponseWriter, *http.Request)
	UpdateBiodata(http.ResponseWriter, *http.Request)
	UpdatePassword(http.ResponseWriter, *http.Request)
	UploadPhotoProfile(http.ResponseWriter, *http.Request)
	AddAddress(http.ResponseWriter, *http.Request)
	UpdateAddress(http.ResponseWriter, *http.Request)
	DeleteAddress(http.ResponseWriter, *http.Request)
	AddBankAccount(http.ResponseWriter, *http.Request)
	UpdateBankaccount(http.ResponseWriter, *http.Request)
	DeleteOne(http.ResponseWriter, *http.Request)
}

type customerAPI struct {
	customerUsecase logic.CustomerUsecase
	authUsecase     logic.AuthenticationUsecase
	serialize       adapterJSON.AdapaterCustomerJSON
	gStorage        *storage.Client
}

func NewCustomerAPI(
	customerUsecase logic.CustomerUsecase,
	authUsecase logic.AuthenticationUsecase,
	gStorage *storage.Client,
) CustomerAPI {
	return &customerAPI{
		customerUsecase: customerUsecase,
		authUsecase:     authUsecase,
		serialize:       adapterJSON.AdapaterCustomerJSON{},
		gStorage:        gStorage,
	}
}

func (c *customerAPI) Create(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Printf("[CUSTOMER API] : READ ALL BODY : %#v \n", err)
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := c.serialize.DecodeCreateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := c.customerUsecase.Create(r.Context(), input)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, customer)
}

func (c *customerAPI) Fetch(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedAsAdmin(credential); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	search := adapter.CustomerSearchOptions{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}

	var defaultNum int64 = 10
	if num, err := strconv.Atoi(r.FormValue("num")); err == nil {
		defaultNum = int64(num)
	}

	cursor := r.FormValue("cursor")
	customers, err := c.customerUsecase.Fetch(r.Context(), cursor, defaultNum, search)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, customers)
}

func (c *customerAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
		if err := c.authUsecase.VerifiedAsAdmin(credential); err != nil {
			http_response.SendErrJSON(w, err)
			return
		}
	}

	customer, err := c.customerUsecase.GetByID(r.Context(), customerID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, customer)
}

func (c *customerAPI) UpdateBiodata(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := c.serialize.DecodeUpdateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := c.customerUsecase.UpdateBiodata(r.Context(), input, customerID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, customer)
}

func (c *customerAPI) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := c.serialize.DecodeUpdatePasswordInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := c.customerUsecase.UpdatePassword(r.Context(), input, customerID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, customer)
}

func (c *customerAPI) UploadPhotoProfile(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "Avatar",
				Message: "Avatar is too large",
			},
		}
		fmt.Printf("[CUSTOMER API] : UPDATE PROFILE PHOTO > FILE LARGE : %#v \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	customerID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		fmt.Printf("[CUSTOMER API] : UPDATE PROFILE PHOTO > VALIDATE LOGIN : %#v \n", err)
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
		fmt.Printf("[CUSTOMER API] : UPDATE PROFILE PHOTO > VALIDATE AUTHOR : %#v \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	multipart := helper.NewMultiPart(c.gStorage)
	err = multipart.ReadAvatar(r)
	if err != nil {
		fmt.Printf("[CUSTOMER API] : UPDATE PROFILE PHOTO > READ MULTIPART : %#v \n", err)
		http_response.SendErrJSON(w, err)
		return
	}
	if isEmage, err := multipart.IsImageAvatar(); !isEmage || err != nil {
		if err != nil {
			fmt.Printf("[CUSTOMER API] : UPDATE PROFILE PHOTO > READ EXTENTION FILE : %#v \n", err)
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
			fmt.Printf("[CUSTOMER API] : UPDATE PROFILE PHOTO > VALIDATE IMAGE EXTENTION : %#v \n", err)
			http_response.SendErrJSON(w, err)
			return
		}
	}
	fileName, err := multipart.StorePhoto(r.Context())
	if err != nil {
		fmt.Printf("[CUSTOMER API] : UPDATE PROFILE PHOTO > STORE IMAGE : %#v \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := c.customerUsecase.UploadAvatar(r.Context(), fileName, customerID)
	if err != nil {
		fmt.Printf("[CUSTOMER API] : UPDATE PROFILE PHOTO > USECASE RETURN : %#v \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, customer)
}

func (c *customerAPI) AddBankAccount(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := c.serialize.DecodeBankInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := c.customerUsecase.AddBankAccount(r.Context(), input, customerID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, customer)
}

func (c *customerAPI) UpdateBankaccount(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	accountBankID := mux.Vars(r)["bankID"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		fmt.Printf("[CUSTOMER API] : UPDATE PROFILE PHOTO > VALIDATE IMAGE EXTENTION : %#v \n", err)
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := c.serialize.DecodeBankUpdate(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := c.customerUsecase.UpdateBankAccount(r.Context(), input, accountBankID, customerID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, customer)
}

func (c *customerAPI) AddAddress(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := c.serialize.DecodeAddressInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := c.customerUsecase.AddAddress(r.Context(), input, customerID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusCreated, customer)
}

func (c *customerAPI) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	addressID := mux.Vars(r)["addID"]
	customerID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := c.serialize.DecodeAddressUpdate(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := c.customerUsecase.UpdateAddress(r.Context(), input, addressID, customerID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, customer)
}

func (c *customerAPI) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	addressID := mux.Vars(r)["addID"]
	customerID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := c.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := c.authUsecase.VerifiedCustomerAuthor(credential, customerID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := c.customerUsecase.RemoveAddress(r.Context(), addressID, customerID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, customer)
}

func (c *customerAPI) DeleteOne(w http.ResponseWriter, r *http.Request) {

}
