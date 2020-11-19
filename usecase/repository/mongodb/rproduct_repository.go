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

type mongoDBRProductRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewRProductRepository(db *mongo.Database) repository.RProductRepository {
	return &mongoDBRProductRepository{
		db:             db,
		collectionName: "review_products",
	}
}

func (c *mongoDBRProductRepository) Create(ctx context.Context, rProduct domain.RProduct) (domain.RProduct, error) {
	rProduct.CreatedAt = time.Now().Truncate(time.Millisecond)
	rProduct.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := c.db.Collection(c.collectionName).InsertOne(ctx, rProduct)
	if err != nil {
		return rProduct, usecase_error.ErrInternalServerError
	}
	return rProduct, nil
}

func (c *mongoDBRProductRepository) Fetch(ctx context.Context, cursor string, num int64, optionsSearch domain.RProductSearchOptions) ([]domain.RProduct, error) {
	query := bson.M{}
	if cursor != "" {
		last, err := time.Parse(time.RFC3339, cursor)
		last = last.Truncate(time.Millisecond)
		if err != nil {
			fmt.Printf("[DEBUG] PARSE CURSOR:  %#v \n", err)
			return []domain.RProduct{}, err
		}
		query["created_at"] = bson.M{
			"$lt": last,
		}
	}
	if optionsSearch.ProductID != "" {
		query["product_id"] = optionsSearch.ProductID
	}

	var rProducts []domain.RProduct
	cur, err := c.db.Collection(c.collectionName).Find(
		ctx,
		query,
		options.Find().SetLimit(num),
		options.Find().SetSort(bson.M{"created_at": -1}),
	)
	if err != nil {
		return rProducts, err
	}

	for cur.Next(ctx) {
		var rProduct domain.RProduct
		if err := cur.Decode(&rProduct); err != nil {
			if err == mongo.ErrNilCursor {
				return rProducts, nil
			}

			return rProducts, usecase_error.ErrInternalServerError
		}

		rProducts = append(rProducts, rProduct)
	}

	return rProducts, nil
}

func (c *mongoDBRProductRepository) GetByID(ctx context.Context, id string) (domain.RProduct, error) {
	query := bson.M{"_id": id}

	var rProduct domain.RProduct
	if err := c.db.Collection(c.collectionName).FindOne(ctx, query).Decode(rProduct); err != nil {
		if err == mongo.ErrNoDocuments {
			return rProduct, usecase_error.ErrNotFound
		}

		return rProduct, usecase_error.ErrInternalServerError
	}

	return rProduct, nil
}

func (c *mongoDBRProductRepository) UpdateOne(ctx context.Context, rProduct domain.RProduct) (domain.RProduct, error) {
	rProduct.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": rProduct.ID}
	data := bson.M{
		"$set": bson.M{
			"merchant_id": rProduct.ProductID,
			"customer":    rProduct.Customer,
			"rating":      rProduct.Rating,
			"comment":     rProduct.Comment,
			"created_at":  rProduct.CreatedAt,
			"updated_at":  rProduct.UpdatedAt,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedRProduct domain.RProduct
	if err := c.db.Collection(c.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedRProduct); err != nil {
		if err == mongo.ErrNoDocuments {
			return rProduct, usecase_error.ErrNotFound
		}

		return rProduct, usecase_error.ErrInternalServerError
	}

	return updatedRProduct, nil
}

func (c *mongoDBRProductRepository) DeleteOne(ctx context.Context, rProduct domain.RProduct) (domain.RProduct, error) {
	query := bson.M{"_id": rProduct.ID}

	var updatedRProduct domain.RProduct
	if err := c.db.Collection(c.collectionName).FindOneAndDelete(ctx, query).Decode(&updatedRProduct); err != nil {
		if err == mongo.ErrNoDocuments {
			return updatedRProduct, usecase_error.ErrNotFound
		}

		return updatedRProduct, usecase_error.ErrInternalServerError
	}

	return updatedRProduct, nil
}

func (c *mongoDBRProductRepository) DeleteAll(ctx context.Context) error {
	if _, err := c.db.Collection(c.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		return usecase_error.ErrInternalServerError
	}
	return nil
}
