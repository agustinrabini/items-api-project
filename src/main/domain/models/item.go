package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	Name        string     `json:"name" binding:"required" bson:"name,$set,omitempty"`
	ShopID      string     `json:"shop_id" bson:"shop_id,$set,omitempty"`
	UserID      string     `json:"user_id" bson:"user_id,$set,omitempty" binding:"required"`
	Category    Category   `json:"category" binding:"required" bson:"category,$set,omitempty"`
	Price       Price      `json:"price" bson:"-"`
	Description string     `json:"description" bson:"description,$set,omitempty"`
	Status      string     `json:"status" bson:"status,omitempty" default:"active"`
	Images      []Image    `json:"images" bson:"images,$set,omitempty"`
	Attributes  Attributes `json:"attributes" bson:"attributes,$set,omitempty"`
	Eligible    []Eligible `json:"eligible,omitempty" bson:"eligible,$set,omitempty"`
}

type Category struct {
	ID   string `json:"id" bson:"_id,$set,omitempty"`
	Name string `json:"name" binding:"required" bson:"name,$set,required"`
}

type Eligible struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Type       string   `json:"type"`
	IsRequired bool     `json:"is_required"`
	Options    []Option `json:"options"`
}

type Attributes map[string]string
type Option string
type Image string

type ItemsIds struct {
	Items []string `json:"items"`
}

func (i *Item) Validate() {
	if len(i.Eligible) == 0 {
		i.Eligible = []Eligible{}
	}
}

func (i *Item) SetEligibleIDs() {
	for e := range i.Eligible {
		i.Eligible[e].ID = primitive.NewObjectID().Hex()
	}
}

func (i *Items) SetPriceToItems(response Prices) Items {

	iis := Items{}

	for _, pr := range response.Prices {

		for _, item := range i.Items {

			if pr.ItemID == item.ID {
				item.Price = pr
				iis.Items = append(iis.Items, item)
			}
		}
	}

	return iis
}

func (i *Items) GetItemsIds() []string {
	var ids []string

	for _, item := range i.Items {
		ids = append(ids, item.ID)
	}

	return ids
}
