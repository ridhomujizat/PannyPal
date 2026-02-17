package engine

import (
	"encoding/json"
	"pannypal/internal/service/chatbot/dto"
)

// Visualizer generates visualization data from raw data
type Visualizer struct{}

// NewVisualizer creates a new Visualizer instance
func NewVisualizer() *Visualizer {
	return &Visualizer{}
}

// GenerateChartData generates chart data from raw data
func (v *Visualizer) GenerateChartData(chartType string, rawData interface{}) (*dto.VisualizationData, error) {
	switch chartType {
	case "bar", "column":
		return v.generateBarChart(rawData)
	case "line":
		return v.generateLineChart(rawData)
	case "pie", "donut":
		return v.generatePieChart(rawData)
	case "table":
		return v.generateTable(rawData)
	default:
		return v.generateBarChart(rawData) // Default to bar chart
	}
}

// generateBarChart generates bar chart data
func (v *Visualizer) generateBarChart(rawData interface{}) (*dto.VisualizationData, error) {
	// Try to parse as category data
	dataMap, ok := rawData.(map[string]interface{})
	if !ok {
		// Try to unmarshal if it's JSON string
		if jsonStr, ok := rawData.(string); ok {
			if err := json.Unmarshal([]byte(jsonStr), &dataMap); err != nil {
				return nil, err
			}
		}
	}

	chartData := dto.ChartData{
		Labels: []string{},
		Values: []float64{},
		Colors: []string{},
	}

	// Extract categories if available
	if categories, ok := dataMap["categories"].([]interface{}); ok {
		for _, cat := range categories {
			catMap := cat.(map[string]interface{})
			if name, ok := catMap["category_name"].(string); ok {
				chartData.Labels = append(chartData.Labels, name)
			}
			if amount, ok := catMap["amount"].(float64); ok {
				chartData.Values = append(chartData.Values, amount)
			}
		}
	}

	// Extract months if available
	if months, ok := dataMap["months"].([]interface{}); ok {
		for _, month := range months {
			monthMap := month.(map[string]interface{})
			if monthNum, ok := monthMap["month"].(float64); ok {
				monthName := getMonthName(int(monthNum))
				chartData.Labels = append(chartData.Labels, monthName)
			}
			if expense, ok := monthMap["expense"].(float64); ok {
				chartData.Values = append(chartData.Values, expense)
			}
		}
	}

	return &dto.VisualizationData{
		Type: "bar",
		Data: chartData,
		Config: map[string]interface{}{
			"x_label": "Kategori/Bulan",
			"y_label": "Jumlah (Rp)",
			"format":  "currency",
		},
	}, nil
}

// generateLineChart generates line chart data
func (v *Visualizer) generateLineChart(rawData interface{}) (*dto.VisualizationData, error) {
	dataMap, ok := rawData.(map[string]interface{})
	if !ok {
		if jsonStr, ok := rawData.(string); ok {
			if err := json.Unmarshal([]byte(jsonStr), &dataMap); err != nil {
				return nil, err
			}
		}
	}

	chartData := dto.ChartData{
		Labels: []string{},
		Values: []float64{},
	}

	// Extract trend data
	if months, ok := dataMap["months"].([]interface{}); ok {
		for _, month := range months {
			monthMap := month.(map[string]interface{})
			if monthNum, ok := monthMap["month"].(float64); ok {
				monthName := getMonthName(int(monthNum))
				chartData.Labels = append(chartData.Labels, monthName)
			}
			// Use net value for line charts
			if net, ok := monthMap["net"].(float64); ok {
				chartData.Values = append(chartData.Values, net)
			} else if expense, ok := monthMap["expense"].(float64); ok {
				chartData.Values = append(chartData.Values, expense)
			}
		}
	}

	return &dto.VisualizationData{
		Type: "line",
		Data: chartData,
		Config: map[string]interface{}{
			"x_label": "Periode",
			"y_label": "Jumlah (Rp)",
			"format":  "currency",
		},
	}, nil
}

// generatePieChart generates pie chart data
func (v *Visualizer) generatePieChart(rawData interface{}) (*dto.VisualizationData, error) {
	dataMap, ok := rawData.(map[string]interface{})
	if !ok {
		if jsonStr, ok := rawData.(string); ok {
			if err := json.Unmarshal([]byte(jsonStr), &dataMap); err != nil {
				return nil, err
			}
		}
	}

	chartData := dto.ChartData{
		Labels: []string{},
		Values: []float64{},
		Colors: generateColors(10), // Generate up to 10 colors
	}

	// Extract categories for pie chart
	if categories, ok := dataMap["categories"].([]interface{}); ok {
		for _, cat := range categories {
			catMap := cat.(map[string]interface{})
			if name, ok := catMap["category_name"].(string); ok {
				chartData.Labels = append(chartData.Labels, name)
			}
			if amount, ok := catMap["amount"].(float64); ok {
				chartData.Values = append(chartData.Values, amount)
			}
		}
	}

	return &dto.VisualizationData{
		Type: "pie",
		Data: chartData,
		Config: map[string]interface{}{
			"format": "currency",
		},
	}, nil
}

// generateTable generates table data
func (v *Visualizer) generateTable(rawData interface{}) (*dto.VisualizationData, error) {
	// Table data is just the raw data formatted nicely
	return &dto.VisualizationData{
		Type: "table",
		Data: rawData,
		Config: map[string]interface{}{
			"format": "currency",
		},
	}, nil
}

// Helper function to get month name
func getMonthName(month int) string {
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	if month >= 1 && month <= 12 {
		return months[month-1]
	}
	return "Unknown"
}

// Helper function to generate colors for charts
func generateColors(count int) []string {
	colors := []string{
		"#FF6384", "#36A2EB", "#FFCE56", "#4BC0C0", "#9966FF",
		"#FF9F40", "#FF6384", "#C9CBCF", "#4BC0C0", "#FF6384",
	}
	if count > len(colors) {
		count = len(colors)
	}
	return colors[:count]
}

// ParseAndVisualize parses AI response and generates visualization if needed
func (v *Visualizer) ParseAndVisualize(aiResponse string) (*dto.VisualizationData, error) {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(aiResponse), &response); err != nil {
		// Not a JSON response, no visualization needed
		return nil, nil
	}

	needsViz, ok := response["needs_visualization"].(bool)
	if !ok || !needsViz {
		return nil, nil
	}

	vizType, ok := response["visualization_type"].(string)
	if !ok {
		vizType = "bar" // Default
	}

	// Get data from response or use entire response
	data := response
	if vizData, ok := response["data"]; ok {
		data = vizData.(map[string]interface{})
	}

	return v.GenerateChartData(vizType, data)
}
