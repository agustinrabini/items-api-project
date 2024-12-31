package dependencies

import (
	"github.com/agustinrabini/items-api-project/src/main/domain/clients"
	"github.com/agustinrabini/items-api-project/src/main/domain/handlers"
	"github.com/agustinrabini/items-api-project/src/main/domain/repositories"
	"github.com/agustinrabini/items-api-project/src/main/domain/services"
)

type Dependencies interface {
	ItemsRepository() repositories.ItemsRepository
	CategoriesRepository() repositories.CategoriesRepository
}

func GetDependencyManager() Dependencies {
	return NewDependencyManager()
}

func BuildDependencies() (HandlersStruct, error) {
	manager := GetDependencyManager()

	// ItemsRepository
	itemsRepository := manager.ItemsRepository()
	categoriesRepository := manager.CategoriesRepository()

	// External Clients
	pricesClient := clients.NewPriceClient()
	shopsClient := clients.NewShopClient()

	// Services
	itemsService := services.NewItemsService(itemsRepository, pricesClient, shopsClient)
	categoriesService := services.NewCategoriesService(categoriesRepository)

	// Handlers
	itemsHandler := handlers.NewItemsHandler(itemsService, categoriesService)
	categoriesHandler := handlers.NewCategoriesHandler(categoriesService, itemsService)

	return HandlersStruct{
		Items:      itemsHandler,
		Categories: categoriesHandler,
	}, nil
}

type HandlersStruct struct {
	Items      handlers.ItemsHandler
	Categories handlers.CategoriesHandler
}
