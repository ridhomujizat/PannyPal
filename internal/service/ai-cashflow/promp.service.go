package aicashflow

import (
	"encoding/json"
	"fmt"
)

func (s *Service) promptUserTransactionInput(input string) (string, error) {
	categoryList, err := s.rp.Category.GetAllCategories()
	if err != nil {
		return "", err
	}

	// Buat mapping sederhana hanya ID dan Name untuk hemat token
	simplifiedCategories := make([]map[string]interface{}, len(categoryList))
	for i, cat := range categoryList {
		simplifiedCategories[i] = map[string]interface{}{
			"id":   cat.ID,
			"name": cat.Name,
		}
	}

	categoriesJSON, err := json.Marshal(simplifiedCategories)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`You are a JSON-only transaction parser. Return ONLY valid JSON, no explanations.

INPUT: "%s"

CATEGORIES: %s

OUTPUT FORMAT:
{
  "req_payload": [{
	"type": "",
	"amount": 0,
	"category_id": 0,
	"description": ""
  }]
}

RULES:
- type: "EXPENSE" or "INCOME"
- amount: integer only
- category_id: match from categories list
- description: original item name
- Multiple transactions: array of objects in req_payload
- If no transaction found, return empty array: {"req_payload": []}
- Return ONLY the JSON object, nothing else`, input, string(categoriesJSON)), nil
}

func (s *Service) promptUserTransactionInputEdit(input string, existJson interface{}) (string, error) {
	categoryList, err := s.rp.Category.GetAllCategories()
	if err != nil {
		return "", err
	}

	// Buat mapping sederhana hanya ID dan Name untuk hemat token
	simplifiedCategories := make([]map[string]interface{}, len(categoryList))
	for i, cat := range categoryList {
		simplifiedCategories[i] = map[string]interface{}{
			"id":   cat.ID,
			"name": cat.Name,
		}
	}

	categoriesJSON, err := json.Marshal(simplifiedCategories)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`You are a JSON-only transaction parser. Return ONLY valid JSON, no explanations.

EXISTING DATA: %s

INPUT: "%s"

CATEGORIES: %s

OUTPUT FORMAT:
{
  "req_payload": [{
	"type": "",
	"amount": 0,
	"category_id": 0,
	"description": ""
  }]
}

RULES:
- MERGE strategy: Keep ALL existing transactions, only UPDATE matching ones
- Match transaction by description (case-insensitive partial match)
- If INPUT matches existing description: UPDATE that transaction
- If INPUT is new: ADD to array
- If INPUT doesn't mention existing transaction: KEEP it unchanged
- type: "EXPENSE" or "INCOME"
- amount: integer only
- category_id: match from categories list
- description: item name
- Return complete req_payload array with ALL transactions (existing + updated + new)
- Return ONLY the JSON object, nothing else`, existJson, input, string(categoriesJSON)), nil
}

// performOCROnImage performs OCR on image using Gemini vision
func (s *Service) performOCROnImage(base64Image string) (string, error) {
	categoryList, err := s.rp.Category.GetAllCategories()
	if err != nil {
		return "", fmt.Errorf("failed to get categories: %w", err)
	}

	// Buat mapping sederhana hanya ID dan Name untuk hemat token
	simplifiedCategories := make([]map[string]interface{}, len(categoryList))
	for i, cat := range categoryList {
		simplifiedCategories[i] = map[string]interface{}{
			"id":   cat.ID,
			"name": cat.Name,
		}
	}

	categoriesJSON, err := json.Marshal(simplifiedCategories)
	if err != nil {
		return "", err
	}

	// Create OCR prompt
	prompt := fmt.Sprintf(`You are an OCR and JSON transaction parser. Extract financial transactions from the image.

IMPORTANT: You MUST include the base64 image data in your analysis.

Categories: %s

OUTPUT FORMAT:
{
  "req_payload": [{
    "type": "",
    "amount": 0,
    "category_id": 0,
    "description": ""
  }]
}

RULES:
- Scan the image carefully for transaction data (receipts, invoices, bank statements, shopping lists, etc.)
- Extract each transaction found
- type: "EXPENSE" or "INCOME"
- amount: integer only
- category_id: match from categories list (best guess if not exact match)
- description: item name from the image
- Multiple transactions: array of objects in req_payload
- If no transaction found, return empty array: {"req_payload": []}
- Return ONLY the JSON object, nothing else`, categoriesJSON)

	// Call Gemini with vision capability
	aiResponse, err := s.ai.GeminiPromptWithImage(prompt, base64Image)
	if err != nil {
		return "", fmt.Errorf("failed to perform OCR: %w", err)
	}

	if aiResponse == nil {
		return "", fmt.Errorf("OCR response is empty")
	}

	return *aiResponse, nil
}
