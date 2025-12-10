package aicashflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pannypal/internal/common/models"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/service/ai-cashflow/dto"
)

func (s *Service) InputTransaction(payload dto.InputTransaction) *types.Response {
	user, err := s.GetUser(payload.PhoneNumber)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get or create user",
			Error:   err,
			Data:    nil,
		})
	}

	prompt, err := s.promptUserTransactionInput(payload.Message)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate prompt",
			Error:   err,
			Data:    nil,
		})
	}

	aiResponse, err := s.ai.GeminiPrompt(prompt)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get AI response",
			Error:   err,
			Data:    nil,
		})
	}
	if aiResponse == nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "AI response is empty",
			Data:    nil,
		})
	}

	fmt.Println("AI Response:", *aiResponse)

	// Clean up AI response - remove backticks and markdown formatting
	cleanResponse := s.cleanAIResponse(*aiResponse)
	fmt.Println("Cleaned Response:", cleanResponse)

	var result dto.TransactionResponseAi

	err = json.Unmarshal([]byte(cleanResponse), &result)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse AI response",
			Error:   err,
			Data:    cleanResponse, // Return cleaned response for debugging
		})
	}

	if !payload.SaveAsDraft {
		// Validate category exists
		validCategoryID, err := s.validateOrCreateCategory(result.ReqPayload.CategoryId)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusInternalServerError,
				Message: "Failed to validate category",
				Error:   err,
				Data:    nil,
			})
		}

		model := models.Transaction{
			UserID:      user.ID,
			Type:        models.TransactionType(result.ReqPayload.Type),
			Amount:      result.ReqPayload.Amount,
			CategoryID:  validCategoryID,
			Description: result.ReqPayload.Description,
		}
		_, err = s.rp.Transaction.CreateTransaction(model)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create transaction",
				Error:   err,
				Data:    nil,
			})
		}
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusOK,
			Message: "Transaction created successfully",
			Data:    result,
		})
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Transaction draft generated successfully",
		Data:    result,
	})

}
