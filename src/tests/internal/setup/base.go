package setup

import (
	"github.com/agustinrabini/items-api-project/src/main/api/app"
	"github.com/agustinrabini/items-api-project/src/main/api/dependencies"
	"github.com/agustinrabini/items-api-project/src/main/api/handlers"
	"github.com/gin-gonic/gin"
)

// BuildRouter Helper function to create a router during testing.
func BuildRouter(depend dependencies.HandlersStruct) *gin.Engine {
	router := app.ConfigureRouter()
	mockRouteMapper(router, depend)
	return router
}

func mockRouteMapper(router *gin.Engine, h dependencies.HandlersStruct) {
	// Health
	health := handlers.NewHealthCheckerHandler()
	router.GET("/ping", health.Ping)

	// Items
	router.GET("/items", handlers.LoggerHandler("GetItemsByUserID"), mockAuthFirebase("01-USER-TEST"), h.Items.GetItemsByUserID)
	router.GET("/items/:id", handlers.LoggerHandler("GetItemByID"), h.Items.GetItemByID)
	router.GET("/items/shop/:id", handlers.LoggerHandler("GetItemsByShopID"), h.Items.GetItemsByShopID)
	router.GET("/items/shop/:id/category/:category_id", handlers.LoggerHandler("GetItemsByShopCategoryID"), h.Items.GetItemsByShopCategoryID)
	router.POST("/items/list", handlers.LoggerHandler("GetItemsByIDs"), h.Items.GetItemsByIDs)
	router.POST("/items", handlers.LoggerHandler("CreateItem"), mockAuthFirebase("01-USER-TEST"), h.Items.CreateItem)
	router.DELETE("/items/:id", handlers.LoggerHandler("DeleteItem"), h.Items.DeleteItem)
	router.PUT("/items/:id", handlers.LoggerHandler("UpdateItem"), mockAuthFirebase("01-USER-TEST"), h.Items.UpdateItem)

	//Categories
	router.PUT("/items/category", handlers.LoggerHandler("UpdateCategory"), h.Categories.Update)
	router.DELETE("/items/category/:id_category", handlers.LoggerHandler("DeleteCategory"), h.Categories.Delete)
	router.POST("/items/category", handlers.LoggerHandler("CreateCategory"), h.Categories.Create)
	router.GET("/items/category/:id_category", handlers.LoggerHandler("GetCategory"), h.Categories.Get)
	router.GET("/items/categories", handlers.LoggerHandler("GetAllCategories"), h.Categories.GetAllCategories)
}

func mockAuthFirebase(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
	}
}
