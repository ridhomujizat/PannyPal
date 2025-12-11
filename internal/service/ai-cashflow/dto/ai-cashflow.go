package dto

import "pannypal/internal/common/enum"

type InputTransaction struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Message     string `json:"message" binding:"required"`
	SaveAsDraft bool   `json:"save_as_draft"`
}

type TransactionResponseAi struct {
	ReqPayload []TransactionPayload `json:"req_payload"`
}

type TransactionPayload struct {
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	CategoryId  int     `json:"category_id"`
	Description string  `json:"description"`
}
type PayloadAICashflow struct {
	TypeBot   enum.BotType `json:"type_bot"`
	Message   string       `json:"message"`
	Type      string       `json:"type"`
	MessageId string       `json:"message_id"`
	From      string       `json:"from"`
	To        string       `json:"to"`
}
