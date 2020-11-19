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

type mongoDBOrderRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewOrderRepository(db *mongo.Database) repository.OrderRepository {
	return &mongoDBOrderRepository{
		db:             db,
		collectionName: "orders",
	}
}

func (o *mongoDBOrderRepository) convertToLocalTime(order *domain.Order) {
	order.CreatedAt = order.CreatedAt.Local().Truncate(time.Millisecond)
	order.UpdatedAt = order.UpdatedAt.Local().Truncate(time.Millisecond)
}

func (o *mongoDBOrderRepository) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	order.CreatedAt = time.Now().Truncate(time.Millisecond)
	order.UpdatedAt = time.Now().Truncate(time.Millisecond)

	_, err := o.db.Collection(o.collectionName).InsertOne(ctx, order)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ORDER CREATE:  %#v \n", err)
		return order, usecase_error.ErrInternalServerError
	}

	o.convertToLocalTime(&order)
	return order, nil
}

func (o *mongoDBOrderRepository) Fetch(ctx context.Context, cursor string, num int64, optionsSearch domain.OrderSearchOptions) ([]domain.Order, error) {
	var orders []domain.Order

	query := bson.M{}
	if optionsSearch.CustomerID != "" {
		query["customer._id"] = optionsSearch.CustomerID
	}
	if optionsSearch.MerchantID != "" {
		query["merchant._id"] = optionsSearch.MerchantID
	}
	if optionsSearch.ShippingID != "" {
		query["shipping._id"] = optionsSearch.ShippingID
	}
	if optionsSearch.ProductID != "" {
		query["product._id"] = optionsSearch.ProductID
	}
	if optionsSearch.TransactionID != "" {
		query["transaction_id"] = optionsSearch.TransactionID
	}
	if optionsSearch.Status != "" {
		query["status_order"] = optionsSearch.Status
	}
	if cursor != "" {
		last, err := time.Parse(time.RFC3339, cursor)
		last = last.Truncate(time.Millisecond)
		if err != nil {
			fmt.Printf("[DEBUG] PARSE CURSOR:  %#v \n", err)
			return orders, err
		}
		query["created_at"] = bson.M{
			"$lt": last,
		}
	}

	cur, err := o.db.Collection(o.collectionName).Find(ctx, query,
		options.Find().SetLimit(num),
		options.Find().SetSort(bson.M{"created_at": -1}),
	)
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ORDER FETCH:  %#v \n", err)
		return orders, err
	}

	for cur.Next(ctx) {
		var order domain.Order
		if err := cur.Decode(&order); err != nil {
			fmt.Printf("[DEBUG] REPOSITORY ORDER FETCH LOOP:  %#v \n", err)
			if err == mongo.ErrNilCursor {
				return orders, nil
			}

			return orders, usecase_error.ErrInternalServerError
		}

		o.convertToLocalTime(&order)
		orders = append(orders, order)
	}

	return orders, nil
}

func (o *mongoDBOrderRepository) GetByID(ctx context.Context, id string) (domain.Order, error) {
	query := bson.M{"_id": id}

	var order domain.Order
	if err := o.db.Collection(o.collectionName).FindOne(ctx, query).Decode(&order); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ORDER GETBYID:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return order, usecase_error.ErrNotFound
		}

		return order, usecase_error.ErrInternalServerError
	}

	o.convertToLocalTime(&order)
	return order, nil
}

func (o *mongoDBOrderRepository) UpdateOne(ctx context.Context, order domain.Order) (domain.Order, error) {
	order.UpdatedAt = time.Now().Truncate(time.Millisecond)

	query := bson.M{"_id": order.ID}
	data := bson.M{
		"$set": bson.M{
			"transaction_id":    order.TransactionsID,
			"order_items":       order.OrderItems,
			"merchant":          order.Merchant,
			"customer":          order.Customer,
			"shipping":          order.Shipping,
			"shipping_cost":     order.ShippingCost,
			"status_order":      order.StatusOrder,
			"resi_number":       order.ResiNumber,
			"delivered":         order.Delivered,
			"reviewed_merchant": order.ReviewedMerchant,
			"reviewed_product":  order.ReviewedProduct,
			"created_at":        order.CreatedAt,
			"updated_at":        order.UpdatedAt,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1))

	var updatedOrder domain.Order
	if err := o.db.Collection(o.collectionName).FindOneAndUpdate(ctx, query, data, opt).Decode(&updatedOrder); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ORDER UPDATE ONE:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return order, usecase_error.ErrNotFound
		}

		return order, usecase_error.ErrInternalServerError
	}

	o.convertToLocalTime(&updatedOrder)
	return updatedOrder, nil
}

