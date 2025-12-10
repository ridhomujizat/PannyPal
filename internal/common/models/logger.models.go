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
