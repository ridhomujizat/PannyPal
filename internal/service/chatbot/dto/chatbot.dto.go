package dto

import "time"

// SendMessageRequest represents the request to send a message to chatbot
type SendMessageRequest struct {
	SessionID string `json:"session_id"` // Optional, create new if empty
	Message   string `json:"message" validate:"required"`
}

// ChatMessageResponse represents a chat message response
type ChatMessageResponse struct {
	SessionID    string      `json:"session_id"`
	Role         string      `json:"role"`
	Content      string      `json:"content"`
	Metadata     interface{} `json:"metadata,omitempty"` // Chart data, tables, etc
	TokenUsed    int         `json:"token_used"`
	ResponseTime int         `json:"response_time"`
	CreatedAt    time.Time   `json:"created_at"`
}

// ConversationResponse represents a conversation summary
type ConversationResponse struct {
	SessionID    string    `json:"session_id"`
	Title        string    `json:"title"`
	IsActive     bool      `json:"is_active"`
	MessageCount int       `json:"message_count"`
	LastMessage  time.Time `json:"last_message"`
}

// VisualizationData represents data for charts and visualizations
type VisualizationData struct {
	Type   string      `json:"type"`             // "bar", "line", "pie", "table"
	Data   interface{} `json:"data"`             // Actual chart data
	Config interface{} `json:"config,omitempty"` // Chart configuration hints
}

// ChartData represents generic chart data structure
type ChartData struct {
	Labels []string    `json:"labels"`
	Values []float64   `json:"values"`
	Colors []string    `json:"colors,omitempty"`
	Extra  interface{} `json:"extra,omitempty"`
}

// Statistics represents statistical analysis data
type Statistics struct {
	Total            float64 `json:"total"`
	Average          float64 `json:"average,omitempty"`
	ChangePercentage float64 `json:"change_percentage,omitempty"`
	TopCategory      string  `json:"top_category,omitempty"`
}

// RecommendationData represents a single recommendation
type RecommendationData struct {
	Title            string  `json:"title"`
	PotentialSaving  float64 `json:"potential_saving,omitempty"`
	Difficulty       string  `json:"difficulty"` // "easy", "medium", "hard"
	Description      string  `json:"description,omitempty"`
	ActionItems      []string `json:"action_items,omitempty"`
}

// MessageMetadata represents the metadata structure for chat messages
type MessageMetadata struct {
	Visualization        *VisualizationData    `json:"visualization,omitempty"`
	Statistics           *Statistics           `json:"statistics,omitempty"`
	Recommendations      []RecommendationData  `json:"recommendations,omitempty"`
	TotalPotentialSaving float64               `json:"total_potential_saving,omitempty"`
	RawData              interface{}           `json:"raw_data,omitempty"`
}
