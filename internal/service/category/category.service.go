package category

import (
	"net/http"
	"pannypal/internal/common/models"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/service/category/dto"

	"gorm.io/gorm"
)

func (s *Service) CreateCategoryRequest(payload dto.CreateCategoryRequest) *types.Response {
	category := models.Category{
		Name: payload.Name,
	}

	createdCategory, err := s.rp.Category.CreateCategory(category)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create category",
			Data:    nil,
			Error:   err,
		})
	}

	response := dto.CategoryResponse{
		ID:        createdCategory.ID,
		Name:      createdCategory.Name,
		CreatedAt: createdCategory.CreatedAt,
		UpdatedAt: createdCategory.UpdatedAt,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusCreated,
		Message: "Category created successfully",
		Data:    response,
	})
}

func (s *Service) GetCategoriesRequest() *types.Response {
	categories, err := s.rp.Category.GetAllCategories()
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get categories",
			Data:    nil,
			Error:   err,
		})
	}

	categoryResponses := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		categoryResponses[i] = dto.CategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}

	response := dto.CategoryListResponse{
		Categories: categoryResponses,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Categories retrieved successfully",
		Data:    response,
	})
}

func (s *Service) GetCategoryByIDRequest(id uint) *types.Response {
	category, err := s.rp.Category.GetCategoryByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusNotFound,
				Message: "Category not found",
				Data:    nil,
				Error:   err,
			})
		}
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get category",
			Data:    nil,
			Error:   err,
		})
	}

	response := dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Category retrieved successfully",
		Data:    response,
	})
}

func (s *Service) UpdateCategoryRequest(id uint, payload dto.UpdateCategoryRequest) *types.Response {
	category, err := s.rp.Category.GetCategoryByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusNotFound,
				Message: "Category not found",
				Data:    nil,
				Error:   err,
			})
		}
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get category",
			Data:    nil,
			Error:   err,
		})
	}

	// Update fields if provided
	if payload.Name != nil {
		category.Name = *payload.Name
	}

	updatedCategory, err := s.rp.Category.UpdateCategory(*category)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update category",
			Data:    nil,
			Error:   err,
		})
	}

	response := dto.CategoryResponse{
		ID:        updatedCategory.ID,
		Name:      updatedCategory.Name,
		CreatedAt: updatedCategory.CreatedAt,
		UpdatedAt: updatedCategory.UpdatedAt,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Category updated successfully",
		Data:    response,
	})
}

func (s *Service) DeleteCategoryRequest(id uint) *types.Response {
	category, err := s.rp.Category.GetCategoryByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusNotFound,
				Message: "Category not found",
				Data:    nil,
				Error:   err,
			})
		}
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get category",
			Data:    nil,
			Error:   err,
		})
	}

	if err := s.rp.Category.DeleteCategory(*category); err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete category",
			Data:    nil,
			Error:   err,
		})
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Category deleted successfully",
		Data:    nil,
	})
}
