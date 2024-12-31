package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/agustinrabini/items-api-project/src/main/domain/models"

	"github.com/jopitnow/go-jopit-toolkit/gonosql"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	CategoriesDatabaseError = "[%s] Error in DB"
)

var CategoriesItemNotFoundError = apierrors.NewNotFoundApiError("categories not found")

type CategoriesRepository interface {
	Get(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError)
	GetAllCategories(ctx context.Context) ([]models.Category, apierrors.ApiError)
	Create(ctx context.Context, input models.Category) (interface{}, apierrors.ApiError)
	Update(ctx context.Context, input models.Category) (int64, apierrors.ApiError)
	Delete(ctx context.Context, categoryID string) (int64, apierrors.ApiError)
}

type categoriesRepository struct {
	Collection *mongo.Collection
}

func NewCategoriesRepository(Collection *mongo.Collection) CategoriesRepository {
	return &categoriesRepository{Collection: Collection}
}

func (storage *categoriesRepository) Get(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
	var category models.Category

	result, err := gonosql.Get(ctx, storage.Collection, categoryID)
	if err != nil {
		return models.Category{}, apierrors.NewInternalServerApiError(fmt.Sprintf(CategoriesDatabaseError, "Get"), err)
	}

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return models.Category{}, CategoriesItemNotFoundError
	}

	if result.Err() != nil { // coverage-ignore
		return models.Category{}, apierrors.NewInternalServerApiError(fmt.Sprintf(CategoriesDatabaseError, "Get"), result.Err())
	}

	err = result.Decode(&category)
	if err != nil {
		return models.Category{}, apierrors.NewInternalServerApiError(fmt.Sprintf(CategoriesDatabaseError, "Get"), err)
	}

	return category, nil
}

func (storage *categoriesRepository) GetAllCategories(ctx context.Context) ([]models.Category, apierrors.ApiError) {
	var categories []models.Category

	filter := bson.M{}

	cursor, err := storage.Collection.Find(ctx, filter)
	if err != nil {
		return nil, apierrors.NewInternalServerApiError(fmt.Sprintf(CategoriesDatabaseError, "GetAll"), err)
	}

	if err = cursor.All(ctx, &categories); err != nil {
		return nil, apierrors.NewInternalServerApiError(fmt.Sprintf(CategoriesDatabaseError, "GetAll"), err)
	}

	return categories, nil
}

func (storage *categoriesRepository) Create(ctx context.Context, input models.Category) (interface{}, apierrors.ApiError) {
	result, err := gonosql.InsertOne(ctx, storage.Collection, input)
	if err != nil {
		return nil, apierrors.NewInternalServerApiError(fmt.Sprintf(CategoriesDatabaseError, "Save"), err)
	}

	if result.InsertedID == nil || result.InsertedID == "" { // coverage-ignore
		return -1, apierrors.NewInternalServerApiError(fmt.Sprintf(CategoriesDatabaseError, "Save"), errors.New("category not created"))
	}

	return result.InsertedID, nil
}

func (storage *categoriesRepository) Update(ctx context.Context, input models.Category) (int64, apierrors.ApiError) {
	primitiveID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return -1, apierrors.NewInternalServerApiError(fmt.Sprintf(CategoriesDatabaseError, "Update"), err)
	}

	update := bson.M{
		"$set": bson.M{
			"name": input.Name,
		},
	}

	result, err := storage.Collection.UpdateOne(ctx, bson.M{"_id": primitiveID}, update)
	if err != nil {
		return -1, apierrors.NewInternalServerApiError(fmt.Sprintf(CategoriesDatabaseError, "Update"), err)
	}

	if result.ModifiedCount == 0 {
		return -1, apierrors.NewNotFoundApiError(fmt.Sprintf(CategoriesDatabaseError, "Update"))
	}

	return result.ModifiedCount, nil
}

func (storage *categoriesRepository) Delete(ctx context.Context, categoryID string) (int64, apierrors.ApiError) {
	result, err := gonosql.Delete(ctx, storage.Collection, categoryID)
	if err != nil {
		return -1, apierrors.NewInternalServerApiError(fmt.Sprintf(CategoriesDatabaseError, "Delete"), err)
	}

	if result.DeletedCount == 0 {
		return -1, apierrors.NewNotFoundApiError(fmt.Sprintf(CategoriesDatabaseError, "Delete"))
	}

	return result.DeletedCount, nil
}
