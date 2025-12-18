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
)

func (s *Service) InputTransaction(payload dto.InputTransaction) *types.Response {
	// user, err := s.GetUser(payload.PhoneNumber)
	// if err != nil {
	// 	return helper.ParseResponse(&types.Response{
	// 		Code:    http.StatusInternalServerError,
	// 		Message: "Failed to get or create user",
	// 		Error:   err,
	// 		Data:    nil,
	// 	})
	// }
	categoryID1 := 1
	categoryID2 := 2
	req := []dto.TransactionPayload{
		{
			Amount:      10000,
			CategoryId:  categoryID1,
			Type:        "EXPENSE",
			Description: "MAKANAN SLEBEW",
		},
		{
			Amount:      5000,
			CategoryId:  categoryID2,
			Type:        "INCOME",
			Description: "Freelance project",
		},
	}

	prompt, err := s.promptUserTransactionInputEdit(payload.Message, req)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate prompt",
			Error:   err,
			Data:    nil,
		})
	}

	schema, err := s.getTransactionSchema()
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get transaction schema",
			Error:   err,
			Data:    nil,
		})
	}
	aiResponse, err := s.ai.GeminiPromptWithSchema(prompt, schema)
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

	// Log the prompt activity
	s.logPromptActivity(prompt, aiResponse.Response, aiResponse.TokenUsed, aiResponse.ResponseTime)

	fmt.Println("AI Response:", aiResponse.Response)

	// Clean up AI response - remove backticks and markdown formatting
	cleanResponse := s.cleanAIResponse(aiResponse.Response)
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

	// if !payload.SaveAsDraft {
	// 	// Validate category exists
	// 	validCategoryID, err := s.validateOrCreateCategory(result.ReqPayload.CategoryId)
	// 	if err != nil {
	// 		return helper.ParseResponse(&types.Response{
	// 			Code:    http.StatusInternalServerError,
	// 			Message: "Failed to validate category",
	// 			Error:   err,
	// 			Data:    nil,
	// 		})
	// 	}

	// 	model := models.Transaction{
	// 		UserID:      user.ID,
	// 		Type:        models.TransactionType(result.ReqPayload.Type),
	// 		Amount:      result.ReqPayload.Amount,
	// 		CategoryID:  validCategoryID,
	// 		Description: result.ReqPayload.Description,
	// 	}
	// 	_, err = s.rp.Transaction.CreateTransaction(model)
	// 	if err != nil {
	// 		return helper.ParseResponse(&types.Response{
	// 			Code:    http.StatusInternalServerError,
	// 			Message: "Failed to create transaction",
	// 			Error:   err,
	// 			Data:    nil,
	// 		})
	// 	}
	// 	return helper.ParseResponse(&types.Response{
	// 		Code:    http.StatusOK,
	// 		Message: "Transaction created successfully",
	// 		Data:    result,
	// 	})
	// }

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Transaction draft generated successfully",
		Data:    result,
	})

}

func (s *Service) PannyPalBotCashflow(payload dto.PayloadAICashflow) {
	// Implementation of PannyPalBotCashflow
	switch payload.Type {
	case "chat":
		s.PannyPalBotCashflowText(payload)
	case "image":
		s.PannyPalBotCashflowImage(payload)
	default:
		fmt.Println("Unsupported payload type:", payload.Type)
	}
}

