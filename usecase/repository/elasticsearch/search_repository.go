package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/market-place/domain"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
)

const (
	TAGS_TYPE     = "TAGS"
	PRODUCT_TYPE  = "PRODUCT"
	MERCHANT_TYPE = "MERCHANT"
)

type Response struct {
	data         interface{}
	err          error
	responseType string
}

type elasticSearchRepository struct {
	db            *elasticsearch.Client
	productIndex  string
	merchantIndex string
}

func NewElasticSearchRepository(db *elasticsearch.Client) repository.SearchRepository {
	return &elasticSearchRepository{
		db:            db,
		productIndex:  "ecommerce.products",
		merchantIndex: "ecommerce.merchants",
	}
}

func (e *elasticSearchRepository) getTags(c chan Response, ctx context.Context, keyword string) {
	var suggestBody bytes.Buffer
	suggestQuery := map[string]interface{}{
		"_source": "tags",
		"suggest": map[string]interface{}{
			"tags-suggestion": map[string]interface{}{
				"prefix": keyword,
				"completion": map[string]interface{}{
					"field":           "suggest",
					"skip_duplicates": true,
					"size":            5,
				},
			},
		},
	}
	if err := json.NewEncoder(&suggestBody).Encode(suggestQuery); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION TAGS:  %#v \n", err)
		c <- Response{
			err:          err,
			responseType: TAGS_TYPE,
		}
		return
	}

	res, err := e.db.Search(
		e.db.Search.WithIndex(e.productIndex),
		e.db.Search.WithBody(&suggestBody),
		e.db.Search.WithContext(ctx),
	)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION TAGS:  %#v \n", err)
		c <- Response{
			err:          err,
			responseType: TAGS_TYPE,
		}
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		var resError map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resError); err != nil {
			fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION TAGS:  %#v \n", err)
			c <- Response{
				err:          err,
				responseType: TAGS_TYPE,
			}
			return
		} else {
			fmt.Println(resError)
			fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION TAGS:  ELASTIC RESPONSE ERROR : %s , %s , %s \n",
				res.Status(),
				resError["error"].(map[string]interface{})["type"],
				resError["error"].(map[string]interface{})["type"],
			)
			c <- Response{
				err:          err,
				responseType: TAGS_TYPE,
			}
			return
		}
	}

	type Payload struct {
		Suggest struct {
			TagsSugestions []struct {
				Options []struct {
					Text string `json:"text"`
				} `json:"options"`
			} `json:"tags-suggestion"`
		} `json:"suggest"`
	}

	var payload Payload
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION TAGS:  %#v \n", err)
		c <- Response{
			err:          err,
			responseType: TAGS_TYPE,
		}
		return
	}

	tags := []string{}
	for _, tag := range payload.Suggest.TagsSugestions[0].Options {
		tags = append(tags, tag.Text)
	}

	c <- Response{
		data:         tags,
		err:          nil,
		responseType: TAGS_TYPE,
	}
}