func (o *mongoDBOrderRepository) DeleteOne(ctx context.Context, order domain.Order) (domain.Order, error) {
	query := bson.M{"_id": order.ID}

	var deleteOrder domain.Order
	if err := o.db.Collection(o.collectionName).FindOneAndDelete(ctx, query).Decode(&deleteOrder); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ORDER DELETE:  %#v \n", err)
		if err == mongo.ErrNoDocuments {
			return order, usecase_error.ErrNotFound
		}

		return order, usecase_error.ErrInternalServerError
	}

	o.convertToLocalTime(&deleteOrder)
	return deleteOrder, nil
}

func (o *mongoDBOrderRepository) DeleteAll(ctx context.Context) error {
	if _, err := o.db.Collection(o.collectionName).DeleteMany(ctx, bson.M{}); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ORDER DELETE ALL:  %#v \n", err)
		return usecase_error.ErrInternalServerError
	}
	return nil
}

func (o *mongoDBOrderRepository) EstimasiPendapatan(ctx context.Context, merchantID string, startDay string, endDay string) ([]map[string]interface{}, error) {
	summary := []map[string]interface{}{}

	start, err := time.Parse(time.RFC3339, startDay)
	start = start.Truncate(time.Millisecond)
	if err != nil {
		fmt.Printf("[DEBUG] PARSE START DAY:  %#v \n", err)
		return summary, err
	}

	end, err := time.Parse(time.RFC3339, endDay)
	end = end.Truncate(time.Millisecond)
	if err != nil {
		fmt.Printf("[DEBUG] PARSE END DAY:  %#v \n", err)
		return summary, err
	}

	week1Constraint := start.AddDate(0, 0, 7)
	week1Label := "1-7"
	week2Constraint := week1Constraint.AddDate(0, 0, 7)
	week2Label := "8-14"
	week3Constraint := week2Constraint.AddDate(0, 0, 7)
	week3Label := "15-21"
	lastDay := end.Day()
	week4Label := fmt.Sprintf("22-%d", lastDay)

	matchStage := bson.D{
		{"$match", bson.D{
			{"merchant._id", merchantID},
			{"status_order", bson.D{
				{"$ne", domain.STATUS_ORDER_DI_CANCEL},
			},
			},
			{"created_at", bson.D{
				{"$gt", start},
				{"$lt", end},
			},
			},
		},
		},
	}

	projectSubtotalStage := bson.D{
		{
			"$project", bson.D{{
				"sub_total", bson.D{
					{"$sum", bson.D{{
						"$map", bson.D{
							{"input", "$order_items"},
							{"as", "item"},
							{"in", bson.D{
								{"$multiply", bson.A{
									bson.D{{"$ifNull", bson.A{"$$item.quantity", 0}}},
									bson.D{{"$ifNull", bson.A{"$$item.product.price", 0}}},
								}},
							}},
						},
					}}},
				},
			}},
		},
	}

	sumSubtotalStage := bson.D{
		{"$group", bson.D{
			{"_id", bson.D{
				{"$cond", bson.A{
					bson.D{{"created_at", bson.D{{"$lt", bson.A{"$created_at", week1Constraint}}}}},
					week1Label,
					bson.D{
						{"$cond", bson.A{
							bson.D{{"created_at", bson.D{{"$lt", bson.A{"$created_at", week2Constraint}}}}},
							week2Label,
							bson.D{
								{"$cond", bson.A{
									bson.D{{"created_at", bson.D{{"$lt", bson.A{"$created_at", week3Constraint}}}}},
									week3Label,
									week4Label,
								},
								}},
						},
						}},
				}},
			}},
			{
				"total", bson.D{{"$sum", "$sub_total"}},
			},
		}},
	}
	res, err := o.db.Collection(o.collectionName).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		projectSubtotalStage,
		sumSubtotalStage,
	})
	if err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ORDER ESTIMASI PENDAPATAN:  %#v \n", err)
		return summary, usecase_error.ErrInternalServerError
	}

	var results []map[string]interface{}
	if err = res.All(ctx, &results); err != nil {
		fmt.Printf("[DEBUG] REPOSITORY ORDER ESTIMASI PENDAPATAN:  %#v \n", err)
		return summary, usecase_error.ErrInternalServerError
	}

	week1 := map[string]interface{}{
		"label": week1Label,
		"total": int64(0),
	}
	week2 := map[string]interface{}{
		"label": week2Label,
		"total": int64(0),
	}
	week3 := map[string]interface{}{
		"label": week3Label,
		"total": int64(0),
	}
	week4 := map[string]interface{}{
		"label": week4Label,
		"total": int64(0),
	}
	for _, result := range results {
		id := result["_id"].(string)
		total := int64(result["total"].(float64))
		if week1Label == id {
			week1["total"] = total
		} else if week2Label == id {
			week2["total"] = total
		} else if week3Label == id {
			week3["total"] = total
		} else if week4Label == id {
			week4["total"] = total
		}
	}
	summary = append(summary, week1, week2, week3, week4)

	return summary, nil
}

