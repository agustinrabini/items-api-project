package handlers

import (
	"context"
	"net/http"

	"github.com/agustinrabini/items-api-project/src/main/domain/models"
	"github.com/agustinrabini/items-api-project/src/main/domain/models/dto"
	"github.com/agustinrabini/items-api-project/src/main/domain/services"
	"github.com/agustinrabini/items-api-project/src/main/domain/utils"

	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/go-jopit-toolkit/goutils/logger"
	"github.com/jopitnow/go-jopit-toolkit/tracing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ItemsHandler struct {
	Service           services.ItemsService
	CategoriesService services.CategoriesService
}

func NewItemsHandler(service services.ItemsService, categoriesService services.CategoriesService) ItemsHandler {
	return ItemsHandler{
		Service:           service,
		CategoriesService: categoriesService,
	}
}

// GetItemByID godoc
// @Summary Get details of item id
// @Description Get details of item
// @Tags Items
// @Accept  json
// @Produce  json
// @Param id path string true "Item ID"
// @Success 200 {object} domain.Item
// @Router /items/{id} [get]
func (h ItemsHandler) GetItemByID(c *gin.Context) {
	xTraceId, _ := c.Get("X-Trace-ID")
	ctx := context.WithValue(c.Request.Context(), tracing.XtraceHeaderKey, xTraceId)

	itemID := c.Param("id")
	err := utils.ValidateHexID([]string{itemID})
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	response, err := h.Service.Get(ctx, itemID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetItemsByUserID godoc
// @Summary Get Items by User ID
// @Description Get Items by User ID in Header Authorization
// @Tags Items
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.Items
// @Router /items [get]
func (h ItemsHandler) GetItemsByUserID(c *gin.Context) {
	xTraceId, _ := c.Get("X-Trace-ID")
	ctx := context.WithValue(c.Request.Context(), tracing.XtraceHeaderKey, xTraceId)
	ctx = context.WithValue(ctx, goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	userID, err1 := goauth.GetUserId(c)
	if err1 != nil {
		apiErr := apierrors.NewUnauthorizedApiError(err1.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	response, apiErr := h.Service.GetItemsByUserID(ctx, userID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetItemsByShopID godoc
// @Summary Get details of items by shop ID
// @Description Get details of items
// @Tags Items
// @Accept  json
// @Produce  json
// @Param id path string true "Shop ID"
// @Success 200 {object} domain.Items
// @Router /items/shop/{id} [get]
func (h ItemsHandler) GetItemsByShopID(c *gin.Context) {
	xTraceId, _ := c.Get("X-Trace-ID")
	ctx := context.WithValue(c.Request.Context(), tracing.XtraceHeaderKey, xTraceId)

	shopID := c.Param("id")
	err := utils.ValidateHexID([]string{shopID})
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	response, err := h.Service.GetItemsByShopID(ctx, shopID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetItemsByShopCategoryID godoc
// @Summary Get details of items by shop ID and category ID
// @Description Get details of items
// @Tags Items
// @Accept  json
// @Produce  json
// @Param id path string true "Shop ID"
// @Param category_id path string true "Category ID"
// @Success 200 {object} domain.Items
// @Router /items/shop/:id/category/:category_id [get]
func (h ItemsHandler) GetItemsByShopCategoryID(c *gin.Context) {
	xTraceId, _ := c.Get("X-Trace-ID")
	ctx := context.WithValue(c.Request.Context(), tracing.XtraceHeaderKey, xTraceId)

	shopID := c.Param("id")
	categoryID := c.Param("category_id")
	err := utils.ValidateHexID([]string{shopID, categoryID})
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	itemsResponse, err := h.Service.GetItemsByShopCategoryID(ctx, shopID, categoryID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	if len(itemsResponse.Items) == 0 {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(http.StatusOK, itemsResponse)
}

// GetItemsByIDs godoc
// @Summary Get items by ids
// @Description Get item by IDs in body
// @Tags Items
// @Accept  json
// @Produce  json
// @Param items body domain.ItemsIds true "Add items"
// @Success 200 {object} domain.Items
// @Router /items/list [post]
func (h ItemsHandler) GetItemsByIDs(c *gin.Context) {
	var input models.ItemsIds

	xTraceId, _ := c.Get("X-Trace-ID")
	ctx := context.WithValue(c.Request.Context(), tracing.XtraceHeaderKey, xTraceId)

	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		apiErr := apierrors.NewGenericErrorMessageDecoder(err)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	err := utils.ValidateHexID(input.Items)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	response, apiErr := h.Service.GetItemsByIDs(ctx, input)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	if len(input.Items) != len(response.Items) {
		c.Header("integrity", "false")
	}

	c.JSON(http.StatusOK, response)
}

// CreateItem godoc
// @Summary Create Item
// @Description Create item in db
// @Tags Items
// @Accept  json
// @Produce  json
// @Param item body dto.ItemDTO true "Add item"
// @Success 201
// @Router /items [post]
func (h ItemsHandler) CreateItem(c *gin.Context) {
	var input dto.ItemDTO

	xTraceId, _ := c.Get("X-Trace-ID")
	ctx := context.WithValue(c.Request.Context(), tracing.XtraceHeaderKey, xTraceId)
	ctx = context.WithValue(ctx, goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		apiErr := apierrors.NewGenericErrorMessageDecoder(err)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	userID, err1 := goauth.GetUserId(c)
	if err1 != nil {
		apiErr := apierrors.NewUnauthorizedApiError(err1.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	input.UserID = userID

	//validate if the category exists
	catcheck, err := h.CategoriesService.Get(c, input.Category.ID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	if catcheck.Name != input.Category.Name {
		c.JSON(http.StatusBadRequest, apierrors.NewApiError("category name does not match with the existing cat for "+catcheck.ID, "bad_request", http.StatusBadRequest, apierrors.CauseList{}))
		return
	}

	response, apiErr := h.Service.CreateItem(ctx, input)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateItem godoc
// @Summary Update item
// @Description Update item by ID
// @Tags Items
// @Accept  json
// @Produce  json
// @Param id path string true "Item ID"
// @Param item body dto.ItemDTO true "Add item"
// @Success 200
// @Router /items/{id} [put]
func (h ItemsHandler) UpdateItem(c *gin.Context) {
	var input dto.ItemDTO

	xTraceId, _ := c.Get("X-Trace-ID")
	ctx := context.WithValue(c.Request.Context(), tracing.XtraceHeaderKey, xTraceId)

	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		apiErr := apierrors.NewGenericErrorMessageDecoder(err)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	userID, err1 := goauth.GetUserId(c)
	if err1 != nil {
		apiErr := apierrors.NewUnauthorizedApiError(err1.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	input.UserID = userID

	itemID := c.Param("id")
	apierr := utils.ValidateHexID([]string{itemID, input.Category.ID})
	if apierr != nil {
		c.JSON(apierr.Status(), apierr)
		return
	}

	//validate if the category exists
	catcheck, err := h.CategoriesService.Get(c, input.Category.ID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	if catcheck.Name != input.Category.Name {
		c.JSON(http.StatusBadRequest, apierrors.NewApiError("category name does not match with the existing cat for "+catcheck.ID, "bad_request", http.StatusBadRequest, apierrors.CauseList{}))
		return
	}

	apiErr := h.Service.Update(ctx, itemID, input)
	if apiErr != nil {
		logger.Error("error update items by id ", apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteItem godoc
// @Summary Delete item
// @Description Delete item by ID
// @Tags Items
// @Accept  json
// @Produce  json
// @Param id path string true "Item ID"
// @Success 204
// @Router /items/{id} [delete]
func (h ItemsHandler) DeleteItem(c *gin.Context) {

	xTraceId, _ := c.Get("X-Trace-ID")

	ctx := context.WithValue(c.Request.Context(), tracing.XtraceHeaderKey, xTraceId)
	ctx = context.WithValue(ctx, goauth.FirebaseAuthHeader, c.GetHeader("Authorization"))

	itemID := c.Param("id")
	err := utils.ValidateHexID([]string{itemID})
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	err = h.Service.Delete(ctx, itemID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusNoContent)
}