func (e *elasticSearchRepository) getProducts(c chan Response, ctx context.Context, keyword string) {
	var productBody bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"type":   "bool_prefix",
				"fields": []string{"name", "name._2gram", "name._3gram"},
				"query":  keyword,
			},
		},
	}
	if err := json.NewEncoder(&productBody).Encode(query); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION PRODUCT :  %#v \n", err)
		c <- Response{
			err:          err,
			responseType: PRODUCT_TYPE,
		}
		return
	}

	res, err := e.db.Search(
		e.db.Search.WithIndex(e.productIndex),
		e.db.Search.WithBody(&productBody),
		e.db.Search.WithContext(ctx),
	)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION PRODUCT:  %#v \n", err)
		c <- Response{
			err:          err,
			responseType: PRODUCT_TYPE,
		}
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		var resError map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resError); err != nil {
			fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION PRODUCT:  %#v \n", err)
			c <- Response{
				err:          err,
				responseType: PRODUCT_TYPE,
			}
			return
		} else {
			fmt.Println(resError)
			fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION PRODUCT:  ELASTIC RESPONSE ERROR : %s , %s , %s \n",
				res.Status(),
				resError["error"].(map[string]interface{})["type"],
				resError["error"].(map[string]interface{})["type"],
			)
			c <- Response{
				err:          err,
				responseType: PRODUCT_TYPE,
			}
			return
		}
	}

	type Payload struct {
		Hits struct {
			Hits []struct {
				ID      string         `json:"_id"`
				Product domain.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	var payload Payload
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION PRODUCT:  %#v \n", err)
		c <- Response{
			err:          err,
			responseType: PRODUCT_TYPE,
		}
		return
	}

	products := []domain.Product{}
	for _, doc := range payload.Hits.Hits {
		doc.Product.ID = doc.ID
		products = append(products, doc.Product)
	}

	c <- Response{
		data:         products,
		err:          err,
		responseType: PRODUCT_TYPE,
	}
}

func (e *elasticSearchRepository) getMerchants(c chan Response, ctx context.Context, keyword string) {
	var merchantBody bytes.Buffer
	merchantQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"type":   "bool_prefix",
				"fields": []string{"name", "name._2gram", "name._3gram"},
				"query":  keyword,
			},
		},
	}
	if err := json.NewEncoder(&merchantBody).Encode(merchantQuery); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION MERCHANTS:  %#v \n", err)
		c <- Response{
			err:          err,
			responseType: PRODUCT_TYPE,
		}
		return
	}

	res, err := e.db.Search(
		e.db.Search.WithIndex(e.merchantIndex),
		e.db.Search.WithBody(&merchantBody),
		e.db.Search.WithContext(ctx),
	)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION MERCHANTS:  %#v \n", err)
		c <- Response{
			err:          err,
			responseType: MERCHANT_TYPE,
		}
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		var resError map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resError); err != nil {
			fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION MERCHANTS:  %#v \n", err)
			c <- Response{
				err:          err,
				responseType: MERCHANT_TYPE,
			}
			return
		} else {
			fmt.Println(resError)
			fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION MERCHANTS:  ELASTIC RESPONSE ERROR : %s , %s , %s \n",
				res.Status(),
				resError["error"].(map[string]interface{})["type"],
				resError["error"].(map[string]interface{})["type"],
			)
			c <- Response{
				err:          err,
				responseType: MERCHANT_TYPE,
			}
			return
		}
	}

	type Payload struct {
		Hits struct {
			Hits []struct {
				ID       string          `json:"_id"`
				Merchant domain.Merchant `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	var payload Payload
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH SUGGESTION MERCHANTS:  %#v \n", err)
		c <- Response{
			err:          err,
			responseType: MERCHANT_TYPE,
		}
		return
	}

	merchants := []domain.Merchant{}
	for _, doc := range payload.Hits.Hits {
		doc.Merchant.ID = doc.ID
		merchants = append(merchants, doc.Merchant)
	}

	c <- Response{
		data:         merchants,
		err:          err,
		responseType: MERCHANT_TYPE,
	}
}

func (e *elasticSearchRepository) SuggestionSearch(ctx context.Context, keyword string) (domain.Search, error) {
	var search domain.Search

	responseLength := 3
	responses := make(chan Response, responseLength)
	go e.getTags(responses, ctx, keyword)
	go e.getProducts(responses, ctx, keyword)
	go e.getMerchants(responses, ctx, keyword)

	for i := 0; i < responseLength; i++ {
		res := <-responses
		switch res.responseType {
		case TAGS_TYPE:
			if res.err == nil {
				search.Tags = res.data.([]string)
			}
		case PRODUCT_TYPE:
			if res.err == nil {
				search.Products = res.data.([]domain.Product)
			}
		case MERCHANT_TYPE:
			if res.err == nil {
				search.Merchants = res.data.([]domain.Merchant)
			}
		}
	}

	return search, nil
}

func (e *elasticSearchRepository) ProductSearch(ctx context.Context, category, secondCategory, thirdCategory string, city string, min, max int64, keyword string, lastDate string) (domain.SearchProduct, error) {
	var searchedProduct domain.SearchProduct

	keywordQuery := map[string]interface{}{
		"multi_match": map[string]interface{}{
			"type":   "bool_prefix",
			"query":  keyword,
			"fields": []string{"name^3", "name._2gram^3", "name._3gram^3", "tags"},
		},
	}

	mustQuery := []interface{}{}
	if keyword != "" {
		mustQuery = append(mustQuery, keywordQuery)
	}

	categoryQuery := map[string]interface{}{
		"term": map[string]interface{}{
			"category.top": category,
		},
	}
	secondCategoryQuery := map[string]interface{}{
		"term": map[string]interface{}{
			"category.second_sub": secondCategory,
		},
	}
	thirdCategoryQuery := map[string]interface{}{
		"term": map[string]interface{}{
			"category.third_sub": thirdCategory,
		},
	}
	cityQuery := map[string]interface{}{
		"term": map[string]interface{}{
			"merchant.address.city.city_name": city,
		},
	}
	priceQuery := map[string]interface{}{}
	if min == 0 && max != 0 {
		priceQuery["range"] = map[string]interface{}{
			"price": map[string]interface{}{
				"gte": 0,
				"lte": max,
			},
		}
	} else if min != 0 && max == 0 {
		priceQuery["range"] = map[string]interface{}{
			"price": map[string]interface{}{
				"gte": min,
			},
		}
	} else {
		priceQuery["range"] = map[string]interface{}{
			"price": map[string]interface{}{
				"gte": min,
				"lte": max,
			},
		}
	}

	lastDateQuery := map[string]interface{}{
		"range": map[string]interface{}{
			"created_at": map[string]interface{}{
				"lt": lastDate,
			},
		},
	}
	filterQuery := []interface{}{}
	if category != "" {
		filterQuery = append(filterQuery, categoryQuery)
	}
	if secondCategory != "" {
		filterQuery = append(filterQuery, secondCategoryQuery)
	}
	if thirdCategory != "" {
		filterQuery = append(filterQuery, thirdCategoryQuery)
	}
	if city != "" {
		filterQuery = append(filterQuery, cityQuery)
	}
	if min != 0 || max != 0 {
		filterQuery = append(filterQuery, priceQuery)
	}
	if lastDate != "" {
		filterQuery = append(filterQuery, lastDateQuery)
	}

	boostByTags := map[string]interface{}{
		"term": map[string]interface{}{
			"tags": map[string]interface{}{
				"value": keyword,
				"boost": 2,
			},
		},
	}
	boostByDescription := map[string]interface{}{
		"match": map[string]interface{}{
			"description": map[string]interface{}{
				"query":    keyword,
				"analyzer": "standard",
				"boost":    1,
			},
		},
	}
	boostByColor := map[string]interface{}{
		"match": map[string]interface{}{
			"colors": map[string]interface{}{
				"query":    keyword,
				"analyzer": "standard",
				"boost":    1,
			},
		},
	}
	boostBySize := map[string]interface{}{
		"term": map[string]interface{}{
			"sizes": map[string]interface{}{
				"value": keyword,
				"boost": 1,
			},
		},
	}
	shouldQuery := []interface{}{
		boostByTags, boostByColor, boostBySize, boostByDescription,
	}

	categoryAggs := map[string]interface{}{
		"terms": map[string]interface{}{
			"field": "category.top",
		},
		"aggs": map[string]interface{}{
			"second_category": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "category.second_sub",
				},
				"aggs": map[string]interface{}{
					"third_category": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "category.third_sub",
						},
					},
				},
			},
		},
	}
	cityAggs := map[string]interface{}{
		"terms": map[string]interface{}{
			"field": "merchant.address.city.city_name",
		},
	}
	aggsQuery := map[string]interface{}{
		"categories": categoryAggs,
		"cities":     cityAggs,
	}

	var body bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   mustQuery,
				"filter": filterQuery,
				"should": shouldQuery,
			},
		},
		"aggs": aggsQuery,
	}
	if err := json.NewEncoder(&body).Encode(query); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH PRODUCT:  %#v \n", err)
		return searchedProduct, usecase_error.ErrInternalServerError
	}

	res, err := e.db.Search(
		e.db.Search.WithContext(ctx),
		e.db.Search.WithIndex(e.productIndex),
		e.db.Search.WithBody(&body),
		e.db.Search.WithTrackTotalHits(true),
		e.db.Search.WithPretty(),
	)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH PRODUCT:  %#v \n", err)
		return searchedProduct, usecase_error.ErrInternalServerError
	}
	defer res.Body.Close()

	if res.IsError() {
		var resError map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resError); err != nil {
			fmt.Printf("[DEBUG] REPOSITORY SEARCH PRODUCT:  %#v \n", err)
		} else {
			fmt.Printf("[DEBUG] : err : %#v \n", resError)
			fmt.Printf("[DEBUG] REPOSITORY SEARCH PRODUCT:  ELASTIC RESPONSE ERROR : %s , %s , %s \n",
				res.Status(),
				resError["error"].(map[string]interface{})["type"],
				resError["error"].(map[string]interface{})["type"],
			)
		}
		return searchedProduct, usecase_error.ErrInternalServerError
	}

	type Payload struct {
		Hits struct {
			Hits []struct {
				ID      string         `json:"_id"`
				Score   float64        `json:"_score"`
				Product domain.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
		Aggs struct {
			Categories domain.CategorySearch `json:"categories"`
			Cities     domain.CitySearch     `json:"cities"`
		} `json:"aggregations"`
	}
	var payload Payload
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH PRODUCT:  %#v \n", err)
		return searchedProduct, usecase_error.ErrInternalServerError
	}

	for _, doc := range payload.Hits.Hits {
		doc.Product.ID = doc.ID
		searchedProduct.Products = append(searchedProduct.Products, doc.Product)
	}
	searchedProduct.Categories = payload.Aggs.Categories
	searchedProduct.Cities = payload.Aggs.Cities

	return searchedProduct, nil
}

func (e *elasticSearchRepository) ProductTerlarisSearch(ctx context.Context, page int64) {
}

func (e *elasticSearchRepository) ProductTerpopulerSearch(ctx context.Context, page int64) {

}

func (e *elasticSearchRepository) MerchantProductSearch(ctx context.Context, merchantId string, etalase string, productName string, lastDate string, number int64) ([]domain.Product, error) {
	var products []domain.Product

	filterQuery := []map[string]interface{}{
		map[string]interface{}{
			"term": map[string]interface{}{
				"merchant._id": map[string]interface{}{
					"value": merchantId,
				},
			},
		},
		map[string]interface{}{
			"term": map[string]interface{}{
				"etalase": map[string]interface{}{
					"value": etalase,
				},
			},
		},
	}
	if lastDate != "" {
		dateFilter := map[string]interface{}{
			"range": map[string]interface{}{
				"created_at": map[string]interface{}{
					"lt": lastDate,
				},
			},
		}
		filterQuery = append(filterQuery, dateFilter)
	}

	mustQuery := []map[string]interface{}{}
	if productName != "" {
		nameMatch := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"type":   "bool_prefix",
				"fields": []string{"name", "name._2gram", "name._3gram"},
				"query":  productName,
			},
		}
		mustQuery = append(mustQuery, nameMatch)
	}

	var body bytes.Buffer
	defSize := int64(10)
	if number == 0 {
		number = defSize
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   mustQuery,
				"filter": filterQuery,
			},
		},
		"size": number,
	}

	if err := json.NewEncoder(&body).Encode(query); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH MERCHANT PRODUCT:  %#v \n", err)
		return products, usecase_error.ErrInternalServerError
	}

	res, err := e.db.Search(
		e.db.Search.WithContext(ctx),
		e.db.Search.WithIndex(e.productIndex),
		e.db.Search.WithBody(&body),
		e.db.Search.WithTrackTotalHits(true),
		e.db.Search.WithPretty(),
	)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH MERCHANT PRODUCT:  %#v \n", err)
		return products, usecase_error.ErrInternalServerError
	}
	defer res.Body.Close()

	if res.IsError() {
		var resError map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resError); err != nil {
			fmt.Printf("[DEBUG] REPOSITORY SEARCH MERCHANT PRODUCT:  %#v \n", err)
		} else {
			fmt.Printf("[DEBUG] REPOSITORY SEARCH MERCHANT PRODUCT:  ELASTIC RESPONSE ERROR : %s , %s , %s \n",
				res.Status(),
				resError["error"].(map[string]interface{})["type"],
				resError["error"].(map[string]interface{})["type"],
			)
		}
		return products, usecase_error.ErrInternalServerError
	}

	type Payload struct {
		Hits struct {
			Hits []struct {
				ID      string         `json:"_id"`
				Score   float64        `json:"_score"`
				Product domain.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	var payload Payload
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY SEARCH MERCHANT PRODUCT:  %#v \n", err)
		return products, usecase_error.ErrInternalServerError
	}

	for _, doc := range payload.Hits.Hits {
		doc.Product.ID = doc.ID
		products = append(products, doc.Product)
	}

	return products, nil
}