func (s *Service) PannyPalBotCashflowText(payload dto.PayloadAICashflow) {
	OutgiingMessage := dtoOutgoingMessage.PayloadOutgoing{
		ReplyToMessage: &payload.MessageId,
		Type:           "TEXT",
		AccountId:      payload.From,
		To:             payload.To,
	}

	prompt, err := s.promptUserTransactionInput(payload.Message)
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses permintaan Anda."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error generating prompt:", err)
		return
	}

	schema, err := s.getTransactionSchema()
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses permintaan Anda."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error getting transaction schema:", err)
		return
	}
	aiResponse, err := s.ai.GeminiPromptWithSchema(prompt, schema)
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses permintaan Anda."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error getting AI response:", err)
		return
	}
	if aiResponse == nil {
		OutgiingMessage.Message = "Maaf, saya tidak dapat memahami permintaan Anda."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("AI response is empty")
		return
	}

	// Log the prompt activity
	s.logPromptActivity(prompt, aiResponse.Response, aiResponse.TokenUsed, aiResponse.ResponseTime)

	var result dto.TransactionResponseAi
	cleanResponse := s.cleanAIResponse(aiResponse.Response)
	err = json.Unmarshal([]byte(cleanResponse), &result)
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses data transaksi."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error unmarshaling AI response:", err)
		return
	}

	// Generate message from ReqPayload to save AI tokens
	messageResult := s.generateTransactionSummary(result.ReqPayload)
	messageResult += "\nBalas dengan _'save'_, _'edit'_, atau _'cancel'_."
	OutgiingMessage.Message = messageResult
	OutgiingMessage.ReplyToMessage = &payload.MessageId

	outResponse, err := s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	if outResponse == nil {
		fmt.Println("No response from outgoing service")
		return
	}

	reqBytes, err := json.Marshal(result.ReqPayload)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses data transaksi."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		return
	}
	rawMessage := json.RawMessage(reqBytes)

	modelMessageToReply := models.MessageToReply{
		MessageID:   outResponse.Id,
		FeatureType: enum.FeatureTypeAIcashflow,
		Messsage:    messageResult,
		Additional:  &rawMessage,
	}

	saveTODraft, err := s.rp.Bot.CreateMessageToReply(modelMessageToReply)
	if err != nil {
		fmt.Println("Error saving MessageToReply:", err)
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat menyimpan draft pesan."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		return
	}
	fmt.Println("MessageToReply saved:", saveTODraft)

}

func (s *Service) PannyPalBotCashflowImage(payload dto.PayloadAICashflow) {
	fmt.Println("Processing image payload for message ID:", payload.MessageId)
	OutgiingMessage := dtoOutgoingMessage.PayloadOutgoing{
		ReplyToMessage: &payload.MessageId,
		Type:           "TEXT",
		AccountId:      payload.From,
		To:             payload.To,
	}

	// Validate media payload
	if payload.Media == nil || payload.Media.URL == "" {
		OutgiingMessage.Message = "Maaf, tidak ada data gambar yang terdeteksi."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error: Media payload is missing")
		return
	}

	// Get account bot for authentication
	accountBot, err := s.rp.Bot.GetBotByAccountID(payload.From)
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses permintaan Anda."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error getting account bot:", err)
		return
	}

	if accountBot == nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses permintaan Anda."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error: Account bot not found")
		return
	}

	// Download image and encode to base64
	base64Image, err := s.downloadImageAndEncodeBase64(payload.Media.URL, accountBot)
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat mengunduh gambar."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error downloading image:", err)
		return
	}

	// Perform OCR on image
	ocrResponse, err := s.performOCROnImage(base64Image)
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses gambar."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error performing OCR:", err)
		return
	}

	// Parse OCR response
	var result dto.TransactionResponseAi
	cleanResponse := s.cleanAIResponse(ocrResponse)
	err = json.Unmarshal([]byte(cleanResponse), &result)
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses data transaksi dari gambar."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error unmarshaling OCR response:", err)
		return
	}

	// Check if any transactions were found
	if len(result.ReqPayload) == 0 {
		OutgiingMessage.Message = "Maaf, saya tidak menemukan data transaksi dari gambar yang Anda kirim."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		return
	}

	// Generate message from ReqPayload
	messageResult := s.generateTransactionSummary(result.ReqPayload)
	messageResult += "\n\nBalas dengan _'save'_, _'edit'_, atau _'cancel'_."
	OutgiingMessage.Message = messageResult
	OutgiingMessage.ReplyToMessage = &payload.MessageId

	outResponse, err := s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	if outResponse == nil {
		fmt.Println("No response from outgoing service")
		return
	}

	// Save draft message
	reqBytes, err := json.Marshal(result.ReqPayload)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses data transaksi."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		return
	}
	rawMessage := json.RawMessage(reqBytes)

	modelMessageToReply := models.MessageToReply{
		MessageID:   outResponse.Id,
		FeatureType: enum.FeatureTypeAIcashflow,
		Messsage:    messageResult,
		Additional:  &rawMessage,
	}

	saveTODraft, err := s.rp.Bot.CreateMessageToReply(modelMessageToReply)
	if err != nil {
		fmt.Println("Error saving MessageToReply:", err)
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat menyimpan draft pesan."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		return
	}
	fmt.Println("MessageToReply saved:", saveTODraft)
}

