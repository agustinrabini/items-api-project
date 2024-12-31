package app

import (
	"github.com/agustinrabini/items-api-project/src/main/api/dependencies"
	"github.com/agustinrabini/items-api-project/src/main/api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jopitnow/go-jopit-toolkit/goauth"
)

func RouterMapper(router *gin.Engine, h dependencies.HandlersStruct) {
	// Health
	health := handlers.NewHealthCheckerHandler()
	router.GET("/ping", health.Ping)

	// Items
	router.GET("/items", handlers.LoggerHandler("GetItemsByUserID"), goauth.AuthWithFirebase(), h.Items.GetItemsByUserID)
	router.GET("/items/:id", handlers.LoggerHandler("GetItemByID"), h.Items.GetItemByID)
	router.GET("/items/shop/:id", handlers.LoggerHandler("GetItemsByShopID"), h.Items.GetItemsByShopID)
	router.GET("/items/shop/:id/category/:category_id", handlers.LoggerHandler("GetItemsByShopCategoryID"), h.Items.GetItemsByShopCategoryID)
	router.POST("/items/list", handlers.LoggerHandler("GetItemsByIDs"), h.Items.GetItemsByIDs)
	router.POST("/items", handlers.LoggerHandler("CreateItem"), goauth.AuthWithFirebase(), h.Items.CreateItem)
	router.PUT("/items/:id", handlers.LoggerHandler("UpdateItem"), goauth.AuthWithFirebase(), h.Items.UpdateItem)
	router.DELETE("/items/:id", handlers.LoggerHandler("DeleteItem"), goauth.AuthWithFirebase(), h.Items.DeleteItem)

	//Categories
	router.GET("/items/category/:id_category", handlers.LoggerHandler("GetCategory"), h.Categories.Get)
	router.GET("/items/categories", handlers.LoggerHandler("GetAllCategories"), h.Categories.GetAllCategories)
	router.POST("/items/category", goauth.PasswordMiddleware(), handlers.LoggerHandler("CreateCategory"), h.Categories.Create)
	router.PUT("/items/category", goauth.PasswordMiddleware(), handlers.LoggerHandler("UpdateCategory"), h.Categories.Update)
	router.DELETE("/items/category/:id_category", goauth.PasswordMiddleware(), handlers.LoggerHandler("DeleteCategory"), h.Categories.Delete)
}
