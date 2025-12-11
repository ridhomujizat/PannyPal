package outgoing

import (
	"net/http"
	"pannypal/internal/common/enum"
	"pannypal/internal/common/models"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/service/outgoing/dto"
)

func (s *Service) HandleWebhookEventWaha(payload dto.PayloadOutgoing) (*dto.ResponseOutgoing, error) {
	accountBot, err := s.rp.Bot.GetBotByAccountID(payload.AccountId)
	if err != nil {
		return nil, err
	}

	if accountBot == nil {
		return nil, nil
	}

	switch accountBot.BotType {
	case enum.BotTypeWaha:
		return s.handleWebhookEventWaha(accountBot, payload)
	default:
		return nil, nil
	}
}

func (s *Service) handleWebhookEventWaha(accountBot *models.AccountBot, payload dto.PayloadOutgoing) (*dto.ResponseOutgoing, error) {
	var req interface{}

	switch payload.Type {
	case "TEXT":
		req = payload.ToReqWahaText(*accountBot)
	default:
		return nil, nil
	}
	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"X-Api-Key":    []string{accountBot.Key},
	}

	resp, err := helper.HTTPRequest(&helper.HTTPRequestPayload{
		Method: enum.POST,
		URL:    accountBot.BaseURL + "/api/sendText",
		Body:   req,
	},
		&helper.HTTPRequestConfig{
			Headers: headers,
			Ctx:     s.ctx,
		})

	if err != nil {
		return nil, err
	}

	Response, err := helper.JSONToStruct[dto.ResponseOutgoingwaha](resp.Data)
	if err != nil {
		return nil, err
	}

	if Response == nil {
		return nil, nil
	}

	response := dto.ResponseOutgoing{
		Message: Response.Body,
		Id:      Response.Data.ID.ID,
	}

	return &response, nil
}
