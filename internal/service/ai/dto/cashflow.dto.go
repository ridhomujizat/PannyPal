package dto

type InputTextCashflow struct {
	Message string `json:"message"`
}

type TransactionResponseAi struct {
	ReqPayload []TransactionPayload `json:"req_payload"`
}

type TransactionPayload struct {
	Type        string `json:"type"`
	Amount      int    `json:"amount"`
	CategoryId  int    `json:"category_id"`
	Description string `json:"description"`
}
