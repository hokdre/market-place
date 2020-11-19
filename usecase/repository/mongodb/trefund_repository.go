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

type mongoDBTRefundRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewTRefundRepository(db *mongo.Database) repository.TRefundRepository {
	return &mongoDBTRefundRepository{
		db:             db,
		collectionName: "transaction_refunds",
	}
}

func (c *mongoDBTRefundRepository) Create(ctx context.Context, tRefund domain.TRefund) (domain.TRefund, error) {
	tRefund.CreatedAt = time.Now().Truncate(time.Millisecond)
	tRefund.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := c.db.Collection(c.collectionName).InsertOne(ctx, tRefund)
	if err != nil {
		return tRefund, usecase_error.ErrInternalServerError
	}
	return tRefund, nil
}

func (c *mongoDBTRefundRepository) Fetch(ctx context.Context, cursor string, num int64, options domain.TRefundSearchOptions) ([]domain.TRefund, error) {
	query := bson.M{}
	if options.OrderID != "" {
		query["order_id"] = options.OrderID
	}
	if options.CustomerID != "" {
		query["customer_id"] = options.CustomerID
	}
	if options.AdminID != "" {
		query["admin_id"] = options.AdminID
	}

	var tRefunds []domain.TRefund
	cur, err := c.db.Collection(c.collectionName).Find(ctx, query)
	if err != nil {
		return tRefunds, err
	}

	for cur.Next(ctx) {
		var tRefund domain.TRefund
		if err := cur.Decode(&tRefund); err != nil {
			if err == mongo.ErrNilCursor {
				return tRefunds, nil
			}

			return tRefunds, usecase_error.ErrInternalServerError
		}

		tRefunds = append(tRefunds, tRefund)
	}

	return tRefunds, nil
}

func (c *mongoDBTRefundRepository) GetByID(ctx context.Context, id string) (domain.TRefund, error) {
	query := bson.M{"_id": id}

	var tRefund domain.TRefund
	if err := c.db.Collection(c.collectionName).FindOne(ctx, query).Decode(tRefund); err != nil {
		if err == mongo.ErrNoDocuments {
			return tRefund, usecase_error.ErrNotFound
		}

		return tRefund, usecase_error.ErrInternalServerError
	}

	return tRefund, nil
}

func (c *mongoDBTRefundRepository) UpdateOne(ctx context.Context, tRefund domain.TRefund) (domain.TRefund, error) {
	tRefund.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": tRefund.ID}
	data := bson.M{
		"$set": bson.M{
			"order_id":       tRefund.OrderID,
			"customer_id":    tRefund.CustomerID,
			"total_transfer": tRefund.TotalTransfer,
			"message":        tRefund.Message,
			"admin_id":       tRefund.AdminID,
			"created_at":     tRefund.CreatedAt,
			"updated_at":     tRefund.UpdatedAt,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedTRefund domain.TRefund
	if err := c.db.Collection(c.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedTRefund); err != nil {
		if err == mongo.ErrNoDocuments {
			return tRefund, usecase_error.ErrNotFound
		}

		return tRefund, usecase_error.ErrInternalServerError
	}

	return updatedTRefund, nil
}

func (c *mongoDBTRefundRepository) DeleteOne(ctx context.Context, tRefund domain.TRefund) (domain.TRefund, error) {
	query := bson.M{"_id": tRefund.ID}

	var updatedTRefund domain.TRefund
	if err := c.db.Collection(c.collectionName).FindOneAndDelete(ctx, query).Decode(&updatedTRefund); err != nil {
		if err == mongo.ErrNoDocuments {
			return updatedTRefund, usecase_error.ErrNotFound
		}

		return updatedTRefund, usecase_error.ErrInternalServerError
	}

	return updatedTRefund, nil
}

func (c *mongoDBTRefundRepository) DeleteAll(ctx context.Context) error {
	if _, err := c.db.Collection(c.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		return usecase_error.ErrInternalServerError
	}
	return nil
}
