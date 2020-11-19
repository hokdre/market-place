package http_api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	infrastructureConfig "github.com/market-place/config/infrastructure_config"
	usecaseConfig "github.com/market-place/config/usecase_config"
)

type HttpConfig interface {
	StartServer(port string, readTimeOut, writeTimeOut time.Duration) error
}

type httpConfig struct {
	Mux *mux.Router
}

func NewHttpAPI(
	r *mux.Router,
	usecaseConfig usecaseConfig.UsecaseConfig,
	infrastructureConf infrastructureConfig.InfrastructureConfig,
) HttpConfig {

	//authentication routing
	{
		authHandler := NewAuthAPI(usecaseConfig.GetAuthUsecase())
		r.HandleFunc("/login-customers", authHandler.CustomerLogin).Methods("POST")
		r.HandleFunc("/login-admins", authHandler.AdminLogin).Methods("POST")
	}

	//city
	{
		cityHandler := NewCityAPI(usecaseConfig.GetCityUsecase())
		r.HandleFunc("/cities", cityHandler.GetCity).Methods("GET")
	}
	//ongkir
	{
		ongkirHandler := NewOngkirAPI(usecaseConfig.GetOngkirUsecase())
		r.HandleFunc("/ongkirs", ongkirHandler.GetOgnkir).Methods("GET")
	}

	//customer routing
	{
		customerHanlder := NewCustomerAPI(
			usecaseConfig.GetCustomersUseCase(),
			usecaseConfig.GetAuthUsecase(),
			infrastructureConf.GetGoogleStorageClient(),
		)
		r.HandleFunc("/customers", customerHanlder.Create).Methods("POST")
		r.HandleFunc("/customers/{id}", customerHanlder.GetByID).Methods("GET")
		r.HandleFunc("/customers", customerHanlder.Fetch).Methods("GET")
		r.HandleFunc("/customers/{id}", customerHanlder.UpdateBiodata).Methods("PUT")
		r.HandleFunc("/customers/{id}/password", customerHanlder.UpdatePassword).Methods("PUT")
		r.HandleFunc("/customers/{id}/photo-profile", customerHanlder.UploadPhotoProfile).Methods("PUT")
		r.HandleFunc("/customers/{id}/addresses", customerHanlder.AddAddress).Methods("POST")
		r.HandleFunc("/customers/{id}/addresses/{addID}", customerHanlder.UpdateAddress).Methods("PUT")
		r.HandleFunc("/customers/{id}/addresses/{addID}", customerHanlder.DeleteAddress).Methods("DELETE")
		r.HandleFunc("/customers/{id}/bank-accounts", customerHanlder.AddBankAccount).Methods("POST")
		r.HandleFunc("/customers/{id}/bank-accounts/{bankID}", customerHanlder.UpdateBankaccount).Methods("PUT")
	}

	//admin routing
	{
		adminHandler := NewAdminAPI(
			usecaseConfig.GetAdminsUseCase(),
			usecaseConfig.GetAuthUsecase(),
			infrastructureConf.GetGoogleStorageClient(),
		)
		r.HandleFunc("/admins", adminHandler.Create).Methods("POST")
		r.HandleFunc("/admins", adminHandler.Fetch).Methods("GET")
		r.HandleFunc("/admins/{id}", adminHandler.GetByID).Methods("GET")
		r.HandleFunc("/admins/{id}", adminHandler.UpdateBiodata).Methods("PUT")
		r.HandleFunc("/admins/{id}/password", adminHandler.UpdatePassword).Methods("PUT")
		r.HandleFunc("/admins/{id}/photo-profile", adminHandler.UploadPhotoProfile).Methods("PUT")
		r.HandleFunc("/admins/{id}/addresses", adminHandler.AddAddress).Methods("POST")
		r.HandleFunc("/admins/{id}/addresses/{addID}", adminHandler.UpdateAddress).Methods("PUT")
		r.HandleFunc("/admins/{id}/addresses/{addID}", adminHandler.DeleteAddress).Methods("DELETE")
	}

	//shippings routing
	{
		shippingHandler := NewShippingAPI(
			usecaseConfig.GetShippingUsecase(),
			usecaseConfig.GetAuthUsecase(),
		)
		r.HandleFunc("/shippings", shippingHandler.Create).Methods("POST")
		r.HandleFunc("/shippings", shippingHandler.Fetch).Methods("GET")
		r.HandleFunc("/shippings/{id}", shippingHandler.GetByID).Methods("GET")
		r.HandleFunc("/shippings/{id}", shippingHandler.UpdateOne).Methods("PUT")
		r.HandleFunc("/shippings/{id}", shippingHandler.DeleteOne).Methods("DELETE")
	}

	//merchant routing
	{
		merchantHandler := NewMerchantAPI(
			usecaseConfig.GetMerchantUseCase(),
			usecaseConfig.GetAuthUsecase(),
			infrastructureConf.GetGoogleStorageClient(),
		)
		r.HandleFunc("/merchants", merchantHandler.Create).Methods("POST")
		r.HandleFunc("/merchants", merchantHandler.Fetch).Methods("GET")
		r.HandleFunc("/merchants/{id}", merchantHandler.GetByID).Methods("GET")
		r.HandleFunc("/merchants/{id}", merchantHandler.UpdateData).Methods("PUT")
		r.HandleFunc("/merchants/{id}/photo-profile", merchantHandler.UploadAvatar).Methods("PUT")
		r.HandleFunc("/merchants/{id}/shippings/{sID}", merchantHandler.AddShipping).Methods("POST")
		r.HandleFunc("/merchants/{id}/shippings/{sID}", merchantHandler.RemoveShipping).Methods("DELETE")
		r.HandleFunc("/merchants/{id}/bank-accounts", merchantHandler.AddBankAccounts).Methods("POST")
		r.HandleFunc("/merchants/{id}/etalase", merchantHandler.AddEtalase).Methods("POST")
		r.HandleFunc("/merchants/{id}/etalase/{etalase}", merchantHandler.DeleteEtalase).Methods("DELETE")
		r.HandleFunc("/merchants/{id}/bank-accounts/{bankID}", merchantHandler.UpdateBankAccount).Methods("PUT")
	}

	//product routing
	{
		productHandler := NewProductAPI(
			usecaseConfig.GetProductUseCase(),
			usecaseConfig.GetAuthUsecase(),
			infrastructureConf.GetGoogleStorageClient(),
		)
		r.HandleFunc("/products", productHandler.Create).Methods("POST")
		r.HandleFunc("/products", productHandler.Fetch).Methods("GET")
		r.HandleFunc("/products/terlaris", productHandler.ProductTerlaris).Methods("GET")
		r.HandleFunc("/products/{id}", productHandler.GetByID).Methods("GET")
		r.HandleFunc("/products/{id}", productHandler.UpdateData).Methods("PUT")
		r.HandleFunc("/products/{id}/photos", productHandler.UploadPhotos).Methods("PUT")
		r.HandleFunc("/products/{id}", productHandler.DeleteOne).Methods("DELETE")
	}

	//cart routing
	{
		cartHandler := NewCartAPI(
			usecaseConfig.GetCartUseCase(),
			usecaseConfig.GetAuthUsecase(),
		)
		r.HandleFunc("/carts/{id}", cartHandler.GetByID).Methods("GET")
		r.HandleFunc("/carts/{id}/items", cartHandler.AddProduct).Methods("POST")
		r.HandleFunc("/carts/{id}/items/{pID}", cartHandler.UpdateItemInCart).Methods("PUT")
		r.HandleFunc("/carts/{id}/items", cartHandler.ClearProduct).Methods("DELETE")
		r.HandleFunc("/carts/{id}/items/{pID}", cartHandler.RemoveProduct).Methods("DELETE")
	}

	//order routing
	{
		orderHandler := NewOrderAPI(
			usecaseConfig.GetOrderUseCase(),
			usecaseConfig.GetAuthUsecase(),
			infrastructureConf.GetGoogleStorageClient(),
		)
		r.HandleFunc("/orders/estimasi-pendapatan", orderHandler.EstimasiPendapatan).Methods("GET")
		r.HandleFunc("/orders/summary", orderHandler.OrderSummary).Methods("GET")
		r.HandleFunc("/orders", orderHandler.Create).Methods("POST")
		r.HandleFunc("/orders/{id}", orderHandler.GetByID).Methods("GET")
		r.HandleFunc("/orders-merchant/{id}", orderHandler.FetchOrderMerchant).Methods("GET")
		r.HandleFunc("/orders-customer/{id}", orderHandler.FetchOrderCustomer).Methods("GET")
		r.HandleFunc("/orders/{id}/resi-number", orderHandler.AddResiNumber).Methods("PUT")
		r.HandleFunc("/orders/{id}/shipping-photo", orderHandler.UploadShippingPhoto).Methods("PUT")
		r.HandleFunc("/orders/{id}/finish-order", orderHandler.FinishOrder).Methods("PUT")
		r.HandleFunc("/orders/{id}/reject-order", orderHandler.RejectOrder).Methods("PUT")
		r.HandleFunc("/orders/{id}/ajukan-sampai", orderHandler.AjukanPaketSampai).Methods("PUT")

	}

	//tbuyer
	{
		tbuyerHandler := NewTBuyerAPI(
			usecaseConfig.GetTBuyerUseCase(),
			usecaseConfig.GetAuthUsecase(),
			infrastructureConf.GetGoogleStorageClient(),
		)
		r.HandleFunc("/transactions-buyers", tbuyerHandler.Fetch).Methods("GET")
		r.HandleFunc("/transactions-buyers/{id}/transfer-photo", tbuyerHandler.UploadTransferPhoto).Methods("PUT")
		r.HandleFunc("/transactions-buyers/{id}/accept-transaction", tbuyerHandler.AcceptTransaction).Methods("PUT")
		r.HandleFunc("/transactions-buyers/{id}/reject-transaction", tbuyerHandler.RejectTransaction).Methods("PUT")

	}

	//search
	{
		searchHandler := NewSearchAPI(
			usecaseConfig.GetSearchUsecase(),
		)
		r.HandleFunc("/merchants/{id}/products", searchHandler.MerchantProductSearch).Methods("GET")
		r.HandleFunc("/search/products", searchHandler.ProductSearch).Methods("GET")
		r.HandleFunc("/suggestion", searchHandler.SuggestionSearch).Methods("GET")
	}

	//review merchant
	{
		reviewMerchantHandler := NewRMerchantAPI(
			usecaseConfig.GetRMerchantUseCase(),
			usecaseConfig.GetAuthUsecase(),
		)
		r.HandleFunc("/reviews-merchants", reviewMerchantHandler.Create).Methods("POST")
		r.HandleFunc("/reviews-merchants", reviewMerchantHandler.Fetch).Methods("GET")
	}

	//review product
	{
		reviewProductHandler := NewRProductAPI(
			usecaseConfig.GetRProductUseCase(),
			usecaseConfig.GetAuthUsecase(),
		)
		r.HandleFunc("/reviews-products", reviewProductHandler.Create).Methods("POST")
		r.HandleFunc("/reviews-products", reviewProductHandler.Fetch).Methods("GET")
	}

	//migrations
	{
		migrantionHandler := NewMigrationAPI(
			infrastructureConf.GetElasticClient(),
		)
		r.HandleFunc("/elastic-product-index", migrantionHandler.ElasticProductIndex).Methods("GET")
		r.HandleFunc("/elastic-merchant-index", migrantionHandler.ElasticMerchantIndex).Methods("GET")
	}

	//seeder
	{
		seederHandler := NewSeederAPI(
			usecaseConfig.GetAdminsUseCase(),
			usecaseConfig.GetShippingUsecase(),
			usecaseConfig.GetCustomersUseCase(),
			usecaseConfig.GetMerchantUseCase(),
			usecaseConfig.GetProductUseCase(),
			usecaseConfig.GetCartUseCase(),
			usecaseConfig.GetOrderUseCase(),
			usecaseConfig.GetTBuyerUseCase(),
			usecaseConfig.GetRMerchantUseCase(),
			usecaseConfig.GetRProductUseCase(),
		)
		r.HandleFunc("/seed-admin", seederHandler.SeedSuperAdmin).Methods("GET")
		r.HandleFunc("/seed-shipping", seederHandler.SeedShipping).Methods("GET")
		r.HandleFunc("/seed-customer", seederHandler.SeedCustomer).Methods("GET")
		r.HandleFunc("/seed-merchant", seederHandler.SeedMerchant).Methods("GET")
		r.HandleFunc("/seed-product", seederHandler.SeedProduct).Methods("GET")

	}

	return &httpConfig{
		Mux: r,
	}

}

func (h *httpConfig) StartServer(port string, readTimeOut, writeTimeOut time.Duration) error {
	headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Content-Type", "token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	server := http.Server{
		Addr:         port,
		Handler:      handlers.CORS(headersOk, originsOk, methodsOk)(h.Mux),
		ReadTimeout:  readTimeOut,
		WriteTimeout: writeTimeOut,
	}

	serverError := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			serverError <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-serverError:
		return fmt.Errorf("error: listening and serving: %s", err)
	case <-shutdown:
		log.Println("Caught Signal, Shutting Down")

		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("error: gracefully shutting down server: %s", err)

			if err := server.Close(); err != nil {
				return fmt.Errorf("error: closing server: %s", err)
			}
		}
	}

	return nil
}
