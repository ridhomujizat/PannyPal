package models

import (
	"encoding/json"
	"pannypal/internal/common/enum"

	"gorm.io/gorm"
)

type MessageToReply struct {
	gorm.Model
	MessageID   string           `gorm:"type:varchar(100);uniqueIndex" json:"message_id"`
	FeatureType enum.FeatureType `gorm:"type:varchar(50)" json:"feature_type"`
	Messsage    string           `gorm:"type:text" json:"message"`
	Additional  *json.RawMessage `gorm:"type:jsonb" json:"additional"`
}

type AccountBot struct {
	gorm.Model
	AccountID uint   `gorm:"type:bigint;index" json:"account_id"`
	BotType   string `gorm:"type:varchar(50)" json:"bot_type"`
	BaseURL   string `gorm:"type:varchar(255)" json:"base_url"`
	Key       string `gorm:"type:varchar(255)" json:"key"`
}
