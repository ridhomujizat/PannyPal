package webhook

import (
	"pannypal/internal/common/enum"
	"pannypal/internal/service/webhook/dto"
	"strings"
)

func (s *Service) IsCashFlowFunction(payload string) bool {
	return strings.Contains(payload, string(enum.TagKeuangan))
}

func (s *Service) IsReplayMessage(message dto.Payloadwaha) (*enum.FeatureType, bool) {
	if message.Payload.Data.QuotedStanzaID == "" || message.Payload.ReplyTo == nil {
		return nil, false
	}
	check, err := s.rp.Bot.MessageToReplyMessage(message.Payload.Data.QuotedStanzaID)
	if err != nil {
		return nil, false
	}

	return &check.FeatureType, true
}
