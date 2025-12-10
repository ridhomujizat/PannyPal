package category

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/pkg/rabbitmq"
	"pannypal/internal/pkg/validation"
	categoryService "pannypal/internal/service/category"
	"pannypal/internal/service/category/dto"
)

type Handler struct {
	ctx             context.Context
	rabbitmq        *rabbitmq.ConnectionManager
	categoryService categoryService.IService
}

type IHandler interface {
	NewRoutes(e *gin.RouterGroup)
	CreateCategory(c *gin.Context)
	GetCategories(c *gin.Context)
	GetCategoryByID(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

func NewHandler(ctx context.Context, rabbitmq *rabbitmq.ConnectionManager, categoryService categoryService.IService) IHandler {
	return &Handler{
		ctx:             ctx,
		rabbitmq:        rabbitmq,
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create new category
// @Description Create a new category for transactions
// @Tags Category APIs
// @Accept json
// @Produce json
// @Param category body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} dto.CategoryResponse "Category created successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /categories [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
			Error:   err,
		})
		return
	}

	if err := validation.Validate(payload); err != nil {
		send(helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    err.Error(),
			Error:   err,
		}))
		return
	}

	send(h.categoryService.CreateCategoryRequest(payload))
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get list of all categories
// @Tags Category APIs
// @Accept json
// @Produce json
// @Success 200 {object} dto.CategoryListResponse "Categories retrieved successfully"
// @Router /categories [get]
func (h *Handler) GetCategories(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	send(h.categoryService.GetCategoriesRequest())
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get specific category by ID
// @Tags Category APIs
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.CategoryResponse "Category retrieved successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Failure 404 {object} types.Response "Not Found"
// @Router /categories/{id} [get]
func (h *Handler) GetCategoryByID(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	id := c.Param("id")
	categoryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid category ID",
			Data:    nil,
			Error:   err,
		})
		return
	}

	send(h.categoryService.GetCategoryByIDRequest(uint(categoryID)))
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update existing category
// @Tags Category APIs
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body dto.UpdateCategoryRequest true "Updated category data"
// @Success 200 {object} dto.CategoryResponse "Category updated successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Failure 404 {object} types.Response "Not Found"
// @Router /categories/{id} [put]
func (h *Handler) UpdateCategory(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	id := c.Param("id")
	categoryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid category ID",
			Data:    nil,
			Error:   err,
		})
		return
	}

	var payload dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
			Error:   err,
		})
		return
	}

	if err := validation.Validate(payload); err != nil {
		send(helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    err.Error(),
			Error:   err,
		}))
		return
	}

	send(h.categoryService.UpdateCategoryRequest(uint(categoryID), payload))
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete category by ID
// @Tags Category APIs
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} types.Response "Category deleted successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Failure 404 {object} types.Response "Not Found"
// @Router /categories/{id} [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	id := c.Param("id")
	categoryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid category ID",
			Data:    nil,
			Error:   err,
		})
		return
	}

	send(h.categoryService.DeleteCategoryRequest(uint(categoryID)))
}
