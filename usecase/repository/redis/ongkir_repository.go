package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/market-place/domain"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
)

type ResultOngkir struct {
	Ongkir domain.Ongkir
	Error  error
}

type ongkirRepo struct {
	db *redis.Client
}

func NewOngkirRepo(
	db *redis.Client,
) repository.OngkirRepository {
	return &ongkirRepo{
		db: db,
	}
}

func (o *ongkirRepo) GetOngkir(ctx context.Context, origin, destination string, providers []string) ([]domain.Ongkir, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	ongkirs := []domain.Ongkir{}

	producer := func(providers []string) chan string {
		cProdierIds := make(chan string)
		go func() {
			defer close(cProdierIds)
			for _, provider := range providers {
				select {
				case <-ctx.Done():
					return
				default:
					cProdierIds <- provider
				}
			}
		}()

		return cProdierIds
	}
	getCache := func(cProviderIds chan string) chan ResultOngkir {
		cCache := make(chan ResultOngkir)
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			for providerID := range cProviderIds {
				select {
				case <-ctx.Done():
					return
				default:
					wg.Add(1)
					go func(providerId string) {
						defer wg.Done()

						ongkir, err := o.getOngkirCache(ctx, origin, destination, providerId)
						select {
						case <-ctx.Done():
							return
						default:
							cCache <- ResultOngkir{
								Ongkir: ongkir,
								Error:  err,
							}
						}
					}(providerID)
				}
			}
		}()

		go func() {
			wg.Wait()
			close(cCache)
		}()

		return cCache
	}
	fetchRajaOngkir := func(cCache chan ResultOngkir) chan ResultOngkir {
		cFetch := make(chan ResultOngkir)
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			for cache := range cCache {
				select {
				case <-ctx.Done():
					return
				default:
					if cache.Error == nil {
						cFetch <- cache
					} else {
						wg.Add(1)
						go func() {
							defer wg.Done()
							providerID := cache.Ongkir.Provider.ID
							ongkir, err := o.fetchOngkir(ctx, origin, destination, providerID)

							select {
							case <-ctx.Done():
								return
							default:
								cFetch <- ResultOngkir{
									Ongkir: ongkir,
									Error:  err,
								}
							}
						}()
					}

				}
			}
		}()

		go func() {
			wg.Wait()
			close(cFetch)
		}()

		return cFetch
	}

	cProvidersIds := producer(providers)
	cCache := getCache(cProvidersIds)
	results := fetchRajaOngkir(cCache)

	for result := range results {
		if result.Error != nil {
			cancel()
			return ongkirs, result.Error
		}

		ongkirs = append(ongkirs, result.Ongkir)
	}

	return ongkirs, nil
}

func (o *ongkirRepo) setOngkirCache(ctx context.Context, ongkir domain.Ongkir) error {
	pipe := o.db.TxPipeline()

	res := pipe.HSet(ctx, ongkir.GenerateRedisKey(), ongkir.EncodeArrayKeyValueFormat())
	if _, err := res.Result(); err != nil {
		fmt.Printf("[OngkirRepository] : setOngkirCache : %#v \n", err)
		return err
	}

	if _, err := pipe.Exec(ctx); err != nil {
		fmt.Printf("[OngkirRepository] : setOngkirCache : %#v \n", err)
		return err
	}

	return nil

}

