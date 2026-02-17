package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"pannypal/internal/common/models"
	ai "pannypal/internal/pkg/ai-connector"
	"pannypal/internal/repository/analytics"
	"pannypal/internal/service/chatbot/dto"
	"pannypal/internal/service/chatbot/prompts"
	"regexp"
	"strings"
	"time"
)

// IntentType represents the type of user intent
type IntentType string

const (
	IntentStatistical    IntentType = "statistical"    // Total, average, sum queries
	IntentTrend          IntentType = "trend"          // Trend analysis, naik/turun
	IntentCategory       IntentType = "category"       // Category breakdown
	IntentBudget         IntentType = "budget"         // Budget vs actual
	IntentRecommendation IntentType = "recommendation" // Recommendations
	IntentPrediction     IntentType = "prediction"     // Predictions
	IntentComparison     IntentType = "comparison"     // Period comparisons
	IntentSummary        IntentType = "summary"        // General summary
	IntentGeneral        IntentType = "general"        // General questions
)

// AnalysisEngine is the core engine for chatbot analysis
type AnalysisEngine struct {
	ctx           context.Context
	aiClient      *ai.AiClient
	dataFetcher   *DataFetcher
	visualizer    *Visualizer
	analyticsRepo analytics.IRepository
}

// NewAnalysisEngine creates a new AnalysisEngine instance
func NewAnalysisEngine(
	ctx context.Context,
	aiClient *ai.AiClient,
	analyticsRepo analytics.IRepository,
) *AnalysisEngine {
	return &AnalysisEngine{
		ctx:           ctx,
		aiClient:      aiClient,
		dataFetcher:   NewDataFetcher(analyticsRepo),
		visualizer:    NewVisualizer(),
		analyticsRepo: analyticsRepo,
	}
}

// AnalyzeQuery analyzes user query and generates response
func (e *AnalysisEngine) AnalyzeQuery(
	userQuery string,
	conversationHistory []models.ChatMessage,
) (*dto.MessageMetadata, string, int, int, error) {
	// 1. Detect intent
	intent := e.detectIntent(userQuery)

	// 2. Fetch relevant data based on intent
	dataContext, err := e.fetchRelevantData(intent, userQuery)
	if err != nil {
		return nil, "", 0, 0, fmt.Errorf("failed to fetch data: %w", err)
	}

	// 3. Build conversation history string
	historyStr := e.buildHistoryString(conversationHistory)

	// 4. Build prompt
	prompt := e.buildPrompt(intent, userQuery, dataContext, historyStr)

	// 5. Call AI
	result, err := e.aiClient.GeminiPrompt(prompt)
	if err != nil {
		return nil, "", 0, 0, fmt.Errorf("failed to call AI: %w", err)
	}

	// 6. Parse AI response
	metadata, textResponse := e.parseAIResponse(result.Response, dataContext)

	return metadata, textResponse, result.TokenUsed, result.ResponseTime, nil
}

// detectIntent detects user intent from query
func (e *AnalysisEngine) detectIntent(query string) IntentType {
	queryLower := strings.ToLower(query)

	// Statistical keywords
	if strings.Contains(queryLower, "total") ||
		strings.Contains(queryLower, "berapa") ||
		strings.Contains(queryLower, "jumlah") ||
		strings.Contains(queryLower, "rata-rata") {
		return IntentStatistical
	}

	// Trend keywords
	if strings.Contains(queryLower, "trend") ||
		strings.Contains(queryLower, "naik") ||
		strings.Contains(queryLower, "turun") ||
		strings.Contains(queryLower, "meningkat") ||
		strings.Contains(queryLower, "menurun") {
		return IntentTrend
	}

	// Category keywords
	if strings.Contains(queryLower, "kategori") ||
		strings.Contains(queryLower, "category") ||
		strings.Contains(queryLower, "percategory") ||
		strings.Contains(queryLower, "per kategori") ||
		strings.Contains(queryLower, "per category") ||
		strings.Contains(queryLower, "breakdown") ||
		(strings.Contains(queryLower, "chart") && strings.Contains(queryLower, "pengeluaran")) ||
		(strings.Contains(queryLower, "pie") && strings.Contains(queryLower, "pengeluaran")) {
		return IntentCategory
	}

	// Budget keywords
	if strings.Contains(queryLower, "budget") ||
		strings.Contains(queryLower, "anggaran") {
		return IntentBudget
	}

	// Recommendation keywords
	if strings.Contains(queryLower, "saran") ||
		strings.Contains(queryLower, "rekomendasi") ||
		strings.Contains(queryLower, "bagaimana") ||
		strings.Contains(queryLower, "cara") {
		return IntentRecommendation
	}

	// Prediction keywords
	if strings.Contains(queryLower, "prediksi") ||
		strings.Contains(queryLower, "akan") ||
		strings.Contains(queryLower, "cukup") {
		return IntentPrediction
	}

	// Comparison keywords
	if strings.Contains(queryLower, "bandingkan") ||
		strings.Contains(queryLower, "vs") ||
		strings.Contains(queryLower, "dibanding") {
		return IntentComparison
	}

	// Summary keywords
	if strings.Contains(queryLower, "ringkasan") ||
		strings.Contains(queryLower, "summary") ||
		strings.Contains(queryLower, "overview") {
		return IntentSummary
	}

	return IntentGeneral
}

