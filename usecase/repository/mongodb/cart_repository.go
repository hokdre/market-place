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

type mongoDBCartRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewCartRepository(db *mongo.Database) repository.CartRepository {
	return &mongoDBCartRepository{
		db:             db,
		collectionName: "carts",
	}
}

func (c *mongoDBCartRepository) convertToLocalTime(cart *domain.Cart) {
	cart.CreatedAt = cart.CreatedAt.Local().Truncate(time.Millisecond)
	cart.UpdatedAt = cart.UpdatedAt.Local().Truncate(time.Millisecond)
}

func (c *mongoDBCartRepository) Create(ctx context.Context, cart domain.Cart) (domain.Cart, error) {
	cart.CreatedAt = time.Now().Truncate(time.Millisecond)
	cart.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := c.db.Collection(c.collectionName).InsertOne(ctx, cart)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY CART CREATE:  %#v \n", err)
		return cart, usecase_error.ErrInternalServerError
	}
	c.convertToLocalTime(&cart)
	return cart, nil
}

func (c *mongoDBCartRepository) Fetch(ctx context.Context, cursor string, num int64, options domain.CartSearchOptions) ([]domain.Cart, error) {
	query := bson.M{}
	if options.ProductID != "" {
		query["items.$.product._id"] = options.ProductID
	}
	if options.MerchantID != "" {
		query["items.$.merchant._id"] = options.MerchantID
	}

	var carts []domain.Cart
	cur, err := c.db.Collection(c.collectionName).Find(ctx, query)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY CART FETCH:  %#v \n", err)
		return carts, usecase_error.ErrInternalServerError
	}

	for cur.Next(ctx) {
		var cart domain.Cart
		if err := cur.Decode(&cart); err != nil {
			if err == mongo.ErrNilCursor {
				return carts, nil
			}
			fmt.Printf("[DEBUG] REPOSITORY CART FETCH LOOP:  %#v \n", err)
			return carts, usecase_error.ErrInternalServerError
		}
		c.convertToLocalTime(&cart)
		carts = append(carts, cart)
	}

	return carts, nil
}

func (c *mongoDBCartRepository) GetByID(ctx context.Context, id string) (domain.Cart, error) {
	query := bson.M{"_id": id}

	var cart domain.Cart
	if err := c.db.Collection(c.collectionName).FindOne(ctx, query).Decode(&cart); err != nil {
		if err == mongo.ErrNoDocuments {
			return cart, usecase_error.ErrNotFound
		}
		fmt.Printf("[DEBUG] REPOSITORY CART GET BY ID:  %#v \n", err)
		return cart, usecase_error.ErrInternalServerError
	}
	c.convertToLocalTime(&cart)
	return cart, nil
}

func (c *mongoDBCartRepository) UpdateOne(ctx context.Context, cart domain.Cart) (domain.Cart, error) {
	cart.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": cart.ID}
	data := bson.M{
		"$set": bson.M{
			"items":      cart.Items,
			"updated_at": cart.UpdatedAt,
		},
	}

	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updateCart domain.Cart
	if err := c.db.Collection(c.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updateCart); err != nil {
		if err == mongo.ErrNoDocuments {
			return cart, usecase_error.ErrNotFound
		}
		fmt.Printf("[DEBUG] REPOSITORY CART UPDATE ONE:  %#v \n", err)
		return cart, usecase_error.ErrInternalServerError
	}
	c.convertToLocalTime(&updateCart)
	return updateCart, nil
}

func (c *mongoDBCartRepository) DeleteOne(ctx context.Context, cart domain.Cart) (domain.Cart, error) {
	query := bson.M{"_id": cart.ID}

	var deletedCart domain.Cart
	if err := c.db.Collection(c.collectionName).FindOneAndDelete(ctx, query).Decode(&deletedCart); err != nil {
		if err == mongo.ErrNoDocuments {
			return cart, usecase_error.ErrNotFound
		}
		fmt.Printf("[DEBUG] REPOSITORY CART DELETE ONE:  %#v \n", err)
		return cart, usecase_error.ErrInternalServerError
	}

	c.convertToLocalTime(&deletedCart)
	return deletedCart, nil
}

func (c *mongoDBCartRepository) DeleteAll(ctx context.Context) error {
	if _, err := c.db.Collection(c.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY CART DELETE ALL:  %#v \n", err)
		return usecase_error.ErrInternalServerError
	}
	return nil
}
