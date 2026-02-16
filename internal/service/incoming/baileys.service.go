package incoming

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pannypal/internal/common/enum"
	"pannypal/internal/common/models"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	dtoAI "pannypal/internal/service/ai/dto"
	"pannypal/internal/service/incoming/dto"
	dtoOutgoing "pannypal/internal/service/outgoing/dto"
	"time"
)

func (s *Service) HandleWebhookEventBaileys(payload interface{}) *types.Response {

	incoming, err := helper.JSONToStruct[dto.BaileysIncomingMessage](payload)
	if err != nil {
		fmt.Println("Error parsing payload:", err)
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Error parsing payload",
			Data:    payload,
		})
	}

	if incoming.MessageType == "extendedTextMessage" {
		return s.HandleExtendedTextMessage(incoming)
	}

	text := incoming.GetText()

	if s.IsCashFlowFunction(text) {
		err := s.HandleCashFlowFunction(incoming)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusInternalServerError,
				Message: "Error processing cashflow function",
				Data:    err,
			})
		}
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "HandleWebhookEventBaileys success",
		Data:    payload,
	})
}

func (s *Service) HandleExtendedTextMessage(message *dto.BaileysIncomingMessage) *types.Response {
	id := message.Message.ExtendedTextMessage.ContextInfo.StanzaID
	fmt.Println("id", id)
	messageToReply, err := s.rp.Bot.MessageToReplyMessage(id)
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

	if messageToReply.FeatureType == enum.FeatureTypeAIcashflow {
		err := s.HandleCashFlowFunctionReplyAction(message, messageToReply)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusInternalServerError,
				Message: "Error processing cashflow reply action",
				Data:    err,
			})
		}
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusOK,
			Message: "HandleCashFlowFunctionReplyAction success",
			Data:    message,
		})
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "HandleExtendedTextMessage success",
		Data:    message,
	})
}

func (s *Service) HandleCashFlowFunction(message *dto.BaileysIncomingMessage) error {
	Outgoing := dtoOutgoing.PayloadOutgoing{
		To:             message.Key.RemoteJid,
		AccountId:      message.SessionID,
		Message:        "Ada yang error (BOT)",
		ReplyToMessage: &message.Key.ID,
		Type:           "text",
		Participant:    message.Key.Participant,
	}

	text := message.GetText()

	payload := dtoAI.InputTextCashflow{
		Message: text,
	}

	result, responseMessage, err := s.ai.InputTextCashflow(payload)
	if err != nil {
		Outgoing.Message = err.Error()
		_, err = s.outgoing.HandleWebhookEventWaha(Outgoing)
		if err != nil {
			return err
		}
		return err
	}

	Outgoing.Message = responseMessage
	reqBytes, err := json.Marshal(result.ReqPayload)
	if err != nil {
		Outgoing.Message = err.Error()
		_, err = s.outgoing.HandleWebhookEventWaha(Outgoing)
		if err != nil {
			return err
		}
		return err
	}
	rawMessage := json.RawMessage(reqBytes)

	messageOut, err := s.outgoing.HandleWebhookEventWaha(Outgoing)
	if err != nil {
		return err
	}
	modelMessageToReply := models.MessageToReply{
		MessageID:   messageOut.Id,
		FeatureType: enum.FeatureTypeAIcashflow,
		Messsage:    responseMessage,
		Additional:  &rawMessage,
		Participant: message.Key.Participant,
	}

	_, err = s.rp.Bot.CreateMessageToReply(modelMessageToReply)
	if err != nil {
		Outgoing.Message = err.Error()
		_, err = s.outgoing.HandleWebhookEventWaha(Outgoing)
		if err != nil {
			return err
		}
		return err
	}

	return nil
}

func (s *Service) HandleCashFlowFunctionReplyAction(message *dto.BaileysIncomingMessage, messageToReply *models.MessageToReply) error {
	Outgoing := dtoOutgoing.PayloadOutgoing{
		To:             message.Key.RemoteJid,
		AccountId:      message.SessionID,
		ReplyToMessage: &message.Key.ID,
		Type:           "text",
		Participant:    message.Key.Participant,
	}

	text := message.GetText()
	typeAction := DetectAction(text)

	switch typeAction {
	case "save":
		fmt.Println("Action detected: save")
		return s.SaveTransaction(message, messageToReply, Outgoing)
	case "cancel":
		fmt.Println("Action detected: cancel")
		return s.CancelTransaction(message, messageToReply, Outgoing)
	case "edit":
		fmt.Println("Action detected: edit")
		return s.EditTransaction(message, messageToReply, Outgoing)
	default:
		fmt.Println("No valid action detected")
		Outgoing.Message = "Maaf, saya tidak mengerti tindakan yang Anda maksud. Silakan balas dengan 'save', 'edit', atau 'cancel'."
		_, err := s.outgoing.HandleWebhookEventWaha(Outgoing)
		if err != nil {
			fmt.Println("Error sending clarification message:", err)
			return err
		}
		return nil
	}
}

