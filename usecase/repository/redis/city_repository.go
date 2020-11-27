package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/market-place/domain"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
)

type ResponseRajaOngkir struct {
	data *[]domain.City
	err  error
}

type cityRepo struct {
	db      *redis.Client
	cityKey string
}

func NewCityRepo(db *redis.Client) repository.CityRepository {
	return &cityRepo{
		db:      db,
		cityKey: "city",
	}
}

func (c *cityRepo) GetCity(ctx context.Context, keyword string) ([]domain.City, error) {
	cities, err := c.getCitiesCache(ctx)
	if err != nil {
		fmt.Printf("[DEBUG] : CITY REPO CHECK CITY CACHED %#v \n", err)
		return cities, usecase_error.ErrInternalServerError
	}

	if len(cities) == 0 {
		fmt.Printf("[DEBUG] : %#v \n", "CACHED MISSED")
		cResponseRajaOngkir := make(chan ResponseRajaOngkir, 1)
		defer close(cResponseRajaOngkir)
		go c.fetchCityRajaOngkir(cResponseRajaOngkir, ctx)

		res := <-cResponseRajaOngkir
		if res.err != nil {
			fmt.Printf("[DEBUG] : CITY REPO FETCH RAJA ONGKIR %#v \n", res.err)
			return cities, usecase_error.ErrInternalServerError
		}

		cities = c.filterRegexCityName(res.data, keyword)

		ctxCached, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := c.setCityCache(ctxCached, res.data); err != nil {
			fmt.Printf("[DEBUG] : CITY REPO SET CITY CACHED %#v \n", err)
			return cities, usecase_error.ErrInternalServerError
		}
	} else {
		cities = c.filterRegexCityName(&cities, keyword)
	}

	return cities, nil
}

func (c *cityRepo) getCitiesCache(ctx context.Context) ([]domain.City, error) {
	cities := []domain.City{}

	res := c.db.HGetAll(ctx, c.cityKey)
	data, err := res.Result()
	if err != nil {
		return cities, err
	}

	if data != nil {
		for _, strData := range data {
			arrCity := strings.Split(strData, ":")
			city := domain.City{
				CityID:       arrCity[0],
				CityName:     arrCity[1],
				ProvinceID:   arrCity[2],
				ProvinceName: arrCity[3],
				PostalCode:   arrCity[4],
			}
			cities = append(cities, city)
		}
	}

	return cities, nil
}

func (c *cityRepo) fetchCityRajaOngkir(ch chan ResponseRajaOngkir, ctx context.Context) {
	url := "https://api.rajaongkir.com/starter/city"
	key := os.Getenv("RAJA_ONGKIR_KEY")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("key", key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- ResponseRajaOngkir{
			err: err,
		}
		return
	}
	defer res.Body.Close()

	type Payload struct {
		RajaOngkir struct {
			Results []domain.City `json:"results"`
		} `json:"rajaongkir"`
	}

	payload := Payload{}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		if err != nil {
			ch <- ResponseRajaOngkir{
				err: err,
			}
			return
		}
	}

	ch <- ResponseRajaOngkir{
		data: &payload.RajaOngkir.Results,
	}
}

func (c *cityRepo) filterRegexCityName(cities *[]domain.City, keyword string) []domain.City {
	matched := []domain.City{}
	pattern := fmt.Sprintf("^(?:.*%s.*)$", strings.ToLower(keyword))
	reg := regexp.MustCompile(pattern)

	for _, city := range *cities {
		if match := reg.MatchString(strings.ToLower(city.CityName)); match {
			matched = append(matched, city)
		}
	}

	return matched
}

func (c *cityRepo) setCityCache(ctx context.Context, cities *[]domain.City) error {
	pipe := c.db.TxPipeline()

	keyValueFormatData := []string{}
	for _, city := range *cities {
		key := city.CityName
		value := city.EncodeToString()
		keyValueFormatData = append(keyValueFormatData, key, value)
	}

	res := pipe.HSet(ctx, c.cityKey, keyValueFormatData)
	if _, err := res.Result(); err != nil {
		return err
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return nil
}
