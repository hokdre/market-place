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

type mongoDBAdminRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewAdminRepository(db *mongo.Database) repository.AdminRepository {
	return &mongoDBAdminRepository{
		db:             db,
		collectionName: "admins",
	}
}

func (a *mongoDBAdminRepository) convertToLocalTime(admin *domain.Admin) {
	admin.CreatedAt = admin.CreatedAt.Local().Truncate(time.Millisecond)
	admin.UpdatedAt = admin.UpdatedAt.Local().Truncate(time.Millisecond)
	admin.BirthDay = admin.BirthDay.Local().Truncate(time.Millisecond)
}

func (a *mongoDBAdminRepository) Create(ctx context.Context, admin domain.Admin) (domain.Admin, error) {
	admin.CreatedAt = time.Now().Truncate(time.Millisecond)
	admin.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := a.db.Collection(a.collectionName).InsertOne(ctx, admin)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ADMIN CREATE:  %#v \n", err)
		return admin, usecase_error.ErrInternalServerError
	}
	a.convertToLocalTime(&admin)
	return admin, nil
}

func (a *mongoDBAdminRepository) Fetch(ctx context.Context, cursor string, num int64, options domain.AdminSearchOptions) ([]domain.Admin, error) {
	query := bson.M{}
	if options.Name != "" {
		pattern := fmt.Sprintf(`^(?:.*%s.*)$`, options.Name)
		query["name"] = bson.M{
			"$regex": primitive.Regex{Pattern: pattern, Options: "i"},
		}
	}
	if options.Email != "" {
		query["email"] = options.Email
	}

	var admins []domain.Admin
	cur, err := a.db.Collection(a.collectionName).Find(ctx, query)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ADMIN FETCH:  %#v \n", err)
		return admins, err
	}

	for cur.Next(ctx) {
		var admin domain.Admin
		if err := cur.Decode(&admin); err != nil {
			fmt.Printf("[DEBUG] REPOSITORY ADMIN FETCH LOOP:  %#v \n", err)
			if err == mongo.ErrNilCursor {
				return admins, nil
			}

			return admins, usecase_error.ErrInternalServerError
		}
		a.convertToLocalTime(&admin)
		admins = append(admins, admin)
	}
	return admins, nil
}

func (a *mongoDBAdminRepository) GetByID(ctx context.Context, id string) (domain.Admin, error) {
	query := bson.M{"_id": id}

	var admin domain.Admin
	if err := a.db.Collection(a.collectionName).FindOne(ctx, query).Decode(&admin); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ADMIN GET BY ID:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return admin, usecase_error.ErrNotFound
		}

		return admin, usecase_error.ErrInternalServerError
	}
	a.convertToLocalTime(&admin)
	return admin, nil
}

func (a *mongoDBAdminRepository) GetByEmail(ctx context.Context, email string) (domain.Admin, error) {
	query := bson.M{"email": email}

	var admin domain.Admin
	if err := a.db.Collection(a.collectionName).FindOne(ctx, query).Decode(admin); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ADMIN GET BY EMAIL:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return admin, usecase_error.ErrNotFound
		}

		return admin, usecase_error.ErrInternalServerError
	}
	a.convertToLocalTime(&admin)
	return admin, nil
}

func (a *mongoDBAdminRepository) UpdateOne(ctx context.Context, admin domain.Admin) (domain.Admin, error) {
	admin.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": admin.ID}
	data := bson.M{
		"$set": bson.M{
			"name":      admin.Name,
			"password":  admin.Password,
			"addresses": admin.Addresses,
			"born":      admin.Born,
			"birth_day": admin.BirthDay,
			"phone":     admin.Phone,
			"avatar":    admin.Avatar,
			"gender":    admin.Gender,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedAdmin domain.Admin
	if err := a.db.Collection(a.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedAdmin); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ADMIN UPDATE:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return admin, usecase_error.ErrNotFound
		}

		return admin, usecase_error.ErrInternalServerError
	}
	a.convertToLocalTime(&updatedAdmin)
	return updatedAdmin, nil
}

func (a *mongoDBAdminRepository) DeleteOne(ctx context.Context, admin domain.Admin) (domain.Admin, error) {
	query := bson.M{"_id": admin.ID}

	var updatedAdmin domain.Admin
	if err := a.db.Collection(a.collectionName).FindOneAndDelete(ctx, query).Decode(&updatedAdmin); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ADMIN DELETE:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return admin, usecase_error.ErrNotFound
		}

		return admin, usecase_error.ErrInternalServerError
	}
	a.convertToLocalTime(&admin)
	return updatedAdmin, nil
}

func (a *mongoDBAdminRepository) DeleteAll(ctx context.Context) error {
	if _, err := a.db.Collection(a.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ADMIN DELETE ALL:  %#v \n", err)
		return usecase_error.ErrInternalServerError
	}
	return nil
}
