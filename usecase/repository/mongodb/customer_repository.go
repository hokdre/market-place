package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDBCustomerRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewCustomerRepository(db *mongo.Database) repository.CustomerRepository {
	return &mongoDBCustomerRepository{
		db:             db,
		collectionName: "users",
	}
}

func (c *mongoDBCustomerRepository) convertToLocalTime(customer *domain.Customer) {
	customer.CreatedAt = customer.CreatedAt.Local().Truncate(time.Millisecond)
	customer.UpdatedAt = customer.UpdatedAt.Local().Truncate(time.Millisecond)
	customer.BirthDay = customer.BirthDay.Local().Truncate(time.Millisecond)
}

func (c *mongoDBCustomerRepository) Create(ctx context.Context, customer domain.Customer) (domain.Customer, error) {
	customer.ID = primitive.NewObjectID().Hex()
	customer.CreatedAt = time.Now().Truncate(time.Millisecond)
	customer.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := c.db.Collection(c.collectionName).InsertOne(ctx, customer)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY CUSTOMER CREATE:  %#v \n", err)
		return customer, usecase_error.ErrInternalServerError
	}
	return customer, nil
}

func (c *mongoDBCustomerRepository) Fetch(ctx context.Context, cursor string, num int64, options domain.CustomerSearchOptions) ([]domain.Customer, error) {
	query := bson.M{}
	if options.Name != "" {
		query["name"] = bson.M{
			"$regex": options.Name,
		}
	}
	if options.Email != "" {
		query["email"] = options.Email
	}

	var customers []domain.Customer
	cur, err := c.db.Collection(c.collectionName).Find(ctx, query)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY CUSTOMER FETCH:  %#v \n", err)
		return customers, err
	}

	for cur.Next(ctx) {
		var customer domain.Customer
		if err := cur.Decode(&customer); err != nil {
			fmt.Printf("[DEBUG] REPOSITORY CUSTOMER FETCH LOOP:  %#v \n", err)
			if err == mongo.ErrNilCursor {
				return customers, nil
			}

			return customers, usecase_error.ErrInternalServerError
		}
		c.convertToLocalTime(&customer)
		customers = append(customers, customer)
	}
	return customers, nil
}

func (c *mongoDBCustomerRepository) GetByID(ctx context.Context, id string) (domain.Customer, error) {
	query := bson.M{"_id": id}

	var customer domain.Customer
	if err := c.db.Collection(c.collectionName).FindOne(ctx, query).Decode(&customer); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY CUSTOMER GET BY ID:  %#v \n", err)

		if err == mongo.ErrNoDocuments {
			return customer, usecase_error.ErrNotFound
		}

		return customer, usecase_error.ErrInternalServerError
	}
	c.convertToLocalTime(&customer)
	return customer, nil
}

func (c *mongoDBCustomerRepository) GetByEmail(ctx context.Context, email string) (domain.Customer, error) {
	query := bson.M{"email": email}

	var customer domain.Customer
	if err := c.db.Collection(c.collectionName).FindOne(ctx, query).Decode(customer); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY CUSTOMER GET BY EMAIL:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return customer, usecase_error.ErrNotFound
		}

		return customer, usecase_error.ErrInternalServerError
	}

	return customer, nil
}

func (c *mongoDBCustomerRepository) UpdateOne(ctx context.Context, customer domain.Customer) (domain.Customer, error) {
	customer.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": customer.ID}
	data := bson.M{
		"$set": bson.M{
			"password":      customer.Password,
			"name":          customer.Name,
			"addresses":     customer.Addresses,
			"merchant_id":   customer.MerchantID,
			"born":          customer.Born,
			"birth_day":     customer.BirthDay,
			"phone":         customer.Phone,
			"avatar":        customer.Avatar,
			"gender":        customer.Gender,
			"bank_accounts": customer.BankAccounts,
			"created_at":    customer.CreatedAt,
			"updated_at":    customer.UpdatedAt,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedUser domain.Customer
	if err := c.db.Collection(c.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedUser); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY CUSTOMER UPDATE:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return customer, usecase_error.ErrNotFound
		}

		return customer, usecase_error.ErrInternalServerError
	}
	c.convertToLocalTime(&updatedUser)
	return updatedUser, nil
}

func (c *mongoDBCustomerRepository) DeleteOne(ctx context.Context, customer domain.Customer) (domain.Customer, error) {
	query := bson.M{"_id": customer.ID}

	var deletedUser domain.Customer
	if err := c.db.Collection(c.collectionName).FindOneAndDelete(ctx, query).Decode(&deletedUser); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY CUSTOMER DELETE:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return customer, usecase_error.ErrNotFound
		}

		return customer, usecase_error.ErrInternalServerError
	}
	c.convertToLocalTime(&deletedUser)
	return deletedUser, nil
}

func (c *mongoDBCustomerRepository) DeleteAll(ctx context.Context) error {
	if _, err := c.db.Collection(c.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY CUSTOMER DELETE ALL:  %#v \n", err)
		return usecase_error.ErrInternalServerError
	}
	return nil
}
