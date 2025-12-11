package webhook

import (
	"encoding/json"
	"fmt"
	"pannypal/internal/common/models"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/service/webhook/dto"
)

func (s *Service) LogWebhookEventWaha(payload interface{}) {
	data, _ := json.Marshal(payload)
	models := models.LogWaha{
		Message: json.RawMessage(data),
	}
	result, err := helper.JSONToStruct[dto.Payloadwaha](payload)
	if err != nil {
		fmt.Println("Error parsing payload:", err)
	}
	if result != nil {
		models.Type = &result.Payload.Data.Type
	}

	_, _ = s.rp.LogData.CreateLogWaha(models)
}
