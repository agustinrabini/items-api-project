package handlers

import (
	"net/http"

	"github.com/agustinrabini/items-api-project/src/main/domain/models"
	"github.com/agustinrabini/items-api-project/src/main/domain/models/dto"
	"github.com/agustinrabini/items-api-project/src/main/domain/services"
	"github.com/agustinrabini/items-api-project/src/main/domain/utils"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CategoriesHandler struct {
	Service      services.CategoriesService
	ItemsService services.ItemsService
}

func NewCategoriesHandler(service services.CategoriesService, itemsService services.ItemsService) CategoriesHandler {
	return CategoriesHandler{
		Service:      service,
		ItemsService: itemsService,
	}
}

// GetAllCategories godoc
// @Summary Categories
// @Description Get All Categories Items
// @Tags Categories
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Category
// @Router /items/categories [get]
func (h CategoriesHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.Service.GetAllCategories(c)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, dto.CategoriesDTO{
		CategoryDTO: categories,
	},
	)
}

// Create CreateCategory godoc
// @Summary Create Category
// @Description Create Category Item in db
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param item body domain.Category true "Add Category Item"
// @Success 200
// @Router /items/categories [post]
func (h CategoriesHandler) Create(c *gin.Context) {
	var input models.Category

	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		apiErr := apierrors.NewBadRequestApiError(err.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	if input.Name == "" {
		apiErr := apierrors.NewBadRequestApiError("category name must not be nil")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	err := h.Service.Create(c, input)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusCreated)
}

// Get GetCategory godoc
// @Summary Get Category
// @Description Get Category
// @Tags Categories
// @Param id_category path string true "Category ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.Category
// @Router /items/category/:id_category [get]
func (h CategoriesHandler) Get(c *gin.Context) {
	categoryID := c.Param("id_category")

	err := utils.ValidateHexID([]string{categoryID})
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	category, err := h.Service.Get(c, categoryID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, category)
}

// Delete DeleteCategory godoc
// @Summary Delete Category
// @Description Delete Category
// @Tags Categories
// @Produce  json
// @Param id_category path string true "Category ID"
// @Success 200
// @Router /items/category/:id_category [delete]
func (h CategoriesHandler) Delete(c *gin.Context) {
	categoryID := c.Param("id_category")

	err := utils.ValidateHexID([]string{categoryID})
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	items, err := h.ItemsService.GetByCategoryID(c, categoryID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	err = h.Service.Delete(c, items, categoryID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Update UpdateCategory godoc
// @Summary Update Category Item
// @Description Update Category Item by ID
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id_category path string true "Category ID"
// @Param item body domain.Category true "Update item"
// @Success 200
// @Router /items/category [put]
func (h CategoriesHandler) Update(c *gin.Context) {
	var input models.Category

	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		apiErr := apierrors.NewBadRequestApiError(err.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	if input.Name == "" {
		apiErr := apierrors.NewBadRequestApiError("category name must not be nil")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	err := utils.ValidateHexID([]string{input.ID})
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	err = h.Service.Update(c, input)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	err = h.ItemsService.UpdateItemsCategories(c, input)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
