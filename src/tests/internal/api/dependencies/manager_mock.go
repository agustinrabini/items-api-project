package dependencies

import (
	"github.com/agustinrabini/items-api-project/src/main/api/dependencies"
	"github.com/agustinrabini/items-api-project/src/main/domain/repositories"
	"github.com/agustinrabini/items-api-project/src/tests/internal/api/platform/storage"
	"github.com/jopitnow/go-jopit-toolkit/gonosql"
	"github.com/tryvium-travels/memongo"
)

type DependencyManagerMock struct {
	*gonosql.Data
}

func NewMockDependencyManager(server *memongo.Server) DependencyManagerMock {
	db := storage.OpenNoSQLMock(server)
	if db.Error != nil {
		panic(db.Error)
	}

	return DependencyManagerMock{
		db,
	}
}

func (m DependencyManagerMock) ItemsRepository() repositories.ItemsRepository {
	return repositories.NewItemsRepository(m.NewCollection(dependencies.KvsItemsCollection))
}

func (m DependencyManagerMock) CategoriesRepository() repositories.CategoriesRepository {
	return repositories.NewCategoriesRepository(m.NewCollection(dependencies.KvsCategoriesCollection))
}
