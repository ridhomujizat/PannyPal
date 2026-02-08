package incoming

import (
	"pannypal/internal/common/enum"
	"strings"
)

func (s *Service) IsCashFlowFunction(payload string) bool {
	return strings.Contains(payload, string(enum.TagKeuangan))
}
