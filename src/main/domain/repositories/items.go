package repositories

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/agustinrabini/items-api-project/src/main/domain/models"
	"github.com/jopitnow/go-jopit-toolkit/gonosql"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	ItemsDatabaseError = "[%s] Error in DB"
)

var ItemNotFoundError = apierrors.NewNotFoundApiError("item not found")
var ItemsNotFoundError = apierrors.NewNotFoundApiError("items not found")

type ItemsRepository interface {
	Get(ctx context.Context, itemID string) (models.Item, apierrors.ApiError)
	GetByUserID(ctx context.Context, userID string) (models.Items, apierrors.ApiError)
	GetByShopID(ctx context.Context, shopID string) (models.Items, apierrors.ApiError)
	GetByShopCategoryID(ctx context.Context, shopID string, categoryID string) (models.Items, apierrors.ApiError)
	GetByIDs(ctx context.Context, itemsIDs []string) (models.Items, apierrors.ApiError)
	Save(ctx context.Context, item models.Item) (interface{}, apierrors.ApiError)
	Update(ctx context.Context, itemID string, updateItem *models.Item) (int64, apierrors.ApiError)
	Delete(ctx context.Context, itemID string) (int64, apierrors.ApiError)

	UpdateItemsCategories(ctx context.Context, category *models.Category) apierrors.ApiError
	GetByCategoryID(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError)
}

type itemsRepository struct {
	Collection *mongo.Collection
}

func NewItemsRepository(collection *mongo.Collection) ItemsRepository {
	return &itemsRepository{Collection: collection}
}

func (storage *itemsRepository) Get(ctx context.Context, itemID string) (models.Item, apierrors.ApiError) {
	var model models.Item

	cursor, err := gonosql.Get(ctx, storage.Collection, itemID)
	if err != nil {
		return models.Item{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), err)
	}

	if errors.Is(cursor.Err(), mongo.ErrNoDocuments) {
		return models.Item{}, ItemNotFoundError
	}

	if cursor.Err() != nil { // coverage-ignore
		return models.Item{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), cursor.Err())
	}

	err = cursor.Decode(&model)
	if err != nil {
		return models.Item{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), err)
	}

	return model, nil
}

func (storage *itemsRepository) GetByUserID(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
	var items []models.Item

	cursor, err := gonosql.GetByKey(ctx, storage.Collection, "user_id", userID)
	if err != nil {
		return models.Items{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), err)
	}

	if cursor.RemainingBatchLength() == 0 {
		return models.Items{}, ItemsNotFoundError
	}

	if err = cursor.All(ctx, &items); err != nil {
		return models.Items{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), err)
	}

	return models.Items{Items: items}, nil
}

func (storage *itemsRepository) GetByShopID(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
	var items []models.Item

	cursor, err := gonosql.GetByKey(ctx, storage.Collection, "shop_id", shopID)
	if err != nil {
		return models.Items{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), err)
	}

	if cursor.RemainingBatchLength() == 0 {
		return models.Items{}, ItemsNotFoundError
	}

	if err = cursor.All(ctx, &items); err != nil {
		return models.Items{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), err)
	}

	return models.Items{Items: items}, nil
}

func (storage *itemsRepository) GetByShopCategoryID(ctx context.Context, shopID string, categoryID string) (models.Items, apierrors.ApiError) {
	var items []models.Item
	var filter = bson.M{"shop_id": shopID, "category._id": categoryID}

	cursor, err := gonosql.GetByFilter(ctx, storage.Collection, filter)
	if err != nil {
		return models.Items{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), err)
	}

	if cursor.RemainingBatchLength() == 0 {
		return models.Items{}, ItemsNotFoundError
	}

	if err = cursor.All(ctx, &items); err != nil {
		return models.Items{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), err)
	}

	return models.Items{Items: items}, nil
}

func (storage *itemsRepository) GetByIDs(ctx context.Context, itemsIDs []string) (models.Items, apierrors.ApiError) {
	var items []models.Item

	cursor, err := gonosql.GetByIDs(ctx, storage.Collection, itemsIDs)
	if err != nil {
		return models.Items{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), err)
	}

	if cursor.RemainingBatchLength() == 0 {
		return models.Items{}, ItemsNotFoundError
	}

	if err = cursor.All(ctx, &items); err != nil {
		return models.Items{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Get"), err)
	}

	return models.Items{Items: items}, nil
}

func (storage *itemsRepository) Save(ctx context.Context, item models.Item) (interface{}, apierrors.ApiError) {
	result, err := gonosql.InsertOne(ctx, storage.Collection, item)
	if err != nil {
		return nil, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Save"), err)
	}

	if result.InsertedID == nil || result.InsertedID == "" { // coverage-ignore
		return -1, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Save"), errors.New("item not created"))
	}

	return result.InsertedID, nil
}

func (storage *itemsRepository) Update(ctx context.Context, itemID string, updateItem *models.Item) (int64, apierrors.ApiError) {
	result, err := gonosql.Update(ctx, storage.Collection, itemID, updateItem)
	if err != nil {
		return -1, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Update"), err)
	}

	if result.MatchedCount == 0 {
		return -1, apierrors.NewNotFoundApiError(fmt.Sprintf(ItemsDatabaseError, "Update"))
	}

	return result.ModifiedCount, nil
}

func (storage *itemsRepository) Delete(ctx context.Context, itemID string) (int64, apierrors.ApiError) {
	result, err := gonosql.Delete(ctx, storage.Collection, itemID)
	if err != nil {
		return -1, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "Delete"), err)
	}

	if result.DeletedCount == 0 {
		return -1, apierrors.NewNotFoundApiError(fmt.Sprintf(ItemsDatabaseError, "Delete"))
	}

	return result.DeletedCount, nil
}

func (storage *itemsRepository) UpdateItemsCategories(ctx context.Context, category *models.Category) apierrors.ApiError {

	filter := bson.M{"category._id": category.ID}
	update := bson.M{"$set": bson.M{
		"category._id":  category.ID,
		"category.name": category.Name,
	}}

	res, err := storage.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "UpdateItemsCategories"), err)
	}

	if res.MatchedCount == 0 {
		return apierrors.NewApiError(fmt.Sprintf(ItemsDatabaseError, "UpdateItemsCategories"), "no update", http.StatusNotFound, apierrors.CauseList{})
	}

	return nil
}

func (storage *itemsRepository) GetByCategoryID(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError) {

	var model []models.Item

	filter := bson.M{"category._id": categoryID}

	cursor, err := storage.Collection.Find(ctx, filter)
	if err != nil {
		return []models.Item{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "GetByCategoryID"), err)
	}

	if cursor.Err() != nil { // coverage-ignore
		return []models.Item{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "GetByCategoryID"), cursor.Err())
	}

	err = cursor.All(ctx, &model)
	if err != nil {
		return []models.Item{}, apierrors.NewInternalServerApiError(fmt.Sprintf(ItemsDatabaseError, "GetByCategoryID"), err)
	}

	return model, nil
}
