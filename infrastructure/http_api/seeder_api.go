package http_api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/market-place/infrastructure/http_api/http_response"
	"github.com/market-place/seed"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/logic"
)

type SeederAPI interface {
	SeedSuperAdmin(w http.ResponseWriter, r *http.Request)
	SeedShipping(w http.ResponseWriter, r *http.Request)
	SeedCustomer(w http.ResponseWriter, r *http.Request)
	SeedMerchant(w http.ResponseWriter, r *http.Request)
	SeedProduct(w http.ResponseWriter, r *http.Request)
}

type seederAPI struct {
	adminUsecase     logic.AdminUsecase
	shippingUsecase  logic.ShippingUsecase
	customerUsecase  logic.CustomerUsecase
	merchantUsecase  logic.MerchantUsecase
	productUsecase   logic.ProductUsecase
	cartUsecase      logic.CartUsecase
	orderUsecase     logic.OrderUsecase
	tbuyerUsecase    logic.TBuyerUsecase
	rMerchantUsecase logic.ReviewMerchantUsecase
	rProductUsecase  logic.ReviewProductUsecase
}

func NewSeederAPI(
	adminUsecase logic.AdminUsecase,
	shippingUsecase logic.ShippingUsecase,
	customerUsecase logic.CustomerUsecase,
	merchantUsecase logic.MerchantUsecase,
	productUsecase logic.ProductUsecase,
	cartUsecase logic.CartUsecase,
	orderUsecase logic.OrderUsecase,
	tbuyerUsecase logic.TBuyerUsecase,
	rMerchantUsecase logic.ReviewMerchantUsecase,
	rProductUsecase logic.ReviewProductUsecase,
) SeederAPI {
	return &seederAPI{
		adminUsecase:     adminUsecase,
		shippingUsecase:  shippingUsecase,
		customerUsecase:  customerUsecase,
		merchantUsecase:  merchantUsecase,
		productUsecase:   productUsecase,
		cartUsecase:      cartUsecase,
		orderUsecase:     orderUsecase,
		tbuyerUsecase:    tbuyerUsecase,
		rMerchantUsecase: rMerchantUsecase,
		rProductUsecase:  rProductUsecase,
	}
}

func (s *seederAPI) SeedSuperAdmin(w http.ResponseWriter, r *http.Request) {
	log.SetOutput(os.Stdout)
	log.Println("Seed Admin : starting!")
	admin, err := seed.SeedSuperAdmin(s.adminUsecase)
	if err != nil {
		log.Printf("Seed Admin : failed cause, %s \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	log.Println("Seed Admin: success!")
	http_response.SendOkJSON(w, http.StatusCreated, admin)
}

func (s *seederAPI) SeedShipping(w http.ResponseWriter, r *http.Request) {
	log.SetOutput(os.Stdout)
	log.Println("Seed Shipping : starting!")

	shippings, err := seed.SeedShipping(s.shippingUsecase)
	if err != nil {
		log.Printf("Seed Shipping : failed cause, %s \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	log.Println("Seed Shipping: success!")
	http_response.SendOkJSON(w, http.StatusCreated, shippings)
}

func (s *seederAPI) SeedCustomer(w http.ResponseWriter, r *http.Request) {

	log.SetOutput(os.Stdout)
	log.Println("Seed Customer : starting!")

	customers, err := seed.SeedCustomer(s.customerUsecase)
	if err != nil {
		log.Printf("Seed Customer : failed cause, %s \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	log.Println("Seed Customer: success!")
	http_response.SendOkJSON(w, http.StatusCreated, customers)
}

func (s *seederAPI) SeedMerchant(w http.ResponseWriter, r *http.Request) {
	ctxFetch, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.SetOutput(os.Stdout)
	log.Println("Seed Merchant : starting!")

	customers, err := s.customerUsecase.Fetch(ctxFetch, "", 100, adapter.CustomerSearchOptions{})
	if err != nil {
		log.Printf("Seed Merchant : failed, cause %s \n ", err)
		http_response.SendErrJSON(w, err)
		return
	}

	merchants, err := seed.SeedMerchant(s.merchantUsecase, customers)
	if err != nil {
		log.Printf("Seed Merchant : failed cause, %s \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	log.Println("Seed Merchant: success!")
	http_response.SendOkJSON(w, http.StatusCreated, merchants)
}

func (s *seederAPI) SeedProduct(w http.ResponseWriter, r *http.Request) {
	ctxFetch, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.SetOutput(os.Stdout)
	log.Println("Seed Product : starting!")

	merchants, err := s.merchantUsecase.Fetch(ctxFetch, "", 100, adapter.MerchantSearchOptions{})
	if err != nil {
		log.Printf("Seed Product : failed, cause %s \n ", err)
		http_response.SendErrJSON(w, err)
		return
	}

	products, err := seed.SeedProduct(s.productUsecase, merchants)
	if err != nil {
		log.Printf("Seed Product : failed cause, %s \n", err)
		http_response.SendErrJSON(w, err)
		return
	}

	log.Println("Seed Product: success!")
	http_response.SendOkJSON(w, http.StatusCreated, products)
}
