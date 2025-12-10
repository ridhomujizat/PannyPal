package ai

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AiClient struct {
	ctx          context.Context
	geminiClient *genai.Client
	geminiModel  string
}
type Config struct {
	GeminiAPIKey string
	GeminiModel  string
}

func NewAiClient(ctx context.Context, cfg *Config) *AiClient {
	aiClient := &AiClient{
		ctx: ctx,
	}

	if cfg.GeminiAPIKey != "" && cfg.GeminiModel != "" {
		client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
		if err != nil {
			log.Fatal("Gagal membuat klien Gemini:", err)
		}
		// Remove defer client.Close() - client will be used later

		aiClient.geminiClient = client
		aiClient.geminiModel = cfg.GeminiModel
	}

	return aiClient
}

func (a *AiClient) GeminiPrompt(prompt string) (*string, error) {
	if a.geminiClient == nil {
		return nil, fmt.Errorf("Gemini client is not initialized")
	}

	model := a.geminiClient.GenerativeModel(a.geminiModel)
	resp, err := model.GenerateContent(a.ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("received empty or invalid response structure from Gemini")
	}
	part := resp.Candidates[0].Content.Parts[0]
	rawGeminiText := ""
	if textPart, ok := part.(genai.Text); ok {
		rawGeminiText = string(textPart)
	} else {
		return nil, fmt.Errorf("unexpected response part type")
	}
	return &rawGeminiText, nil
}

// Close properly closes the Gemini client
func (a *AiClient) Close() error {
	if a.geminiClient != nil {
		return a.geminiClient.Close()
	}
	return nil
}
