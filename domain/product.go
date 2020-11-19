package domain

import "time"

type Product struct {
	ID          string                  `json:"_id" bson:"_id" validate:"required"`
	Name        string                  `json:"name" bson:"name" validate:"required"`
	Weight      float64                 `json:"weight" bson:"weight" validate:"min=0"`
	Width       float64                 `json:"width" bson:"width" validate:"min=0"`
	Height      float64                 `json:"height" bson:"height" validate:"min=0"`
	Long        float64                 `json:"long" bson:"long" validate:"min=0"`
	Description string                  `json:"description" bson:"description" validate:"required,max=10000"`
	Etalase     string                  `json:"etalase" bson:"etalase" validate:"required"`
	Category    Category                `json:"category" bson:"category" validate:"category"`
	Tags        []string                `json:"tags" bson:"tags" validate:"min=1"`
	Colors      []string                `json:"colors" bson:"colors" validate:"unique_colors"`
	Sizes       []string                `json:"sizes" bson:"sizes" validate:"unique_sizes"`
	Photos      []string                `json:"photos" bson:"photos" validate:"required"`
	Price       float64                 `json:"price" bson:"price" validate:"min=1"`
	Stock       float64                 `json:"stock" bson:"stock" validate:"min=1"`
	Merchant    DenormalizationMerchant `json:"merchant" bson:"merchant" validate:"required"`
	Reviews     []RProduct              `json:"reviews" bson:"-"`
	Rating      float64                 `json:"rating" bson:"rating"`
	NumReview   float64                 `json:"num_review" bson:"num_review"`
	CreatedAt   time.Time               `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at" bson:"updated_at"`
}

func (p *Product) DenormalizationData() DenormalizationProduct {
	return DenormalizationProduct{
		ID:          p.ID,
		Name:        p.Name,
		Weight:      p.Weight,
		Width:       p.Width,
		Height:      p.Height,
		Long:        p.Long,
		Description: p.Description,
		Etalase:     p.Etalase,
		Tags:        p.Tags,
		Category:    p.Category,
		Colors:      p.Colors,
		Sizes:       p.Sizes,
		Photos:      p.Photos,
		Price:       p.Price,
		Stock:       p.Stock,
		Rating:      p.Rating,
		NumReview:   p.NumReview,
	}
}

type DenormalizationProduct struct {
	ID          string   `json:"_id" bson:"_id"`
	Name        string   `json:"name" bson:"name"`
	Weight      float64  `json:"weight" bson:"weight"`
	Width       float64  `json:"width" bson:"width"`
	Height      float64  `json:"height" bson:"height"`
	Long        float64  `json:"long" bson:"long"`
	Description string   `json:"description" bson:"description"`
	Etalase     string   `json:"etalase" bson:"etalase"`
	Category    Category `json:"category" bson:"category"`
	Tags        []string `json:"tags" bson:"tags"`
	Colors      []string `json:"colors" bson:"colors"`
	Sizes       []string `json:"sizes" bson:"sizes"`
	Photos      []string `json:"photos" bson:"photos"`
	Price       float64  `json:"price" bson:"price"`
	Stock       float64  `json:"stock" bson:"stock"`
	Rating      float64  `json:"rating" bson:"rating"`
	NumReview   float64  `json:"num_review" bson:"num_review"`
}

type ProductSearchOptions struct {
	//product's name contains regex search name keyword
	Name string
	//product's categories have item with category name contains regex search category keyword
	Category string
	//product's description contains regex search keyword
	Tags        string
	Description string
	//product's price in range between search keyword price
	Price int64
	//product's merchant city equals to search city keyword
	City string
	//product'merchant id equals to search merchantID keyword
	MerchantID string
	//product'etalase equals to search etalase keyword
	Etalase string
	//product'reviews have item with id equals to search reviewID keyword
	ReviewID string
}
