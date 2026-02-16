package incoming

import (
	"fmt"
	"pannypal/internal/common/enum"
	"pannypal/internal/common/models"
	"strings"

	"gorm.io/gorm"
)

// DetectAction detects the user's action from their message
func DetectAction(msg string) string {
	msg = strings.ToLower(msg)

	saveKeywords := []string{"save"}
	cancelKeywords := []string{"cancel"}
	editKeywords := []string{"edit"}

	for _, w := range saveKeywords {
		if strings.Contains(msg, w) {
			return "save"
		}
	}

	for _, w := range cancelKeywords {
		if strings.Contains(msg, w) {
			return "cancel"
		}
	}

	for _, w := range editKeywords {
		if strings.Contains(msg, w) {
			return "edit"
		}
	}

	return "none"
}

// GetUser retrieves a user by phone number, or creates one if it doesn't exist
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

// validateOrCreateCategory validates if a category exists, or creates/returns a default one
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

// formatCurrency formats an amount as currency with dots
func (s *Service) formatCurrency(amount int) string {
	amountStr := fmt.Sprintf("%d", amount)
	n := len(amountStr)
	if n <= 3 {
		return amountStr
	}

	var result []rune
	for i, digit := range amountStr {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, digit)
	}
	return string(result)
}

func (s *Service) IsCashFlowFunction(payload string) bool {
	return strings.Contains(payload, string(enum.TagKeuangan))
}
