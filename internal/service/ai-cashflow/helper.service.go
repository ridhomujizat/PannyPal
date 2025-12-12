package aicashflow

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"pannypal/internal/common/models"
	"pannypal/internal/service/ai-cashflow/dto"
	"strconv"
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

func (s *Service) DetectAction(msg string) string {
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

func (s *Service) generateTransactionSummary(transactions []dto.TransactionPayload) string {
	if len(transactions) == 0 {
		return "Tidak ada transaksi yang terdeteksi."
	}

	summary := "*Summary:*\n\n"

	for _, tx := range transactions {
		// Get category name
		category, err := s.rp.Category.GetCategoryByID(uint(tx.CategoryId))
		categoryName := "Unknown"
		if err == nil && category != nil {
			categoryName = category.Name
		}

		// Format amount with dots
		amountStr := s.formatCurrency(int(tx.Amount))

		if tx.Type == "EXPENSE" {
			summary += "✅ Tercatat pengeluaran\n"
		} else {
			summary += "💰 Tercatat pemasukan\n"
		}

		summary += " n: " + tx.Description + "\n"
		summary += " a: Rp. " + amountStr + "\n"
		summary += " c: " + categoryName + "\n\n"
	}

	return summary
}

func (s *Service) formatCurrency(amount int) string {
	if amount < 1000 {
		return strconv.Itoa(amount)
	}

	// Convert to string
	str := strconv.Itoa(amount)

	// Add dots from right to left
	result := ""
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += "."
		}
		result += string(char)
	}

	return result
}

// downloadImageAndEncodeBase64 downloads image from media URL with authentication
func (s *Service) downloadImageAndEncodeBase64(mediaURL string, accountBot *models.AccountBot) (string, error) {
	// Replace localhost with base URL if needed
	if strings.Contains(mediaURL, "localhost") {
		mediaURL = strings.ReplaceAll(mediaURL, "http://localhost:3000", accountBot.BaseURL)
	}

	// Create HTTP request with authentication
	req, err := http.NewRequest("GET", mediaURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", accountBot.Key)
	req.Header.Set("Content-Type", "application/octet-stream")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	// Encode to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	return base64Image, nil
}
