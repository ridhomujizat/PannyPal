package webhook

import (
	"net/http"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
)

func (s *Service) HandleWebhookEventWaha(payload interface{}) *types.Response {
	s.LogWebhookEventWaha(payload)

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "HandleWebhookEventWaha success",
		Data:    payload,
	})
}
