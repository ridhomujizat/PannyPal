package outgoing

import (
	"fmt"
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
	case enum.BotTypeBaileys:
		return s.handleWebhookEventBaileys(accountBot, payload)
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

func (s *Service) handleWebhookEventBaileys(accountBot *models.AccountBot, payload dto.PayloadOutgoing) (*dto.ResponseOutgoing, error) {
	var req interface{}

	switch payload.Type {
	case "text":
		req = payload.ToReqBaileysText(*accountBot)
	default:
		return nil, nil
	}

	fmt.Println("req", req)
	fmt.Println("accountBot", accountBot)
	headers := http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer " + accountBot.Key},
	}

	resp, err := helper.HTTPRequest(&helper.HTTPRequestPayload{
		Method: enum.POST,
		URL:    accountBot.BaseURL + "/api/message/send",
		Body:   req,
	},
		&helper.HTTPRequestConfig{
			Headers: headers,
			Ctx:     s.ctx,
		})

	if err != nil {
		return nil, err
	}

	Response, err := helper.JSONToStruct[dto.BaileysOutgoingResponse](resp.Data)
	if err != nil {
		return nil, err
	}

	if Response == nil {
		return nil, nil
	}

	fmt.Println("Response", Response)

	response := dto.ResponseOutgoing{
		Message: "test",
		Id:      Response.Data.Key.ID,
	}

	return &response, nil
}
