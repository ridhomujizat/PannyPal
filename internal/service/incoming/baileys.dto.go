package incoming

import (
	"fmt"
	"net/http"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	dtoAI "pannypal/internal/service/ai/dto"
	"pannypal/internal/service/incoming/dto"
	dtoOutgoing "pannypal/internal/service/outgoing/dto"
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

	_, responseMessage, err := s.ai.InputTextCashflow(payload)
	if err != nil {
		Outgoing.Message = err.Error()
		_, err = s.outgoing.HandleWebhookEventWaha(Outgoing)
		if err != nil {
			return err
		}
		return err
	}

	Outgoing.Message = responseMessage

	_, err = s.outgoing.HandleWebhookEventWaha(Outgoing)
	if err != nil {
		return err
	}

	return nil
}
