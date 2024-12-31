package dto

import (
	"fmt"

	"github.com/agustinrabini/items-api-project/src/main/domain/models"
	"github.com/mitchellh/mapstructure"
)

type ItemDTO struct {
	Name        string        `json:"name" binding:"required"`
	Description string        `json:"description" binding:"required"`
	UserID      string        `json:"user_id"`
	Status      string        `json:"status"`
	Category    CategoryDTO   `json:"category" binding:"required"`
	Price       PriceDTO      `json:"price" binding:"required"`
	Images      []ImageDTO    `json:"images"`
	Attributes  AttributesDTO `json:"attributes"`
	Eligible    []EligibleDTO `json:"eligible"`
}

type ImageDTO string
type AttributesDTO map[string]string

func (i ItemDTO) ToItem() (models.Item, error) {
	item := models.Item{}
	err := mapstructure.Decode(i, &item)
	if err != nil {
		return models.Item{}, err
	}
	if item.Price.Amount == 0 {
		return models.Item{}, fmt.Errorf("error maping to item on service")
	}
	return item, nil
}
