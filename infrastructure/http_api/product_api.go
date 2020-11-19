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

type ProductAPI interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	UpdateData(w http.ResponseWriter, r *http.Request)
	UploadPhotos(w http.ResponseWriter, r *http.Request)
	DeleteOne(w http.ResponseWriter, r *http.Request)
	ProductTerlaris(w http.ResponseWriter, r *http.Request)
}

type productAPI struct {
	productUsecase logic.ProductUsecase
	authUsecase    logic.AuthenticationUsecase
	serialize      adapterJSON.AdapterProductJSON
	gStorage       *storage.Client
}

func NewProductAPI(
	productUsecase logic.ProductUsecase,
	authUsecase logic.AuthenticationUsecase,
	gStorage *storage.Client,
) ProductAPI {
	return &productAPI{
		productUsecase: productUsecase,
		authUsecase:    authUsecase,
		serialize:      adapterJSON.AdapterProductJSON{},
		gStorage:       gStorage,
	}
}

func (p *productAPI) Create(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	credential, err := p.authUsecase.ValidateLogin(token)
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
	input, err := p.serialize.DecodeCreateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	product, err := p.productUsecase.Create(ctx, input)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusCreated, product)
}

func (p *productAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]
	product, err := p.productUsecase.GetByID(r.Context(), productID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, product)
}

func (p *productAPI) Fetch(w http.ResponseWriter, r *http.Request) {
	price := int64(0)
	if qPrice, err := strconv.Atoi(r.FormValue("price")); err == nil {
		price = int64(qPrice)
	}

	search := adapter.ProductSearchOptions{
		Name:        r.FormValue("name"),
		Category:    r.FormValue("category"),
		Description: r.FormValue("description"),
		MerchantID:  r.FormValue("merchant_id"),
		Etalase:     r.FormValue("etalase"),
		Price:       price,
		City:        r.FormValue("city"),
	}

	var defaultNum int64 = 10
	if num, err := strconv.Atoi(r.FormValue("num")); err == nil {
		defaultNum = int64(num)
	}

	cursor := r.FormValue("cursor")
	products, err := p.productUsecase.Fetch(r.Context(), cursor, defaultNum, search)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, products)
}

func (p *productAPI) UpdateData(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := p.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := p.authUsecase.VerifiedProductOwner(r.Context(), credential, productID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	input, err := p.serialize.DecodeUpdateInput(requestBody)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	product, err := p.productUsecase.UpdateData(r.Context(), input, productID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, product)
}

func (p *productAPI) UploadPhotos(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		err := usecase_error.ErrBadEntityInput{
			usecase_error.ErrEntityField{
				Field:   "Photos",
				Message: "Photos is too large",
			},
		}
		http_response.SendErrJSON(w, err)
		return
	}

	productID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := p.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := p.authUsecase.VerifiedProductOwner(r.Context(), credential, productID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	multipart := helper.NewMultiPart(p.gStorage)
	err = multipart.ReadPhotos(r)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if isImage, err := multipart.IsImagePhotos(); !isImage || err != nil {
		if !isImage {
			err = usecase_error.ErrBadEntityInput{
				usecase_error.ErrEntityField{
					Field:   "Photos",
					Message: "Photos is not correct file type",
				},
			}
		}

		http_response.SendErrJSON(w, err)
		return
	}
	fileNames, err := multipart.StorePhotos(r.Context())
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	product, err := p.productUsecase.UploadPhotos(r.Context(), fileNames, productID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, product)
}

func (p *productAPI) DeleteOne(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]
	token := r.Header.Get("token")
	credential, err := p.authUsecase.ValidateLogin(token)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	if err := p.authUsecase.VerifiedProductOwner(r.Context(), credential, productID); err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	product, err := p.productUsecase.DeleteOne(r.Context(), productID)
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}

	http_response.SendOkJSON(w, http.StatusOK, product)
}

func (p *productAPI) ProductTerlaris(w http.ResponseWriter, r *http.Request) {
	products, err := p.productUsecase.ProductTerlaris(r.Context())
	if err != nil {
		http_response.SendErrJSON(w, err)
		return
	}
	http_response.SendOkJSON(w, http.StatusOK, products)
}
