package dto

import "github.com/agustinrabini/items-api-project/src/main/domain/models"

type CategoryDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name" binding:"required"`
}

type CategoriesDTO struct {
	CategoryDTO []models.Category `json:"categories"`
}
