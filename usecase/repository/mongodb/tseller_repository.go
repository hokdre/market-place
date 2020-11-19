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

type mongoDBTSellerRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewTSellerRepository(db *mongo.Database) repository.TSellerRepository {
	return &mongoDBTSellerRepository{
		db:             db,
		collectionName: "transaction_sellers",
	}
}

func (c *mongoDBTSellerRepository) Create(ctx context.Context, tSeller domain.TSeller) (domain.TSeller, error) {
	tSeller.CreatedAt = time.Now().Truncate(time.Millisecond)
	tSeller.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := c.db.Collection(c.collectionName).InsertOne(ctx, tSeller)
	if err != nil {
		return tSeller, usecase_error.ErrInternalServerError
	}
	return tSeller, nil
}

func (c *mongoDBTSellerRepository) Fetch(ctx context.Context, cursor string, num int64, options domain.TSellerSearchOptions) ([]domain.TSeller, error) {
	query := bson.M{}
	if options.OrderID != "" {
		query["order_id"] = options.OrderID
	}
	if options.MerchantID != "" {
		query["merchant_id"] = options.MerchantID
	}
	if options.AdminID != "" {
		query["admin_id"] = options.AdminID
	}

	var tSellers []domain.TSeller
	cur, err := c.db.Collection(c.collectionName).Find(ctx, query)
	if err != nil {
		return tSellers, err
	}

	for cur.Next(ctx) {
		var tSeller domain.TSeller
		if err := cur.Decode(&tSeller); err != nil {
			if err == mongo.ErrNilCursor {
				return tSellers, nil
			}

			return tSellers, usecase_error.ErrInternalServerError
		}

		tSellers = append(tSellers, tSeller)
	}

	return tSellers, nil
}

func (c *mongoDBTSellerRepository) GetByID(ctx context.Context, id string) (domain.TSeller, error) {
	query := bson.M{"_id": id}

	var tSeller domain.TSeller
	if err := c.db.Collection(c.collectionName).FindOne(ctx, query).Decode(tSeller); err != nil {
		if err == mongo.ErrNoDocuments {
			return tSeller, usecase_error.ErrNotFound
		}

		return tSeller, usecase_error.ErrInternalServerError
	}

	return tSeller, nil
}

func (c *mongoDBTSellerRepository) UpdateOne(ctx context.Context, tSeller domain.TSeller) (domain.TSeller, error) {
	tSeller.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": tSeller.ID}
	data := bson.M{
		"$set": bson.M{
			"order_id":       tSeller.OrderID,
			"merchant_id":    tSeller.MerchantID,
			"total_transfer": tSeller.TotalTransfer,
			"message":        tSeller.Message,
			"admin_id":       tSeller.AdminID,
			"created_at":     tSeller.CreatedAt,
			"updated_at":     tSeller.UpdatedAt,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedTSeller domain.TSeller
	if err := c.db.Collection(c.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedTSeller); err != nil {
		if err == mongo.ErrNoDocuments {
			return tSeller, usecase_error.ErrNotFound
		}

		return tSeller, usecase_error.ErrInternalServerError
	}

	return updatedTSeller, nil
}

func (c *mongoDBTSellerRepository) DeleteOne(ctx context.Context, tSeller domain.TSeller) (domain.TSeller, error) {
	query := bson.M{"_id": tSeller.ID}

	var updatedTSeller domain.TSeller
	if err := c.db.Collection(c.collectionName).FindOneAndDelete(ctx, query).Decode(&updatedTSeller); err != nil {
		if err == mongo.ErrNoDocuments {
			return updatedTSeller, usecase_error.ErrNotFound
		}

		return updatedTSeller, usecase_error.ErrInternalServerError
	}

	return updatedTSeller, nil
}

func (c *mongoDBTSellerRepository) DeleteAll(ctx context.Context) error {
	if _, err := c.db.Collection(c.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		return usecase_error.ErrInternalServerError
	}
	return nil
}
