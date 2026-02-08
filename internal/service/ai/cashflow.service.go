package AI

import (
	"encoding/json"
	"fmt"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/service/ai/dto"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

func (s *Service) promptUserTransactionInput(input string) (string, error) {
	prompt := fmt.Sprintf(`Extract financial transactions from this input text.

INPUT: "%s"

Extract with JSON schema.`, input)

	return prompt, nil
}

func (s *Service) getTransactionSchema() (*genai.Schema, error) {
	// Get categories from database
	categoryList, err := s.rp.Category.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Build category description from database
	var categoryDescParts []string
	for _, cat := range categoryList {
		categoryDescParts = append(categoryDescParts, fmt.Sprintf("%d=%s", cat.ID, cat.Name))
	}
	categoryDescription := "Category ID: " + strings.Join(categoryDescParts, ", ")

	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"req_payload": {
				Type:        genai.TypeArray,
				Description: "Array of transactions extracted from input",
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"type": {
							Type:        genai.TypeString,
							Description: "Transaction type",
							Enum:        []string{"EXPENSE", "INCOME"},
						},
						"amount": {
							Type:        genai.TypeInteger,
							Description: "Transaction amount in integer",
						},
						"category_id": {
							Type:        genai.TypeInteger,
							Description: categoryDescription,
						},
						"description": {
							Type:        genai.TypeString,
							Description: "Item or transaction description",
						},
					},
					Required: []string{"type", "amount", "category_id", "description"},
				},
			},
		},
		Required: []string{"req_payload"},
	}, nil
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
		amountStr := helper.FormatCurrency(int(tx.Amount))

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

func (s *Service) InputTextCashflow(payload dto.InputTextCashflow) (*dto.TransactionResponseAi, string, error) {
	prompt, err := s.promptUserTransactionInput(payload.Message)
	if err != nil {
		return nil, "", err
	}

	schema, err := s.getTransactionSchema()
	if err != nil {
		return nil, "", err
	}

	aiResponse, err := s.ai.GeminiPromptWithSchema(prompt, schema)
	if err != nil {
		return nil, "", err
	}

	if aiResponse == nil {
		return nil, "", fmt.Errorf("OCR response is empty")
	}

	s.logPromptActivity(prompt, aiResponse.Response, aiResponse.TokenUsed, aiResponse.ResponseTime)

	var result dto.TransactionResponseAi
	cleanResponse := helper.CleanAIResponse(aiResponse.Response)
	err = json.Unmarshal([]byte(cleanResponse), &result)
	if err != nil {
		return nil, "", err
	}

	// Generate message from ReqPayload to save AI tokens
	messageResult := s.generateTransactionSummary(result.ReqPayload)
	messageResult += "\nBalas dengan _'save'_, _'edit'_, atau _'cancel'_."

	return &result, messageResult, nil
}
