package mongodb

import (
	"context"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDBShippingRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewShippingRepository(db *mongo.Database) repository.ShippingRepository {
	return &mongoDBShippingRepository{
		db:             db,
		collectionName: "shippings",
	}
}

func (s *mongoDBShippingRepository) convertToLocalTime(shipping *domain.ShippingProvider) {
	shipping.CreatedAt = shipping.CreatedAt.Local().Truncate(time.Millisecond)
	shipping.UpdatedAt = shipping.UpdatedAt.Local().Truncate(time.Millisecond)
}

func (s *mongoDBShippingRepository) Create(ctx context.Context, shipping domain.ShippingProvider) (domain.ShippingProvider, error) {
	shipping.CreatedAt = time.Now().Truncate(time.Millisecond)
	shipping.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := s.db.Collection(s.collectionName).InsertOne(ctx, shipping)
	if err != nil {
		return shipping, usecase_error.ErrInternalServerError
	}
	s.convertToLocalTime(&shipping)
	return shipping, nil
}

func (s *mongoDBShippingRepository) Fetch(ctx context.Context, cursor string, num int64, options domain.ShippingProviderSearchOptions) ([]domain.ShippingProvider, error) {
	query := bson.M{}
	if options.Name != "" {
		query["name"] = bson.M{
			"$regex": options.Name,
		}
	}

	var shippings []domain.ShippingProvider
	cur, err := s.db.Collection(s.collectionName).Find(ctx, query)
	if err != nil {
		return shippings, err
	}

	for cur.Next(ctx) {
		var shipping domain.ShippingProvider
		if err := cur.Decode(&shipping); err != nil {
			if err == mongo.ErrNilCursor {
				return shippings, nil
			}

			return shippings, usecase_error.ErrInternalServerError
		}
		s.convertToLocalTime(&shipping)
		shippings = append(shippings, shipping)
	}

	return shippings, nil
}

func (s *mongoDBShippingRepository) GetByID(ctx context.Context, id string) (domain.ShippingProvider, error) {
	query := bson.M{"_id": id}

	var shipping domain.ShippingProvider
	if err := s.db.Collection(s.collectionName).FindOne(ctx, query).Decode(&shipping); err != nil {
		if err == mongo.ErrNoDocuments {
			return shipping, usecase_error.ErrNotFound
		}
		return shipping, usecase_error.ErrInternalServerError
	}
	s.convertToLocalTime(&shipping)
	return shipping, nil
}

func (s *mongoDBShippingRepository) UpdateOne(ctx context.Context, shipping domain.ShippingProvider) (domain.ShippingProvider, error) {
	shipping.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": shipping.ID}
	data := bson.M{
		"$set": bson.M{
			"name":       shipping.Name,
			"created_at": shipping.CreatedAt,
			"updated_at": shipping.UpdatedAt,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedShipping domain.ShippingProvider
	if err := s.db.Collection(s.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedShipping); err != nil {
		if err == mongo.ErrNoDocuments {
			return shipping, usecase_error.ErrNotFound
		}

		return shipping, usecase_error.ErrInternalServerError
	}
	s.convertToLocalTime(&updatedShipping)
	return updatedShipping, nil
}

func (s *mongoDBShippingRepository) DeleteOne(ctx context.Context, shipping domain.ShippingProvider) (domain.ShippingProvider, error) {
	query := bson.M{"_id": shipping.ID}

	var updatedShipping domain.ShippingProvider
	if err := s.db.Collection(s.collectionName).FindOneAndDelete(ctx, query).Decode(&updatedShipping); err != nil {
		if err == mongo.ErrNoDocuments {
			return updatedShipping, usecase_error.ErrNotFound
		}

		return updatedShipping, usecase_error.ErrInternalServerError
	}
	s.convertToLocalTime(&updatedShipping)
	return updatedShipping, nil
}

func (s *mongoDBShippingRepository) DeleteAll(ctx context.Context) error {
	if _, err := s.db.Collection(s.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		return usecase_error.ErrInternalServerError
	}
	return nil
}
