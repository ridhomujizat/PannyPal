package dto

import (
	"pannypal/internal/common/models"
)

type PayloadOutgoing struct {
	Message        string  `json:"message"`
	ReplyToMessage *string `json:"reply_to_message,omitempty"`
	Type           string  `json:"type"`
	AccountId      string  `json:"account_id"`
	To             string  `json:"to"`
	Participant    *string `json:"participant,omitempty"`
}

type ReqWahaText struct {
	ChatID  string  `json:"chatId"`
	Text    string  `json:"text"`
	ReplyTo *string `json:"replyTo,omitempty"`
	Session string  `json:"session"`
}

func (p *PayloadOutgoing) ToReqWahaText(account models.AccountBot) *ReqWahaText {
	return &ReqWahaText{
		ChatID:  p.To,
		Text:    p.Message,
		ReplyTo: p.ReplyToMessage,
		Session: account.SessionID,
	}
}

type ResponseOutgoing struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

type ReqBaileysText struct {
	Session     string  `json:"sessionId"`
	Type        string  `json:"type"`
	ChatID      string  `json:"to"`
	Text        string  `json:"text"`
	ReplyTo     *string `json:"replyTo,omitempty"`
	Participant *string `json:"participant,omitempty"`
}

func (p *PayloadOutgoing) ToReqBaileysText(account models.AccountBot) *ReqBaileysText {
	return &ReqBaileysText{
		ChatID:      p.To,
		Text:        p.Message,
		ReplyTo:     p.ReplyToMessage,
		Session:     account.SessionID,
		Type:        p.Type,
		Participant: p.Participant,
	}
}