func (s *Service) PannyPalBotCashflowReplayAction(payload dto.PayloadAICashflow, messageToReply models.MessageToReply) {

	// Implementation of replay action
	typeAction := s.DetectAction(payload.Message)

	switch typeAction {
	case "save":
		fmt.Println("Action detected: save")
		s.SaveTransaction(payload, messageToReply)
	case "cancel":
		fmt.Println("Action detected: cancel")
		s.CancelTransaction(payload, messageToReply)
	case "edit":
		fmt.Println("Action detected: edit")
		s.EditTransaction(payload, messageToReply)
	default:
		fmt.Println("No valid action detected")
		outgoingMessage := dtoOutgoingMessage.PayloadOutgoing{
			Message:        "Maaf, saya tidak mengerti tindakan yang Anda maksud. Silakan balas dengan 'save', 'edit', atau 'cancel'.",
			ReplyToMessage: &payload.MessageId,
			Type:           "TEXT",
			AccountId:      payload.From,
			To:             payload.To,
		}
		_, err := s.outgoingService.HandleWebhookEventWaha(outgoingMessage)
		if err != nil {
			fmt.Println("Error sending clarification message:", err)
			return
		}
	}

}

func (s *Service) SaveTransaction(payload dto.PayloadAICashflow, messageToReply models.MessageToReply) {
	user, err := s.GetUser(payload.To)
	if err != nil {
		fmt.Println("Failed to get or create user:", err)
		return
	}

	dataTransaction, err := helper.JSONToStruct[[]dto.TransactionPayload](messageToReply.Additional)
	if err != nil {
		fmt.Println("Error converting JSON to struct:", err)
		return
	}
	if dataTransaction == nil {
		fmt.Println("No transaction data found in Additional field")
		return
	}
	for _, tx := range *dataTransaction {
		validCategoryID, err := s.validateOrCreateCategory(tx.CategoryId)
		if err != nil {
			fmt.Println("Failed to validate category:", err)
			return
		}

		model := models.Transaction{
			UserID:      user.ID,
			Type:        models.TransactionType(tx.Type),
			Amount:      tx.Amount,
			CategoryID:  validCategoryID,
			Description: tx.Description,
		}
		_, err = s.rp.Transaction.CreateTransaction(model)
		if err != nil {
			fmt.Println("Failed to create transaction:", err)
			return
		}

	}

	payloadBot := dtoOutgoingMessage.PayloadOutgoing{
		Message:        "Transaksi berhasil disimpan.",
		ReplyToMessage: &payload.MessageId,
		Type:           "TEXT",
		AccountId:      payload.From,
		To:             payload.To,
	}

	_, err = s.outgoingService.HandleWebhookEventWaha(payloadBot)
	if err != nil {
		fmt.Println("Error sending confirmation message:", err)
		return
	}

	// Delete the MessageToReply after saving transactions
	err = s.rp.Bot.DeleteMessageToReply(messageToReply.MessageID)
	if err != nil {
		fmt.Println("Error deleting MessageToReply:", err)
		return
	}
}

