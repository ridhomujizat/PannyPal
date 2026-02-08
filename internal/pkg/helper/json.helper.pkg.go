package helper

import (
	"encoding/json"
	"fmt"
	"strings"
)

func JSONToString(payload any) (string, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	jsonString := string(jsonBytes)
	return jsonString, nil
}

func JSONToStruct[I any](payload any) (result *I, err error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func JSONToByte(payload any) ([]byte, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return jsonBytes, nil
}

// cleanAIResponse removes markdown formatting and backticks from AI response
func CleanAIResponse(response string) string {
	// Remove code block backticks
	response = strings.ReplaceAll(response, "```json", "")
	response = strings.ReplaceAll(response, "```", "")

	// Remove any leading/trailing whitespace
	response = strings.TrimSpace(response)

	// Find JSON start and end
	startIndex := strings.Index(response, "{")
	if startIndex == -1 {
		return response // No JSON found
	}

	// Find the last closing brace
	endIndex := strings.LastIndex(response, "}")
	if endIndex == -1 || endIndex < startIndex {
		return response // No valid JSON end found
	}

	// Extract only the JSON part
	return response[startIndex : endIndex+1]
}
