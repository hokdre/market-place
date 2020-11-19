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

type mongoDBReturRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewReturRepository(db *mongo.Database) repository.ReturRepository {
	return &mongoDBReturRepository{
		db:             db,
		collectionName: "returs",
	}
}

func (c *mongoDBReturRepository) Create(ctx context.Context, retur domain.Retur) (domain.Retur, error) {
	retur.CreatedAt = time.Now().Truncate(time.Millisecond)
	retur.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := c.db.Collection(c.collectionName).InsertOne(ctx, retur)
	if err != nil {
		return retur, usecase_error.ErrInternalServerError
	}
	return retur, nil
}

func (c *mongoDBReturRepository) Fetch(ctx context.Context, cursor string, num int64, options domain.ReturSearchOptions) ([]domain.Retur, error) {
	query := bson.M{}
	if options.OrderID != "" {
		query["order_id"] = options.OrderID
	}

	var returs []domain.Retur
	cur, err := c.db.Collection(c.collectionName).Find(ctx, query)
	if err != nil {
		return returs, err
	}

	for cur.Next(ctx) {
		var retur domain.Retur
		if err := cur.Decode(&retur); err != nil {
			if err == mongo.ErrNilCursor {
				return returs, nil
			}

			return returs, usecase_error.ErrInternalServerError
		}

		returs = append(returs, retur)
	}

	return returs, nil
}

func (c *mongoDBReturRepository) GetByID(ctx context.Context, id string) (domain.Retur, error) {
	query := bson.M{"_id": id}

	var retur domain.Retur
	if err := c.db.Collection(c.collectionName).FindOne(ctx, query).Decode(retur); err != nil {
		if err == mongo.ErrNoDocuments {
			return retur, usecase_error.ErrNotFound
		}

		return retur, usecase_error.ErrInternalServerError
	}

	return retur, nil
}

func (c *mongoDBReturRepository) UpdateOne(ctx context.Context, retur domain.Retur) (domain.Retur, error) {
	retur.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": retur.ID}
	data := bson.M{
		"$set": bson.M{
			"order_id":           retur.OrderID,
			"customer_reason":    retur.CustomerReason,
			"merchant_accepment": retur.MerchantAccepment,
			"merchant_reason":    retur.MerchantReason,
			"shipping":           retur.Shipping,
			"resi_number":        retur.ResiNumber,
			"created_at":         retur.CreatedAt,
			"updated_at":         retur.UpdatedAt,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedRetur domain.Retur
	if err := c.db.Collection(c.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedRetur); err != nil {
		if err == mongo.ErrNoDocuments {
			return retur, usecase_error.ErrNotFound
		}

		return retur, usecase_error.ErrInternalServerError
	}

	return updatedRetur, nil
}

func (c *mongoDBReturRepository) DeleteOne(ctx context.Context, retur domain.Retur) (domain.Retur, error) {
	query := bson.M{"_id": retur.ID}

	var updatedRetur domain.Retur
	if err := c.db.Collection(c.collectionName).FindOneAndDelete(ctx, query).Decode(&updatedRetur); err != nil {
		if err == mongo.ErrNoDocuments {
			return retur, usecase_error.ErrNotFound
		}

		return retur, usecase_error.ErrInternalServerError
	}

	return updatedRetur, nil
}

func (c *mongoDBReturRepository) DeleteAll(ctx context.Context) error {
	if _, err := c.db.Collection(c.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		return usecase_error.ErrInternalServerError
	}
	return nil
}
