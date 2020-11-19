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

type mongoDBMerchantRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewMerchantRepository(db *mongo.Database) repository.MerchantRepository {
	return &mongoDBMerchantRepository{
		db:             db,
		collectionName: "merchants",
	}
}

func (m *mongoDBMerchantRepository) convertToLocalTime(merchant *domain.Merchant) {
	merchant.CreatedAt = merchant.CreatedAt.Local().Truncate(time.Millisecond)
	merchant.UpdatedAt = merchant.UpdatedAt.Local().Truncate(time.Millisecond)
}

func (m *mongoDBMerchantRepository) Create(ctx context.Context, merchant domain.Merchant) (domain.Merchant, error) {
	merchant.CreatedAt = time.Now().Truncate(time.Millisecond)
	merchant.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := m.db.Collection(m.collectionName).InsertOne(ctx, merchant)
	if err != nil {
		fmt.Printf("[REPOSITORY] MERCHANT CREATE:  %#v \n", err)
		return merchant, usecase_error.ErrInternalServerError
	}

	m.convertToLocalTime(&merchant)
	return merchant, nil
}

func (m *mongoDBMerchantRepository) Fetch(ctx context.Context, cursor string, num int64, options domain.MerchantSearchOptions) ([]domain.Merchant, error) {
	query := bson.M{}
	if options.Name != "" {
		query["name"] = bson.M{
			"$regex": options.Name,
		}
	}
	if options.Description != "" {
		query["description"] = bson.M{
			"$regex": options.Description,
		}
	}
	if options.City != "" {
		query["address.city"] = options.City
	}
	if options.ShippingID != "" {
		query["shippings._id"] = options.ShippingID
	}
	if options.ReviewID != "" {
		query["reviews._id"] = options.ReviewID
	}
	if options.ProductID != "" {
		query["products._id"] = options.ProductID
	}

	var merchants []domain.Merchant
	cur, err := m.db.Collection(m.collectionName).Find(ctx, query)
	if err != nil {
		fmt.Printf("[REPOSITORY] MERCHANT FETCH:  %#v \n", err)
		return merchants, err
	}

	for cur.Next(ctx) {
		var merchant domain.Merchant
		if err := cur.Decode(&merchant); err != nil {
			if err == mongo.ErrNilCursor {
				return merchants, nil
			}

			return merchants, usecase_error.ErrInternalServerError
		}
		m.convertToLocalTime(&merchant)
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

func (m *mongoDBMerchantRepository) GetByID(ctx context.Context, id string) (domain.Merchant, error) {
	query := bson.M{"_id": id}

	var merchant domain.Merchant
	if err := m.db.Collection(m.collectionName).FindOne(ctx, query).Decode(&merchant); err != nil {
		fmt.Printf("[REPOSITORY] MERCHANT GET BY ID:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return merchant, usecase_error.ErrNotFound
		}

		return merchant, usecase_error.ErrInternalServerError
	}

	m.convertToLocalTime(&merchant)
	return merchant, nil
}

func (m *mongoDBMerchantRepository) GetByName(ctx context.Context, name string) (domain.Merchant, error) {
	query := bson.M{"name": name}

	var merchant domain.Merchant
	if err := m.db.Collection(m.collectionName).FindOne(ctx, query).Decode(&merchant); err != nil {
		if err == mongo.ErrNoDocuments {
			return merchant, usecase_error.ErrNotFound
		}
		fmt.Printf("[REPOSITORY] MERCHANT GET BY NAME:  %#v \n", err)
		return merchant, usecase_error.ErrInternalServerError
	}

	m.convertToLocalTime(&merchant)
	return merchant, nil
}

func (m *mongoDBMerchantRepository) UpdateOne(ctx context.Context, merchant domain.Merchant) (domain.Merchant, error) {
	merchant.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": merchant.ID}
	data := bson.M{
		"$set": bson.M{
			"name":           merchant.Name,
			"address":        merchant.Address,
			"avatar":         merchant.Avatar,
			"phone":          merchant.Phone,
			"description":    merchant.Description,
			"etalase":        merchant.Etalase,
			"products":       merchant.Products,
			"bank_accounts":  merchant.BankAccounts,
			"shippings":      merchant.Shippings,
			"rating":         merchant.Rating,
			"num_review":     merchant.NumReview,
			"location_point": merchant.LocationPoint,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedMerchant domain.Merchant
	if err := m.db.Collection(m.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedMerchant); err != nil {
		fmt.Printf("[REPOSITORY] MERCHANT UPDATE:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return merchant, usecase_error.ErrNotFound
		}

		return merchant, usecase_error.ErrInternalServerError
	}

	m.convertToLocalTime(&updatedMerchant)
	return updatedMerchant, nil
}

func (m *mongoDBMerchantRepository) DeleteOne(ctx context.Context, merchant domain.Merchant) (domain.Merchant, error) {
	query := bson.M{"_id": merchant.ID}

	var updatedMerchant domain.Merchant
	if err := m.db.Collection(m.collectionName).FindOneAndDelete(ctx, query).Decode(&updatedMerchant); err != nil {
		fmt.Printf("[REPOSITORY] MERCHANT DELETE ONE:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return merchant, usecase_error.ErrNotFound
		}

		return merchant, usecase_error.ErrInternalServerError
	}

	m.convertToLocalTime(&merchant)
	return updatedMerchant, nil
}

func (m *mongoDBMerchantRepository) DeleteAll(ctx context.Context) error {
	if _, err := m.db.Collection(m.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		fmt.Printf("[REPOSITORY] MERCHANT DELETE ALL:  %#v \n", err)
		return usecase_error.ErrInternalServerError
	}
	return nil
}
