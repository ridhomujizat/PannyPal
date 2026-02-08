package AI

import (
	"fmt"
	"pannypal/internal/common/models"
)

func (s *Service) logPromptActivity(prompt, response string, tokenUsed, responseTime int) {
	modelName := "gemini"
	logEntry := models.LogPrompt{
		ModelLLM:     &modelName,
		Prompt:       prompt,
		Response:     response,
		TokenUsed:    tokenUsed,
		ResponseTime: responseTime,
	}

	_, err := s.rp.LogData.CreateLogPrompt(logEntry)
	if err != nil {
		fmt.Println("Error logging prompt activity:", err)
	}
}
