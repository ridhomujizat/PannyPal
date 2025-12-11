package dto

type ReplayActionSimpleRequest struct {
	Payload        PayloadAICashflow `json:"payload" binding:"required"`
	QuotedStanzaID string            `json:"quoted_stanza_id" binding:"required"`
}
