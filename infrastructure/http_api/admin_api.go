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

type AdminAPI interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	UpdateBiodata(w http.ResponseWriter, r *http.Request)
	UploadPhotoProfile(w http.ResponseWriter, r *http.Request)
	UpdatePassword(w http.ResponseWriter, r *http.Request)
	AddAddress(w http.ResponseWriter, r *http.Request)
	UpdateAddress(w http.ResponseWriter, r *http.Request)
	DeleteAddress(w http.ResponseWriter, r *http.Request)
	DeleteOne(w http.ResponseWriter, r *http.Request)
}

type adminAPI struct {
	adminUsecase logic.AdminUsecase
	serialize    adapter.AdminAdapter
	authUsecase  logic.AuthenticationUsecase
	gStorage     *storage.Client
}

func NewAdminAPI(
	adminUsecase logic.AdminUsecase,
	authUsecase logic.AuthenticationUsecase,
	gStorage *storage.Client,
) AdminAPI {
	return &adminAPI{
		adminUsecase: adminUsecase,
		serialize:    &adapterJSON.AdapterAdminJSON{},
		authUsecase:  authUsecase,
		gStorage:     gStorage,
	}
}

func (a *adminAPI) Create(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	fmt.Printf("ADMIN CREATE : WITH TOKEN %v \n", token)
	credential, err := a.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := a.authUsecase.VerifiedAsAdmin(credential); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := a.serialize.DecodeCreateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	admin, err := a.adminUsecase.Create(r.Context(), input)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, admin)
}

func (a *adminAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	adminID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	fmt.Printf("ADMIN GET BY ID :%v - WITH TOKEN : %v \n", adminID, token)
	credential, err := a.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := a.authUsecase.VerifiedAdminAuthor(credential, adminID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	admin, err := a.adminUsecase.GetByID(r.Context(), adminID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, admin)
}

func (a *adminAPI) Fetch(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	fmt.Printf("ADMIN FETCH :  WITH TOKEN : %v \n", token)
	credential, err := a.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := a.authUsecase.VerifiedAsAdmin(credential); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	search := adapter.AdminSearchOptions{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}

	var defaultNum int64 = 10
	if num, err := strconv.Atoi(r.FormValue("num")); err == nil {
		defaultNum = int64(num)
	}

	cursor := r.FormValue("cursor")
	admins, err := a.adminUsecase.Fetch(r.Context(), cursor, defaultNum, search)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, admins)
}

func (a *adminAPI) UpdateBiodata(w http.ResponseWriter, r *http.Request) {
	adminID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	fmt.Printf("ADMIN UPDATE BIODATA : ID : %v - WITH TOKEN : %v \n", adminID, token)
	credential, err := a.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := a.authUsecase.VerifiedAdminAuthor(credential, adminID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := a.serialize.DecodeUpdateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	admin, err := a.adminUsecase.UpdateBiodata(r.Context(), input, adminID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, admin)
}

func (a *adminAPI) UploadPhotoProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ADMIN UPDATE PHOTO Profile")
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

	adminID := mux.Vars(r)["id"]
	token := r.Header.Get("token")

	credential, err := a.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := a.authUsecase.VerifiedAdminAuthor(credential, adminID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	multipart := helper.NewMultiPart(a.gStorage)
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

	admin, err := a.adminUsecase.UploadAvatar(r.Context(), fileName, adminID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, admin)
}

func (a *adminAPI) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	adminID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	fmt.Printf("ADMIN UPDATE PASSWORD : ID : %v - TOKEN : %v \n", adminID, token)

	credential, err := a.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := a.authUsecase.VerifiedAdminAuthor(credential, adminID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := a.serialize.DecodeUpdatePasswordInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	admin, err := a.adminUsecase.UpdatePassword(r.Context(), input, adminID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, admin)
}

func (a *adminAPI) AddAddress(w http.ResponseWriter, r *http.Request) {
	adminID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	fmt.Printf("ADMIN ADD ADDRESS : ID : %v - TOKEN : %v \n", adminID, token)

	credential, err := a.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := a.authUsecase.VerifiedAdminAuthor(credential, adminID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := a.serialize.DecodeAddressInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := a.adminUsecase.AddAddress(r.Context(), input, adminID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusCreated, customer)
}

func (a *adminAPI) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	addressID := mux.Vars(r)["addID"]
	adminID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	fmt.Printf("ADMIN UPDATE ADDRESS : ID : %v - TOKEN : %v - addID : %v \n", adminID, token, addressID)

	credential, err := a.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := a.authUsecase.VerifiedAdminAuthor(credential, adminID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := a.serialize.DecodeAddressUpdate(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := a.adminUsecase.UpdateAddress(r.Context(), input, addressID, adminID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, customer)
}

func (a *adminAPI) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	addressID := mux.Vars(r)["addID"]
	adminID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	fmt.Printf("ADMIN DELETE ADDRESS : ID : %v - TOKEN : %v - addID : %v \n", adminID, token, addressID)

	credential, err := a.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := a.authUsecase.VerifiedAdminAuthor(credential, adminID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	customer, err := a.adminUsecase.RemoveAddress(r.Context(), addressID, adminID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, customer)
}

func (a *adminAPI) DeleteOne(w http.ResponseWriter, r *http.Request) {

}
