package aicashflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pannypal/internal/common/enum"
	"pannypal/internal/common/models"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/service/ai-cashflow/dto"
	dtoOutgoingMessage "pannypal/internal/service/outgoing/dto"
	dtoTransaction "pannypal/internal/service/transaction/dto"
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

func (s *Service) PannyPalBotCashflow(payload dto.PayloadAICashflow) {

	//xample
	categoryID1 := 1
	categoryID2 := 2
	req := []dtoTransaction.CreateTransactionRequest{
		{
			Amount:      10000,
			CategoryID:  &categoryID1,
			Type:        "EXPENSE",
			Description: "Lunch at restaurant",
		},
		{
			Amount:      5000,
			CategoryID:  &categoryID2,
			Type:        "INCOME",
			Description: "Freelance project",
		},
	}
	messageBot := fmt.Sprintf("(CUMAN SAMPLE KATAKATA)Berikut adalah draft transaksi cashflow yang telah dibuat berdasarkan pesan Anda:\n\n%v\n\nSilakan tinjau dan simpan draft ini jika sudah sesuai.", req)

	OutgiingMessage := dtoOutgoingMessage.PayloadOutgoing{
		Message:        messageBot,
		ReplyToMessage: &payload.MessageId,
		Type:           "TEXT",
		AccountId:      payload.From,
		To:             payload.To,
	}

	outResponse, err := s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	if outResponse == nil {
		fmt.Println("No response from outgoing service")
		return
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return
	}
	rawMessage := json.RawMessage(reqBytes)

	modelMessageToReply := models.MessageToReply{
		MessageID:   outResponse.Id,
		FeatureType: enum.FeatureTypeAIcashflow,
		Messsage:    messageBot,
		Additional:  &rawMessage,
	}

	saveTODraft, err := s.rp.Bot.CreateMessageToReply(modelMessageToReply)
	if err != nil {
		fmt.Println("Error saving MessageToReply:", err)
		return
	}
	fmt.Println("MessageToReply saved:", saveTODraft)

}
