package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	ONGKIR_KEY_REDIS = "ongkir"

	ORIGIN_KEY_NAME          = "origin"
	DESTINATINATION_KEY_NAME = "destination"
	PROVIDER_KEY_NAME        = "provider"
	SERVICES_KEY_NAME        = "services"
	separator_inner          = "_"
	separator_outer          = ":"
)

type Service struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cost        int64  `json:"cost"`
	Etd         string `json:"etd"`
}

type Ongkir struct {
	Origin      string           `json:"origin"`
	Destination string           `json:"destination"`
	Provider    ShippingProvider `json:"provider"`
	Services    []Service        `json:"services"`
}

func (o *Ongkir) GenerateRedisKey() string {
	return fmt.Sprintf(
		"%s:%s:%s:%s",
		ONGKIR_KEY_REDIS,
		o.Origin,
		o.Destination,
		o.Provider.ID,
	)
}

func (o *Ongkir) EncodeProvider() string {
	encodedProvider := fmt.Sprintf(
		"%s_%s_%s_%s",
		o.Provider.ID,
		o.Provider.Name,
		o.Provider.CreatedAt.Format(time.RFC3339),
		o.Provider.UpdatedAt.Format(time.RFC3339),
	)

	return encodedProvider
}

func (o *Ongkir) DecodeProvider(encodedProvider string) (ShippingProvider, error) {
	provider := ShippingProvider{}
	arrValue := strings.Split(encodedProvider, separator_inner)
	provider.ID = arrValue[0]
	provider.Name = arrValue[1]

	createdAt, err := time.Parse(time.RFC3339, arrValue[2])
	if err != nil {
		return provider, err
	}
	createdAt = createdAt.Truncate(time.Millisecond)
	provider.CreatedAt = createdAt

	updatedAt, err := time.Parse(time.RFC3339, arrValue[3])
	if err != nil {
		return provider, err
	}
	updatedAt = updatedAt.Truncate(time.Millisecond)
	provider.UpdatedAt = updatedAt

	return provider, nil
}

func (o *Ongkir) EncodeServices() string {
	encodedServices := []string{}
	if o.Services != nil {
		for _, service := range o.Services {
			encodedService := fmt.Sprintf(
				"%s_%s_%s_%s",
				service.Name,
				service.Description,
				strconv.Itoa(int(service.Cost)),
				service.Etd,
			)
			encodedServices = append(encodedServices, encodedService)
		}
	}

	return strings.Join(encodedServices, ",")
}

func (o *Ongkir) DecodeServices(encodedServices string) ([]Service, error) {
	strServices := strings.Split(encodedServices, ",")

	services := []Service{}
	for _, strService := range strServices {
		serviceData := strings.Split(strService, separator_inner)

		service := Service{}
		service.Name = serviceData[0]
		service.Description = serviceData[1]
		cost, err := strconv.Atoi(serviceData[2])
		if err != nil {
			return services, err
		}
		service.Cost = int64(cost)
		service.Etd = strings.Join(serviceData[3:], "-")

		services = append(services, service)
	}

	return services, nil
}

func (o *Ongkir) EncodeArrayKeyValueFormat() []string {
	keyValueFormatData := []string{
		ORIGIN_KEY_NAME,
		o.Origin,
		DESTINATINATION_KEY_NAME,
		o.Destination,
		PROVIDER_KEY_NAME,
		o.EncodeProvider(),
		SERVICES_KEY_NAME,
		o.EncodeServices(),
	}

	return keyValueFormatData
}
