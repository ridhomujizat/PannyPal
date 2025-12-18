package aicashflow

import (
	"encoding/json"
	"fmt"
	"pannypal/internal/common/models"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

// getTransactionSchema builds JSON Schema for structured output with dynamic categories from database
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

// promptUserTransactionInput generates prompt for text-based transaction input
func (s *Service) promptUserTransactionInput(input string) (string, error) {
	prompt := fmt.Sprintf(`Extract financial transactions from this input text.

INPUT: "%s"

Extract with JSON schema.`, input)

	return prompt, nil
}

// promptUserTransactionInputEdit generates prompt for editing existing transactions
func (s *Service) promptUserTransactionInputEdit(input string, existJson interface{}) (string, error) {
	existingData, err := json.Marshal(existJson)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`Merge and update transactions based on user input.

EXISTING DATA: %s

NEW INPUT: "%s"

RULES:
- Keep ALL existing transactions
- UPDATE matching transactions (match by description)
- ADD new transactions from input
- If input doesn't mention existing transaction, KEEP it unchanged

Extract with JSON schema.`, string(existingData), input)

	return prompt, nil
}

// performOCROnImage performs OCR on image using Gemini vision with structured output
func (s *Service) performOCROnImage(base64Image string) (string, error) {
	prompt := `Extract financial transactions from this image (receipt, invoice, bank statement, shopping list, etc).

Extract with JSON schema.`

	// Get transaction schema for structured output (with dynamic categories)
	schema, err := s.getTransactionSchema()
	if err != nil {
		return "", fmt.Errorf("failed to get transaction schema: %w", err)
	}

	// Call Gemini with vision capability and schema
	aiResponse, err := s.ai.GeminiPromptWithImageAndSchema(prompt, base64Image, schema)
	if err != nil {
		return "", fmt.Errorf("failed to perform OCR: %w", err)
	}

	if aiResponse == nil {
		return "", fmt.Errorf("OCR response is empty")
	}

	// Log the prompt activity
	s.logPromptActivity(prompt, aiResponse.Response, aiResponse.TokenUsed, aiResponse.ResponseTime)

	return aiResponse.Response, nil
}

// logPromptActivity saves the prompt activity to database
func (s *Service) logPromptActivity(prompt, response string, tokenUsed, responseTime int) {
	modelName := "gemini"
	logEntry := models.LogPrompt{
		ModelLLM:     &modelName,
		Prompt:       prompt,
		Response:     response,
		TokenUsed:    tokenUsed,
		ResponseTime: responseTime,
	}

	_, err := s.rp.LogData.CreateLogPrompt(logEntry)
	if err != nil {
		fmt.Println("Error logging prompt activity:", err)
	}
}
