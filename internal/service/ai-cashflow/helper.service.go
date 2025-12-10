package aicashflow

import (
	"pannypal/internal/common/models"
	"strings"

	"gorm.io/gorm"
)

func (s *Service) GetUser(phoneNumber string) (*models.User, error) {
	var user *models.User
	var err error

	if phoneNumber != "" {
		user, err = s.rp.User.GetUserByPhone(phoneNumber)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				createdUser, err := s.rp.User.CreateUser(models.User{
					PhoneNumber: phoneNumber,
				})
				if err != nil {
					return nil, err
				}
				user = createdUser
			} else {
				// Handle other database errors
				return nil, err
			}
		}
	} else {
		// If no phone number provided, create transaction without user
		return nil, nil
	}
	return user, nil
}

func (s *Service) validateOrCreateCategory(categoryID int) (*uint, error) {
	if categoryID <= 0 {
		// If no category or invalid ID, get the first available category or create default
		categories, err := s.rp.Category.GetAllCategories()
		if err != nil {
			return nil, err
		}

		if len(categories) == 0 {
			// Create default category if none exists
			defaultCategory, err := s.rp.Category.CreateCategory(models.Category{
				Name: "Other",
			})
			if err != nil {
				return nil, err
			}
			return &defaultCategory.ID, nil
		}

		// Use first available category
		return &categories[0].ID, nil
	}

	// Check if category exists
	_, err := s.rp.Category.GetCategoryByID(uint(categoryID))
	if err != nil {
		// Category doesn't exist, get first available category or create default
		categories, err := s.rp.Category.GetAllCategories()
		if err != nil {
			return nil, err
		}

		if len(categories) == 0 {
			// Create default category if none exists
			defaultCategory, err := s.rp.Category.CreateCategory(models.Category{
				Name: "Lainnya",
			})
			if err != nil {
				return nil, err
			}
			return &defaultCategory.ID, nil
		}

		// Use first available category
		return &categories[0].ID, nil
	}

	// Category exists, return as pointer
	validID := uint(categoryID)
	return &validID, nil
}

// cleanAIResponse removes markdown formatting and backticks from AI response
func (s *Service) cleanAIResponse(response string) string {
	// Remove code block backticks
	response = strings.ReplaceAll(response, "```json", "")
	response = strings.ReplaceAll(response, "```", "")

	// Remove any leading/trailing whitespace
	response = strings.TrimSpace(response)

	// Find JSON start and end
	startIndex := strings.Index(response, "{")
	if startIndex == -1 {
		return response // No JSON found
	}

	// Find the last closing brace
	endIndex := strings.LastIndex(response, "}")
	if endIndex == -1 || endIndex < startIndex {
		return response // No valid JSON end found
	}

	// Extract only the JSON part
	return response[startIndex : endIndex+1]
}