func (s *Service) SaveTransaction(message *dto.BaileysIncomingMessage, messageToReply *models.MessageToReply, Outgoing dtoOutgoing.PayloadOutgoing) error {
	// Extract phone number from RemoteJid (format: 6281234567890@s.whatsapp.net)
	phoneNumber := message.Key.RemoteJid

	user, err := s.GetUser(phoneNumber)
	if err != nil {
		fmt.Println("Failed to get or create user:", err)
		Outgoing.Message = "Maaf, terjadi kesalahan saat memproses pengguna."
		_, _ = s.outgoing.HandleWebhookEventWaha(Outgoing)
		return err
	}

	dataTransaction, err := helper.JSONToStruct[[]dtoAI.TransactionPayload](messageToReply.Additional)
	if err != nil {
		fmt.Println("Error converting JSON to struct:", err)
		Outgoing.Message = "Maaf, terjadi kesalahan saat memproses data transaksi."
		_, _ = s.outgoing.HandleWebhookEventWaha(Outgoing)
		return err
	}
	if dataTransaction == nil {
		fmt.Println("No transaction data found in Additional field")
		Outgoing.Message = "Maaf, data transaksi tidak ditemukan."
		_, _ = s.outgoing.HandleWebhookEventWaha(Outgoing)
		return fmt.Errorf("no transaction data found")
	}

	for _, tx := range *dataTransaction {
		validCategoryID, err := s.validateOrCreateCategory(tx.CategoryId)
		if err != nil {
			fmt.Println("Failed to validate category:", err)
			Outgoing.Message = "Maaf, terjadi kesalahan saat memvalidasi kategori."
			_, _ = s.outgoing.HandleWebhookEventWaha(Outgoing)
			return err
		}

		model := models.Transaction{
			UserID:          user.ID,
			Type:            models.TransactionType(tx.Type),
			Amount:          float64(tx.Amount),
			CategoryID:      validCategoryID,
			Description:     tx.Description,
			TransactionDate: time.Now(),
		}
		_, err = s.rp.Transaction.CreateTransaction(model)
		if err != nil {
			fmt.Println("Failed to create transaction:", err)
			Outgoing.Message = "Maaf, terjadi kesalahan saat menyimpan transaksi."
			_, _ = s.outgoing.HandleWebhookEventWaha(Outgoing)
			return err
		}
	}

	Outgoing.Message = "Transaksi berhasil disimpan."
	_, err = s.outgoing.HandleWebhookEventWaha(Outgoing)
	if err != nil {
		fmt.Println("Error sending confirmation message:", err)
		return err
	}

	// Delete the MessageToReply after saving transactions
	err = s.rp.Bot.DeleteMessageToReply(messageToReply.MessageID)
	if err != nil {
		fmt.Println("Error deleting MessageToReply:", err)
		return err
	}

	return nil
}

func (s *Service) EditTransaction(message *dto.BaileysIncomingMessage, messageToReply *models.MessageToReply, Outgoing dtoOutgoing.PayloadOutgoing) error {
	text := message.GetText()

	payload := dtoAI.InputTextCashflow{
		Message: text,
	}

	result, responseMessage, err := s.ai.InputTextCashflow(payload)
	if err != nil {
		Outgoing.Message = "Maaf, terjadi kesalahan saat memproses permintaan Anda."
		_, _ = s.outgoing.HandleWebhookEventWaha(Outgoing)
		return err
	}

	Outgoing.Message = responseMessage

	outResponse, err := s.outgoing.HandleWebhookEventWaha(Outgoing)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}
	if outResponse == nil {
		fmt.Println("No response from outgoing service")
		return fmt.Errorf("no response from outgoing service")
	}

	reqBytes, err := json.Marshal(result.ReqPayload)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return err
	}
	rawMessage := json.RawMessage(reqBytes)

	messageToReply.MessageID = outResponse.Id
	messageToReply.Messsage = responseMessage
	messageToReply.Additional = &rawMessage

	updatedMessageToReply, err := s.rp.Bot.UpdateMessageToReply(*messageToReply)
	if err != nil {
		fmt.Println("Error updating MessageToReply:", err)
		return err
	}
	fmt.Println("MessageToReply updated:", updatedMessageToReply)

	return nil
}

func (s *Service) CancelTransaction(message *dto.BaileysIncomingMessage, messageToReply *models.MessageToReply, Outgoing dtoOutgoing.PayloadOutgoing) error {
	err := s.rp.Bot.DeleteMessageToReply(messageToReply.MessageID)
	if err != nil {
		fmt.Println("Error deleting MessageToReply:", err)
		Outgoing.Message = "Maaf, terjadi kesalahan saat membatalkan draft."
		_, _ = s.outgoing.HandleWebhookEventWaha(Outgoing)
		return err
	}

	Outgoing.Message = "Draft transaksi telah dibatalkan."
	_, err = s.outgoing.HandleWebhookEventWaha(Outgoing)
	if err != nil {
		fmt.Println("Error sending cancellation message:", err)
		return err
	}

	return nil
}
