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

type mongoDBTBuyerRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewTBuyerRepository(db *mongo.Database) repository.TBuyerRepository {
	return &mongoDBTBuyerRepository{
		db:             db,
		collectionName: "transaction_buyers",
	}
}

func (c *mongoDBTBuyerRepository) Create(ctx context.Context, tBuyer domain.TBuyer) (domain.TBuyer, error) {
	tBuyer.CreatedAt = time.Now().Truncate(time.Millisecond)
	tBuyer.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := c.db.Collection(c.collectionName).InsertOne(ctx, tBuyer)
	if err != nil {
		return tBuyer, usecase_error.ErrInternalServerError
	}
	return tBuyer, nil
}

func (c *mongoDBTBuyerRepository) Fetch(ctx context.Context, cursor string, num int64, optionsSearch domain.TBuyerSearchOptions) ([]domain.TBuyer, error) {
	var tBuyers []domain.TBuyer

	query := bson.M{}
	if optionsSearch.CustomerID != "" {
		query["customer_id"] = optionsSearch.CustomerID
	}
	if optionsSearch.AdminID != "" {
		query["admin_id"] = optionsSearch.AdminID
	}
	if optionsSearch.Status != "" {
		query["payment_status"] = optionsSearch.Status
	}
	if cursor != "" {
		last, err := time.Parse(time.RFC3339, cursor)
		last = last.Truncate(time.Millisecond)
		if err != nil {
			fmt.Printf("[DEBUG] PARSE CURSOR:  %#v \n", err)
			return tBuyers, err
		}
		query["created_at"] = bson.M{
			"$lt": last,
		}
	}

	cur, err := c.db.Collection(c.collectionName).Find(ctx,
		query,
		options.Find().SetLimit(num),
		options.Find().SetSort(bson.M{"created_at": -1}),
	)
	if err != nil {
		return tBuyers, err
	}

	for cur.Next(ctx) {
		var tBuyer domain.TBuyer
		if err := cur.Decode(&tBuyer); err != nil {
			if err == mongo.ErrNilCursor {
				return tBuyers, nil
			}

			return tBuyers, usecase_error.ErrInternalServerError
		}

		tBuyers = append(tBuyers, tBuyer)
	}

	return tBuyers, nil
}

func (c *mongoDBTBuyerRepository) GetByID(ctx context.Context, id string) (domain.TBuyer, error) {
	query := bson.M{"_id": id}

	var tBuyer domain.TBuyer
	if err := c.db.Collection(c.collectionName).FindOne(ctx, query).Decode(&tBuyer); err != nil {
		if err == mongo.ErrNoDocuments {
			return tBuyer, usecase_error.ErrNotFound
		}

		return tBuyer, usecase_error.ErrInternalServerError
	}

	return tBuyer, nil
}

func (c *mongoDBTBuyerRepository) UpdateOne(ctx context.Context, tBuyer domain.TBuyer) (domain.TBuyer, error) {
	tBuyer.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": tBuyer.ID}
	data := bson.M{
		"$set": bson.M{
			"customer_id":    tBuyer.CustomerID,
			"total_transfer": tBuyer.TotalTransfer,
			"payment_status": tBuyer.PaymentStatus,
			"transfer_photo": tBuyer.TransferPhoto,
			"message":        tBuyer.Message,
			"admin_id":       tBuyer.AdminID,
			"created_at":     tBuyer.CreatedAt,
			"updated_at":     tBuyer.UpdatedAt,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedTBuyer domain.TBuyer
	if err := c.db.Collection(c.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedTBuyer); err != nil {
		if err == mongo.ErrNoDocuments {
			return tBuyer, usecase_error.ErrNotFound
		}

		return tBuyer, usecase_error.ErrInternalServerError
	}

	return updatedTBuyer, nil
}

func (c *mongoDBTBuyerRepository) DeleteOne(ctx context.Context, tBuyer domain.TBuyer) (domain.TBuyer, error) {
	query := bson.M{"_id": tBuyer.ID}

	var updatedTBuyer domain.TBuyer
	if err := c.db.Collection(c.collectionName).FindOneAndDelete(ctx, query).Decode(&updatedTBuyer); err != nil {
		if err == mongo.ErrNoDocuments {
			return updatedTBuyer, usecase_error.ErrNotFound
		}

		return updatedTBuyer, usecase_error.ErrInternalServerError
	}

	return updatedTBuyer, nil
}

func (c *mongoDBTBuyerRepository) DeleteAll(ctx context.Context) error {
	if _, err := c.db.Collection(c.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		return usecase_error.ErrInternalServerError
	}
	return nil
}
