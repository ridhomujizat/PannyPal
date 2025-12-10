package webhook

import (
	"encoding/json"
	"pannypal/internal/common/models"
)

func (s *Service) LogWebhookEventWaha(payload interface{}) {
	data, _ := json.Marshal(payload)
	models := models.LogWaha{
		Message: json.RawMessage(data),
	}

	_, _ = s.rp.LogData.CreateLogWaha(models)
}
