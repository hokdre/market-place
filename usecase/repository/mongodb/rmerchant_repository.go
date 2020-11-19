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

type mongoDBRMerchantRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewRMerchantRepository(db *mongo.Database) repository.RMerchantRepository {
	return &mongoDBRMerchantRepository{
		db:             db,
		collectionName: "review_merchants",
	}
}

func (c *mongoDBRMerchantRepository) Create(ctx context.Context, rMerchant domain.RMerchant) (domain.RMerchant, error) {
	rMerchant.CreatedAt = time.Now().Truncate(time.Millisecond)
	rMerchant.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := c.db.Collection(c.collectionName).InsertOne(ctx, rMerchant)
	if err != nil {
		return rMerchant, usecase_error.ErrInternalServerError
	}
	return rMerchant, nil
}

func (c *mongoDBRMerchantRepository) Fetch(ctx context.Context, cursor string, num int64, optionsSearch domain.RMerchantSearchOptions) ([]domain.RMerchant, error) {
	query := bson.M{}
	if optionsSearch.MerchantID != "" {
		query["merchant_id"] = optionsSearch.MerchantID
	}
	if cursor != "" {
		last, err := time.Parse(time.RFC3339, cursor)
		last = last.Truncate(time.Millisecond)
		if err != nil {
			fmt.Printf("[DEBUG] PARSE CURSOR:  %#v \n", err)
			return []domain.RMerchant{}, err
		}
		query["created_at"] = bson.M{
			"$lt": last,
		}
	}

	var rMerchants []domain.RMerchant
	cur, err := c.db.Collection(c.collectionName).Find(
		ctx,
		query,
		options.Find().SetLimit(num),
		options.Find().SetSort(bson.M{"created_at": -1}),
	)
	if err != nil {
		return rMerchants, err
	}

	for cur.Next(ctx) {
		var rMerchant domain.RMerchant
		if err := cur.Decode(&rMerchant); err != nil {
			if err == mongo.ErrNilCursor {
				return rMerchants, nil
			}

			return rMerchants, usecase_error.ErrInternalServerError
		}

		rMerchants = append(rMerchants, rMerchant)
	}

	return rMerchants, nil
}

func (c *mongoDBRMerchantRepository) GetByID(ctx context.Context, id string) (domain.RMerchant, error) {
	query := bson.M{"_id": id}

	var rMerchant domain.RMerchant
	if err := c.db.Collection(c.collectionName).FindOne(ctx, query).Decode(rMerchant); err != nil {
		if err == mongo.ErrNoDocuments {
			return rMerchant, usecase_error.ErrNotFound
		}

		return rMerchant, usecase_error.ErrInternalServerError
	}

	return rMerchant, nil
}

func (c *mongoDBRMerchantRepository) UpdateOne(ctx context.Context, rMerchant domain.RMerchant) (domain.RMerchant, error) {
	rMerchant.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": rMerchant.ID}
	data := bson.M{
		"$set": bson.M{
			"merchant_id": rMerchant.MerchantID,
			"customer":    rMerchant.Customer,
			"rating":      rMerchant.Rating,
			"comment":     rMerchant.Comment,
			"created_at":  rMerchant.CreatedAt,
			"updated_at":  rMerchant.UpdatedAt,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedRMerchant domain.RMerchant
	if err := c.db.Collection(c.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedRMerchant); err != nil {
		if err == mongo.ErrNoDocuments {
			return rMerchant, usecase_error.ErrNotFound
		}

		return rMerchant, usecase_error.ErrInternalServerError
	}

	return updatedRMerchant, nil
}

func (c *mongoDBRMerchantRepository) DeleteOne(ctx context.Context, rMerchant domain.RMerchant) (domain.RMerchant, error) {
	query := bson.M{"_id": rMerchant.ID}

	var updatedRMerchant domain.RMerchant
	if err := c.db.Collection(c.collectionName).FindOneAndDelete(ctx, query).Decode(&updatedRMerchant); err != nil {
		if err == mongo.ErrNoDocuments {
			return updatedRMerchant, usecase_error.ErrNotFound
		}

		return updatedRMerchant, usecase_error.ErrInternalServerError
	}

	return updatedRMerchant, nil
}

func (c *mongoDBRMerchantRepository) DeleteAll(ctx context.Context) error {
	if _, err := c.db.Collection(c.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		return usecase_error.ErrInternalServerError
	}
	return nil
}
