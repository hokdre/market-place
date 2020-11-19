package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDBProductRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewProductRepository(db *mongo.Database) repository.ProductRepository {
	return &mongoDBProductRepository{
		db:             db,
		collectionName: "products",
	}
}

func (p *mongoDBProductRepository) convertToLocalTime(product *domain.Product) {
	product.CreatedAt = product.CreatedAt.Local().Truncate(time.Millisecond)
	product.UpdatedAt = product.UpdatedAt.Local().Truncate(time.Millisecond)
}

func (p *mongoDBProductRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	product.CreatedAt = time.Now().Truncate(time.Millisecond)
	product.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := p.db.Collection(p.collectionName).InsertOne(ctx, product)
	if err != nil {
		fmt.Printf("[REPOSITORY] REPOSITORY PRODUCT CREATE:  %#v \n", err)
		return product, usecase_error.ErrInternalServerError
	}
	p.convertToLocalTime(&product)
	return product, nil
}

func (p *mongoDBProductRepository) Fetch(ctx context.Context, cursor string, num int64, options domain.ProductSearchOptions) ([]domain.Product, error) {
	query := bson.M{}
	if options.Name != "" {
		query["name"] = bson.M{
			"$regex": options.Name,
		}
	}
	if options.Category != "" {
		query["categories"] = bson.M{
			"$regex": options.Category,
		}
	}
	if options.Description != "" {
		query["description"] = bson.M{
			"$regex": options.Description,
		}
	}
	if options.Price != 0 {
		query["price"] = bson.M{
			"$gte": options.Price,
		}
	}
	if options.City != "" {
		query["merchant.address.city"] = options.City
	}
	if options.MerchantID != "" {
		query["merchant._id"] = options.MerchantID
	}
	if options.Etalase != "" {
		query["etalase"] = options.Etalase
	}
	if options.ReviewID != "" {
		query["reviews._id"] = options.ReviewID
	}

	var products []domain.Product
	cur, err := p.db.Collection(p.collectionName).Find(ctx, query)
	if err != nil {
		fmt.Printf("[REPOSITORY] REPOSITORY PRODUCT FETCH:  %#v \n", err)
		return products, err
	}

	for cur.Next(ctx) {
		var product domain.Product
		if err := cur.Decode(&product); err != nil {
			fmt.Printf("[REPOSITORY] REPOSITORY PRODUCT LOOP:  %#v \n", err)
			if err == mongo.ErrNilCursor {
				return products, nil
			}

			return products, usecase_error.ErrInternalServerError
		}
		p.convertToLocalTime(&product)
		products = append(products, product)
	}

	return products, nil
}

func (p *mongoDBProductRepository) GetByID(ctx context.Context, id string) (domain.Product, error) {
	query := bson.M{"_id": id}

	var product domain.Product
	if err := p.db.Collection(p.collectionName).FindOne(ctx, query).Decode(&product); err != nil {
		fmt.Printf("[REPOSITORY] REPOSITORY PRODUCT GETBYID:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return product, usecase_error.ErrNotFound
		}

		return product, usecase_error.ErrInternalServerError
	}
	p.convertToLocalTime(&product)
	return product, nil
}

func (p *mongoDBProductRepository) UpdateOne(ctx context.Context, product domain.Product) (domain.Product, error) {
	product.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": product.ID}
	data := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"weight":      product.Weight,
			"width":       product.Width,
			"height":      product.Height,
			"long":        product.Long,
			"description": product.Description,
			"category":    product.Category,
			"tags":        product.Tags,
			"etalase":     product.Etalase,
			"colors":      product.Colors,
			"sizes":       product.Sizes,
			"photos":      product.Photos,
			"price":       product.Price,
			"stock":       product.Stock,
			"merchant":    product.Merchant,
			"rating":      product.Rating,
			"num_review":  product.NumReview,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedProduct domain.Product
	if err := p.db.Collection(p.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedProduct); err != nil {
		fmt.Printf("[REPOSITORY] REPOSITORY PRODUCT UPDATE ONE:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return product, usecase_error.ErrNotFound
		}

		return product, usecase_error.ErrInternalServerError
	}
	p.convertToLocalTime(&updatedProduct)
	return updatedProduct, nil
}

func (p *mongoDBProductRepository) DeleteOne(ctx context.Context, product domain.Product) (domain.Product, error) {
	query := bson.M{"_id": product.ID}

	var updatedProduct domain.Product
	if err := p.db.Collection(p.collectionName).FindOneAndDelete(ctx, query).Decode(&updatedProduct); err != nil {
		fmt.Printf("[REPOSITORY] REPOSITORY PRODUCT DELETE ONE:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return product, usecase_error.ErrNotFound
		}

		return product, usecase_error.ErrInternalServerError
	}
	p.convertToLocalTime(&updatedProduct)
	return updatedProduct, nil
}

func (p *mongoDBProductRepository) DeleteAll(ctx context.Context) error {
	if _, err := p.db.Collection(p.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		fmt.Printf("[REPOSITORY] REPOSITORY PRODUCT DELETE ALL:  %#v \n", err)

		return usecase_error.ErrInternalServerError
	}
	return nil
}