func (o *ongkirRepo) fetchOngkir(ctx context.Context, origin string, destination string, providerID string) (domain.Ongkir, error) {
	fmt.Printf("[DEBUG] : %#v \n", "Cached Miss, do fetch ongkir!")
	ongkir := domain.Ongkir{
		Provider: domain.ShippingProvider{
			ID: providerID,
		},
	}

	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 2 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 2 * time.Second,
	}
	netClient := &http.Client{
		Timeout:   2 * time.Second,
		Transport: netTransport,
	}

	key := os.Getenv("RAJA_ONGKIR_KEY")
	url := "https://api.rajaongkir.com/starter/cost"
	payload := strings.NewReader(fmt.Sprintf(
		`origin=%s&destination=%s&weight=1000&courier=%s`,
		origin,
		destination,
		providerID,
	))

	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		fmt.Printf("[OngkirRepository] : FetchOngkir : %#v \n", err)
		return ongkir, err
	}
	req.Header.Add("key", key)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	response, err := netClient.Do(req)
	if err != nil {
		fmt.Printf("[OngkirRepository] : FetchOngkir : %#v \n", err)
		fmt.Printf("[DEBUG] : %#v \n", err.Error())
		return ongkir, err
	}
	defer response.Body.Close()

	type Result struct {
		RajaOngkir struct {
			Status struct {
				Code int64 `json:"code"`
			}
			Results []struct {
				Code  string `json:"code"`
				Name  string `json:"name"`
				Costs []struct {
					Service     string `json:"service"`
					Description string `json:"description"`
					Cost        []struct {
						Value float64 `json:"value"`
						Etd   string  `json:"etd"`
						Note  string  `json:"note"`
					} `json:"cost"`
				} `json:"costs"`
			} `json:"results"`
		} `json:"rajaongkir"`
	}

	result := Result{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		fmt.Printf("[OngkirRepository] : FetchOngkir : %#v \n", err)
		return ongkir, err
	}

	providerData := result.RajaOngkir.Results[0]
	servicesData := providerData.Costs
	services := []domain.Service{}
	for _, serviceData := range servicesData {
		service := domain.Service{}

		service.Name = serviceData.Service
		service.Description = serviceData.Description
		service.Cost = int64(serviceData.Cost[0].Value)
		service.Etd = serviceData.Cost[0].Etd

		services = append(services, service)
	}

	strCreatedAtRFC3339 := time.Now().Format(time.RFC3339)
	dateCreatedAt, _ := time.Parse(time.RFC3339, strCreatedAtRFC3339)
	dateCreatedAt = dateCreatedAt.Truncate(time.Millisecond)

	strUpdatedAtRFC3339 := time.Now().Format(time.RFC3339)
	dateUpdatedAt, _ := time.Parse(time.RFC3339, strUpdatedAtRFC3339)
	dateUpdatedAt = dateUpdatedAt.Truncate(time.Millisecond)

	ongkir = domain.Ongkir{
		Origin:      origin,
		Destination: destination,
		Provider: domain.ShippingProvider{
			ID:        providerID,
			Name:      providerData.Name,
			CreatedAt: dateCreatedAt,
			UpdatedAt: dateUpdatedAt,
		},
		Services: services,
	}

	o.setOngkirCache(ctx, ongkir)
	return ongkir, nil
}

func (o *ongkirRepo) getOngkirCache(ctx context.Context, origin, destination string, providerID string) (domain.Ongkir, error) {
	fmt.Printf("[DEBUG] : %#v \n", "Cached Catched!")
	ongkir := domain.Ongkir{
		Origin:      origin,
		Destination: destination,
		Provider: domain.ShippingProvider{
			ID: providerID,
		},
		Services: []domain.Service{},
	}

	key := ongkir.GenerateRedisKey()
	res := o.db.HGetAll(ctx, key)
	cacheOngkir, err := res.Result()
	if err != nil {
		fmt.Printf("[OngkirRepository] : GetOngkirCache : %#v \n", err)
		return ongkir, err
	}
	if len(cacheOngkir) == 0 {
		return ongkir, usecase_error.ErrNotFound
	}

	encodedProvider := cacheOngkir[domain.PROVIDER_KEY_NAME]
	if encodedProvider != "" {
		provider, err := ongkir.DecodeProvider(encodedProvider)
		if err != nil {
			fmt.Printf("[OngkirRepository] : GetOngkirCache : %#v \n", err)
			return ongkir, err
		}
		ongkir.Provider = provider
	}

	encodedServices := cacheOngkir[domain.SERVICES_KEY_NAME]
	if encodedServices != "" {
		services, err := ongkir.DecodeServices(encodedServices)
		if err != nil {
			fmt.Printf("[OngkirRepository] : GetOngkirCache : %#v \n", err)
			return ongkir, err
		}
		ongkir.Services = services
	}

	return ongkir, nil
}
