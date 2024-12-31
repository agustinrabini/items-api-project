package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/agustinrabini/items-api-project/src/main/domain/models"
	"github.com/agustinrabini/items-api-project/src/main/domain/repositories"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

var ErrorCategoryExists = apierrors.NewApiError("Error creating the item category. ", fmt.Errorf("category alredy exists. ").Error(), 409, apierrors.CauseList{})

type CategoriesService interface {
	Get(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError)
	GetAllCategories(ctx context.Context) ([]models.Category, apierrors.ApiError)
	Create(ctx context.Context, input models.Category) apierrors.ApiError
	Update(ctx context.Context, input models.Category) apierrors.ApiError
	Delete(ctx context.Context, items []models.Item, categoryID string) apierrors.ApiError
}

type categoriesService struct {
	repository repositories.CategoriesRepository
}

func NewCategoriesService(repository repositories.CategoriesRepository) CategoriesService {
	return &categoriesService{repository: repository}
}

func (s *categoriesService) Get(ctx context.Context, categoryID string) (models.Category, apierrors.ApiError) {
	category, err := s.repository.Get(ctx, categoryID)
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func (s *categoriesService) GetAllCategories(ctx context.Context) ([]models.Category, apierrors.ApiError) {
	categories, err := s.repository.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	if len(categories) <= 0 {
		return []models.Category{}, apierrors.NewApiError("no categories found", "categories should never be nil, please contact and administrator", http.StatusInternalServerError, apierrors.CauseList{})
	}

	return categories, nil
}

func (s *categoriesService) Create(ctx context.Context, input models.Category) apierrors.ApiError {
	categories, err := s.repository.GetAllCategories(ctx)
	if err != nil {
		return err
	}

	err = validateCategoryExistence(input.Name, categories)
	if err != nil {
		return err
	}

	_, err = s.repository.Create(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoriesService) Update(ctx context.Context, input models.Category) apierrors.ApiError {
	categories, err := s.repository.GetAllCategories(ctx)
	if err != nil {
		return err
	}

	err = validateCategoryExistence(input.Name, categories)
	if err != nil {
		return err
	}

	_, err = s.repository.Update(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoriesService) Delete(ctx context.Context, items []models.Item, categoryID string) apierrors.ApiError {

	if len(items) != 0 {

		ids := []string{}

		for _, itm := range items {
			ids = append(ids, itm.ID)
		}

		return apierrors.NewApiError("error attempting to delete the category, update this implementations before deleting the category", "the following items are using the category "+fmt.Sprint(ids), 409, apierrors.CauseList{})
	}

	_, err := s.repository.Delete(ctx, categoryID)
	if err != nil {
		return err
	}

	return nil
}

func validateCategoryExistence(inputCategoryName string, categories []models.Category) apierrors.ApiError {
	for _, c := range categories {
		if strings.ToLower(inputCategoryName) == strings.ToLower(c.Name) {
			return ErrorCategoryExists
		}
	}

	return nil
}
