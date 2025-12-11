package webhook

import (
	"fmt"
	"net/http"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	dtoAiCashflow "pannypal/internal/service/ai-cashflow/dto"
	"pannypal/internal/service/webhook/dto"
)

func (s *Service) HandleWebhookEventWaha(payload interface{}) *types.Response {
	// note
	// quotedStanzaID bisa buat nampung reply to message id
	s.LogWebhookEventWaha(payload)
	message, err := helper.JSONToStruct[dto.Payloadwaha](payload)
	if err != nil {
		fmt.Println("Error parsing payload:", err)
	}

	if isCasflowFunction := s.IsCashFlowFunction(message.Payload.Body); isCasflowFunction {
		payload := dtoAiCashflow.PayloadAICashflow{
			TypeBot: "WAHA",
			// From is phone number
			From:    message.Payload.From,
			To:      message.Payload.To,
			Message: message.Payload.Body,
			// MessageId is waha message id
			MessageId: message.Payload.ID,
			Type:      message.Payload.Data.Type,
		}
		go s.aiCashFlowService.PannyPalBotCashflow(payload)
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "HandleWebhookEventWaha success",
		Data:    payload,
	})
}
