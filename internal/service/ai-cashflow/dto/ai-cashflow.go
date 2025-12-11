package dto

type InputTransaction struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Message     string `json:"message" binding:"required"`
	SaveAsDraft bool   `json:"save_as_draft"`
}

type TransactionResponseAi struct {
	Message    string `json:"message"`
	ReqPayload struct {
		Type        string  `json:"type"`
		Amount      float64 `json:"amount"`
		CategoryId  int     `json:"category_id"`
		Description string  `json:"description"`
	} `json:"req_payload"`
}

type PayloadAICashflow struct {
	TypeBot   string `json:"type_bot"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	MessageId string `json:"message_id"`
	From      string `json:"from"`
	To        string `json:"to"`
}
