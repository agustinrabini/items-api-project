package dependencies

import (
	"github.com/agustinrabini/items-api-project/src/main/api/platform/storage"
	"github.com/agustinrabini/items-api-project/src/main/domain/repositories"

	"github.com/jopitnow/go-jopit-toolkit/gonosql"
)

const (
	KvsItemsCollection      = "items"
	KvsCategoriesCollection = "categories"
)

type DependencyManager struct {
	*gonosql.Data
}

func NewDependencyManager() DependencyManager {
	db := storage.NewNoSQL()
	if db.Error != nil {
		panic(db.Error)
	}
	return DependencyManager{
		db,
	}
}

func (m DependencyManager) ItemsRepository() repositories.ItemsRepository {
	return repositories.NewItemsRepository(m.NewCollection(KvsItemsCollection))
}

func (m DependencyManager) CategoriesRepository() repositories.CategoriesRepository {
	return repositories.NewCategoriesRepository(m.NewCollection(KvsCategoriesCollection))
}
