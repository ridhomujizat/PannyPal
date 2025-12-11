package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type LogWaha struct {
	gorm.Model
	Type    *string         `gorm:"type:varchar(50)" json:"type"`
	Message json.RawMessage `gorm:"type:jsonb" json:"message"`
}

type LogPrompt struct {
	gorm.Model
	ModelLLM *string `gorm:"type:varchar(100)" json:"model"`
	Prompt   string  `gorm:"type:text" json:"prompt"`
	Response string  `gorm:"type:text" json:"response"`
}

type LogWahaResponse struct {
	gorm.Model
	LogWahaID string          `gorm:"type:varchar(100);index" json:"log_waha_id"`
	IsSuccess bool            `gorm:"type:boolean" json:"is_success"`
	Response  json.RawMessage `gorm:"type:jsonb" json:"response"`
}
