package services

import (
	"context"
	"fmt"

	"github.com/agustinrabini/items-api-project/src/main/domain/clients"
	"github.com/agustinrabini/items-api-project/src/main/domain/models"
	"github.com/agustinrabini/items-api-project/src/main/domain/models/dto"
	"github.com/agustinrabini/items-api-project/src/main/domain/repositories"

	"github.com/creasty/defaults"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemsService interface {
	Get(ctx context.Context, itemID string) (models.Item, apierrors.ApiError)
	GetItemsByUserID(ctx context.Context, userID string) (models.Items, apierrors.ApiError)
	GetItemsByShopID(ctx context.Context, shopID string) (models.Items, apierrors.ApiError)
	GetItemsByIDs(ctx context.Context, items models.ItemsIds) (models.Items, apierrors.ApiError)
	GetItemsByShopCategoryID(ctx context.Context, shopID string, categoryID string) (models.Items, apierrors.ApiError)
	Delete(ctx context.Context, itemID string) apierrors.ApiError
	CreateItem(ctx context.Context, itemRequest dto.ItemDTO) (interface{}, apierrors.ApiError)
	Update(ctx context.Context, itemID string, itemRequest dto.ItemDTO) apierrors.ApiError

	UpdateItemsCategories(ctx context.Context, category models.Category) apierrors.ApiError
	GetByCategoryID(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError)
}

type itemsService struct {
	repository   repositories.ItemsRepository
	pricesClient clients.PriceClient
	shopsClient  clients.ShopClient
}

func NewItemsService(repository repositories.ItemsRepository, pricesClient clients.PriceClient, shopsClient clients.ShopClient) ItemsService {
	return &itemsService{repository: repository, pricesClient: pricesClient, shopsClient: shopsClient}
}

func (s *itemsService) Get(ctx context.Context, itemID string) (models.Item, apierrors.ApiError) {
	item, err := s.repository.Get(ctx, itemID)
	if err != nil {
		return models.Item{}, err
	}

	item.Price, err = s.pricesClient.GetPriceByItemID(ctx, itemID)
	if err != nil {
		return models.Item{}, err
	}

	item.Validate()

	return item, nil
}

func (s *itemsService) GetItemsByUserID(ctx context.Context, userID string) (models.Items, apierrors.ApiError) {
	items, err := s.repository.GetByUserID(ctx, userID)
	if err != nil {
		return models.Items{}, err
	}

	response, err := s.pricesClient.GetItemsPrices(ctx, items.GetItemsIds())
	if err != nil {
		return models.Items{}, err
	}

	finalItems := items.SetPriceToItems(response)

	return finalItems, nil
}

func (s *itemsService) GetItemsByShopID(ctx context.Context, shopID string) (models.Items, apierrors.ApiError) {
	items, err := s.repository.GetByShopID(ctx, shopID)
	if err != nil {
		return models.Items{}, err
	}

	response, err := s.pricesClient.GetItemsPrices(ctx, items.GetItemsIds())
	if err != nil {
		return models.Items{}, err
	}

	finalItems := items.SetPriceToItems(response)

	return finalItems, nil
}

func (s *itemsService) GetItemsByShopCategoryID(ctx context.Context, shopID string, categoryID string) (models.Items, apierrors.ApiError) {
	items, err := s.repository.GetByShopCategoryID(ctx, shopID, categoryID)
	if err != nil {
		return models.Items{}, err
	}

	response, err := s.pricesClient.GetItemsPrices(ctx, items.GetItemsIds())
	if err != nil {
		return models.Items{}, err
	}

	finalItems := items.SetPriceToItems(response)

	return finalItems, nil
}

func (s *itemsService) GetItemsByIDs(ctx context.Context, itemsIds models.ItemsIds) (models.Items, apierrors.ApiError) {
	items, err := s.repository.GetByIDs(ctx, itemsIds.Items)
	if err != nil {
		return models.Items{}, err
	}

	response, err := s.pricesClient.GetItemsPrices(ctx, items.GetItemsIds())
	if err != nil {
		return models.Items{}, err
	}

	finalItems := items.SetPriceToItems(response)

	return finalItems, nil
}

func (s *itemsService) CreateItem(ctx context.Context, request dto.ItemDTO) (interface{}, apierrors.ApiError) {
	shopID, apiErr := s.shopsClient.GetShopByUserID(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	item, err := request.ToItem()
	if err != nil {
		return nil, apierrors.NewBadRequestApiError(fmt.Sprintf("error converting dto to domain: %s", err.Error()))
	}

	item.ShopID = shopID.ID

	err = defaults.Set(&item)
	if err != nil { // coverage-ignore
		return nil, apierrors.NewInternalServerApiError("error settling defaults", err)
	}

	item.SetEligibleIDs()

	insertedID, apiErr := s.repository.Save(ctx, item)
	if apiErr != nil {
		return nil, apiErr
	}

	item.Price.ItemID = insertedID.(primitive.ObjectID).Hex()

	apiErr = s.pricesClient.CreatePrice(ctx, &item.Price)
	if apiErr != nil {
		return nil, apiErr
	}

	return insertedID, nil
}

func (s *itemsService) Update(ctx context.Context, itemID string, request dto.ItemDTO) apierrors.ApiError {
	item, err := request.ToItem()
	if err != nil {
		return apierrors.NewBadRequestApiError("Error convert body to domain: " + err.Error())
	}

	oldPrice, apiErr := s.pricesClient.GetPriceByItemID(ctx, itemID)
	if apiErr != nil {
		return apiErr
	}

	item.Price.ID = oldPrice.ID
	item.Price.ItemID = itemID

	_, apiErr = s.repository.Update(ctx, itemID, &item)
	if apiErr != nil {
		return apiErr
	}

	return s.pricesClient.UpdatePrice(ctx, &item.Price)
}

func (s *itemsService) Delete(ctx context.Context, itemID string) apierrors.ApiError {

	item, err := s.repository.Get(ctx, itemID)
	if err != nil {
		return err
	}

	_, err = s.repository.Delete(ctx, itemID)
	if err != nil {
		return err
	}

	return s.pricesClient.DeletePrice(ctx, item.ID)
}

func (s *itemsService) UpdateItemsCategories(ctx context.Context, category models.Category) apierrors.ApiError {
	return s.repository.UpdateItemsCategories(ctx, &category)
}

func (s *itemsService) GetByCategoryID(ctx context.Context, categoryID string) ([]models.Item, apierrors.ApiError) {

	items, err := s.repository.GetByCategoryID(ctx, categoryID)
	if err != nil {
		return []models.Item{}, err
	}

	return items, nil
}