func (o *mongoDBOrderRepository) OrderSummary(ctx context.Context, merchantID string, startDay string, endDay string) (map[string]int64, error) {
	summary := map[string]int64{}

	start, err := time.Parse(time.RFC3339, startDay)
	start = start.Truncate(time.Millisecond)
	if err != nil {
		fmt.Printf("[DEBUG] NumOrderBaru Parse Start Day:  %#v \n", err)
		return summary, err
	}

	end, err := time.Parse(time.RFC3339, endDay)
	end = end.Truncate(time.Millisecond)
	if err != nil {
		fmt.Printf("[DEBUG] NumOrderBaru Parse End Day:  %#v \n", err)
		return summary, err
	}

	matchStage := bson.D{
		{"$match", bson.D{
			{"merchant._id", merchantID},
			{"status_order", bson.D{
				{"$ne", domain.STATUS_ORDER_DI_CANCEL},
			},
			},
			{"created_at", bson.D{
				{"$gt", start},
				{"$lt", end},
			},
			},
		},
		},
	}

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", bson.D{{"status_order", "$status_order"}}},
			{"jumlah", bson.D{
				{"$sum", 1},
			}},
		}},
	}
	res, err := o.db.Collection(o.collectionName).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		groupStage,
	})
	if err != nil {
		fmt.Printf("[DEBUG] NumOrderBaru Fetch Group Order:  %#v \n", err)
		return summary, usecase_error.ErrInternalServerError
	}

	var results []map[string]interface{}
	if err = res.All(ctx, &results); err != nil {
		fmt.Printf("[DEBUG] Result:  %#v \n", err)
		return summary, usecase_error.ErrInternalServerError

	}

	for _, result := range results {
		data := result["_id"].(map[string]interface{})
		status := data["status_order"].(string)
		jumlah := int64(result["jumlah"].(int32))
		summary[status] = jumlah
	}

	return summary, nil
}

func (o *mongoDBOrderRepository) ProductTerlaris(ctx context.Context) ([]map[string]interface{}, error) {
	summary := []map[string]interface{}{}

	matchStage := bson.D{
		{"$match", bson.M{"status_order": domain.STATUS_ORDER_SELESAI}},
	}
	unWindStage := bson.D{{"$unwind", "$order_items"}}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$order_items.product._id"},
			{"jumlah_terjual", bson.M{"$sum": "$order_items.quantity"}},
			{"product", bson.M{"$first": "$order_items.product"}},
			{"merchant", bson.M{"$first": "$merchant"}},
		}},
	}

	res, err := o.db.Collection(o.collectionName).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		unWindStage,
		groupStage,
	})
	if err != nil {
		fmt.Printf("[DEBUG] Product Terlaris Order:  %#v \n", err)
		return summary, usecase_error.ErrInternalServerError
	}

	var results []map[string]interface{}
	if err = res.All(ctx, &results); err != nil {
		fmt.Printf("[DEBUG] Result:  %#v \n", err)
		return summary, usecase_error.ErrInternalServerError

	}

	for _, result := range results {
		product := map[string]interface{}{}
		product["product"] = result["product"]
		product["jumlah_terjual"] = result["jumlah_terjual"]
		product["merchant"] = result["merchant"]
		summary = append(summary, product)
	}

	return summary, nil
}