// fetchRelevantData fetches data based on intent
func (e *AnalysisEngine) fetchRelevantData(intent IntentType, query string) (string, error) {
	now := time.Now()

	switch intent {
	case IntentStatistical, IntentSummary, IntentGeneral:
		// Fetch current month data by default
		return e.dataFetcher.FetchCurrentMonthData()

	case IntentTrend:
		// Fetch last 6 months trend
		return e.dataFetcher.FetchMonthlyTrend(now.Year())

	case IntentCategory:
		// Fetch category breakdown for current month
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		return e.dataFetcher.FetchAllCategoriesComparison(&startOfMonth, &endOfMonth)

	case IntentBudget:
		// Fetch current month data (budget comparison will be done by AI)
		return e.dataFetcher.FetchCurrentMonthData()

	case IntentRecommendation, IntentPrediction:
		// Fetch last 3 months for better recommendations
		return e.dataFetcher.FetchLastNMonths(3)

	case IntentComparison:
		// Fetch last 6 months for comparison
		return e.dataFetcher.FetchLastNMonths(6)

	default:
		// Default to current month
		return e.dataFetcher.FetchCurrentMonthData()
	}
}

// buildHistoryString builds conversation history string
func (e *AnalysisEngine) buildHistoryString(messages []models.ChatMessage) string {
	if len(messages) == 0 {
		return "Tidak ada history percakapan sebelumnya."
	}

	// Get last 5 messages
	start := 0
	if len(messages) > 5 {
		start = len(messages) - 5
	}

	historyParts := []string{}
	for _, msg := range messages[start:] {
		role := "User"
		if msg.Role == "assistant" {
			role = "Assistant"
		}
		historyParts = append(historyParts, fmt.Sprintf("%s: %s", role, msg.Content))
	}

	return strings.Join(historyParts, "\n")
}

// buildPrompt builds the complete prompt for AI
func (e *AnalysisEngine) buildPrompt(intent IntentType, userQuery, dataContext, history string) string {
	systemPrompt := prompts.SystemPrompt

	var mainPrompt string
	switch intent {
	case IntentRecommendation:
		mainPrompt = prompts.BuildRecommendationPrompt(dataContext)
	case IntentTrend:
		mainPrompt = prompts.BuildTrendAnalysisPrompt(dataContext)
	case IntentCategory:
		period := "bulan ini"
		mainPrompt = prompts.BuildCategoryAnalysisPrompt(dataContext, period)
	case IntentSummary:
		period := time.Now().Format("January 2006")
		mainPrompt = prompts.BuildSummaryPrompt(period, dataContext)
	default:
		mainPrompt = prompts.BuildAnalysisPrompt(userQuery, dataContext, history)
	}

	return fmt.Sprintf("%s\n\n%s", systemPrompt, mainPrompt)
}

// parseAIResponse parses AI response and extracts metadata
func (e *AnalysisEngine) parseAIResponse(aiResponse string, rawData string) (*dto.MessageMetadata, string) {
	metadata := &dto.MessageMetadata{}

	// Try to extract JSON from markdown code block (```json ... ```)
	jsonStr := extractJSONFromMarkdown(aiResponse)

	// If no markdown block found, try parsing the entire response as JSON
	if jsonStr == "" {
		jsonStr = aiResponse
	}

	// Try to parse the JSON
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonResponse); err != nil {
		// Not a valid JSON response, return as plain text
		return metadata, aiResponse
	}

	// Extract insights for metadata
	if insights, ok := jsonResponse["insights"].([]interface{}); ok {
		insightList := []string{}
		for _, insight := range insights {
			if insightStr, ok := insight.(string); ok {
				insightList = append(insightList, insightStr)
			}
		}
		if len(insightList) > 0 {
			metadata.Statistics = &dto.Statistics{}
		}
	}

	// Extract recommendations
	if recs, ok := jsonResponse["recommendations"].([]interface{}); ok {
		recList := []dto.RecommendationData{}
		for _, rec := range recs {
			if recMap, ok := rec.(map[string]interface{}); ok {
				recData := dto.RecommendationData{}
				if title, ok := recMap["title"].(string); ok {
					recData.Title = title
				}
				if desc, ok := recMap["description"].(string); ok {
					recData.Description = desc
				}
				if saving, ok := recMap["potential_saving"].(float64); ok {
					recData.PotentialSaving = saving
				}
				if diff, ok := recMap["difficulty"].(string); ok {
					recData.Difficulty = diff
				}
				recList = append(recList, recData)
			} else if recStr, ok := rec.(string); ok {
				// Simple string recommendation
				recList = append(recList, dto.RecommendationData{
					Title:      recStr,
					Difficulty: "medium",
				})
			}
		}
		metadata.Recommendations = recList
	}

	// Check if visualization is needed
	needsViz := false
	if needs, ok := jsonResponse["needs_visualization"].(bool); ok {
		needsViz = needs
	}

	if needsViz {
		vizType := "bar"
		if vt, ok := jsonResponse["visualization_type"].(string); ok {
			vizType = vt
		}

		// Generate visualization from raw data
		viz, err := e.visualizer.GenerateChartData(vizType, rawData)
		if err == nil {
			metadata.Visualization = viz
		}
	}

	// Return the FULL AI response as content (intro text + JSON block)
	// so the frontend can parse and display both the intro and structured data
	return metadata, aiResponse
}

// extractJSONFromMarkdown extracts JSON content from a markdown code block
// Handles patterns like: ```json\n{...}\n``` or ```\n{...}\n```
func extractJSONFromMarkdown(text string) string {
	// Pattern: ```json ... ``` or ``` ... ```
	re := regexp.MustCompile("(?s)```(?:json)?\\s*\\n(\\{.*?\\})\\s*\\n?```")
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