func (s *Service) EditTransaction(payload dto.PayloadAICashflow, messageToReply models.MessageToReply) {

	OutgiingMessage := dtoOutgoingMessage.PayloadOutgoing{
		ReplyToMessage: &payload.MessageId,
		Type:           "TEXT",
		AccountId:      payload.From,
		To:             payload.To,
	}

	prompt, err := s.promptUserTransactionInputEdit(payload.Message, messageToReply.Additional)
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses permintaan Anda."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		return
	}

	schema, err := s.getTransactionSchema()
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses permintaan Anda."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		fmt.Println("Error getting transaction schema:", err)
		return
	}
	aiResponse, err := s.ai.GeminiPromptWithSchema(prompt, schema)
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses permintaan Anda."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		return
	}
	if aiResponse == nil {
		OutgiingMessage.Message = "Maaf, saya tidak dapat memahami permintaan Anda."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		return
	}

	// Log the prompt activity
	s.logPromptActivity(prompt, aiResponse.Response, aiResponse.TokenUsed, aiResponse.ResponseTime)

	var result dto.TransactionResponseAi
	cleanResponse := s.cleanAIResponse(aiResponse.Response)
	err = json.Unmarshal([]byte(cleanResponse), &result)
	if err != nil {
		OutgiingMessage.Message = "Maaf, terjadi kesalahan saat memproses data transaksi."
		s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
		return
	}

	// Generate message from ReqPayload to save AI tokens
	messageBot := s.generateTransactionSummary(result.ReqPayload)
	messageBot += "\nBalas dengan _'save'_, _'edit'_, atau _'cancel'_."
	OutgiingMessage.Message = messageBot
	OutgiingMessage.ReplyToMessage = &payload.MessageId

	outResponse, err := s.outgoingService.HandleWebhookEventWaha(OutgiingMessage)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	if outResponse == nil {
		fmt.Println("No response from outgoing service")
		return
	}

	reqBytes, err := json.Marshal(result.ReqPayload)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return
	}
	rawMessage := json.RawMessage(reqBytes)

	messageToReply.MessageID = outResponse.Id
	messageToReply.Messsage = messageBot
	messageToReply.Additional = &rawMessage

	updatedMessageToReply, err := s.rp.Bot.UpdateMessageToReply(messageToReply)
	if err != nil {
		fmt.Println("Error updating MessageToReply:", err)
		return
	}
	fmt.Println("MessageToReply updated:", updatedMessageToReply)

}

func (s *Service) CancelTransaction(payload dto.PayloadAICashflow, messageToReply models.MessageToReply) {
	// Implementation of cancel transaction
	err := s.rp.Bot.DeleteMessageToReply(messageToReply.MessageID)
	if err != nil {
		fmt.Println("Error deleting MessageToReply:", err)
		return
	}

	payloadBot := dtoOutgoingMessage.PayloadOutgoing{
		Message:        "Draft transaksi telah dibatalkan.",
		ReplyToMessage: &payload.MessageId,
		Type:           "TEXT",
		AccountId:      payload.From,
		To:             payload.To,
	}

	_, err = s.outgoingService.HandleWebhookEventWaha(payloadBot)
	if err != nil {
		fmt.Println("Error sending cancellation message:", err)
		return
	}
}

func (s *Service) ReplayAction(payload dto.PayloadAICashflow, quotedStanzaID string) *types.Response {
	messageToReply, err := s.rp.Bot.MessageToReplyMessage(quotedStanzaID)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get MessageToReply",
			Error:   err,
			Data:    nil,
		})
	}
	if messageToReply == nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "No MessageToReply found for the given ID",
			Data:    nil,
		})
	}

	s.PannyPalBotCashflowReplayAction(payload, *messageToReply)

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Replay action processed successfully",
		Data:    nil,
	})
}
