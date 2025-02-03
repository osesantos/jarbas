package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"jarbas-go/main/model"
)

const (
	version   = "2023-06-01"
	maxTokens = 1024
)

func DoSingleQuestion(input string, settings model.Settings, system string) (string, error) {
	body := _getBody(settings.Model, []map[string]interface{}{{
		"role":    "user",
		"content": input,
	}}, system)

	respData, err := _doRequest(body, settings.APIKey)
	if err != nil {
		return "", err
	}

	response, err := ParseResponse(respData)
	if err != nil {
		return "", err
	}

	return response.Content, nil
}

func DoChatQuestion(messages []map[string]interface{}, question string, settings model.Settings, system string) (model.Answer, error) {
	newQuestion := map[string]interface{}{
		"role": "user", "content": question,
	}

	conversation := append(messages, newQuestion)

	body := _getBody(settings.Model, conversation, system)
	respData, err := _doRequest(body, settings.APIKey)
	if err != nil {
		return model.Answer{}, err
	}

	response, err := ParseResponse(respData)
	if err != nil {
		return model.Answer{}, err
	}

	answer := model.Answer{
		PreviousMessages: append(conversation, response.GetMessageRequest()),
		LastMessage:      response.Content,
		PromptToken:      fmt.Sprint(response.InputTokens),
		CompletionToken:  fmt.Sprint(response.OutputTokens),
		TotalToken:       fmt.Sprint(response.TotalTokens),
	}

	return answer, nil
}

func _getBody(model string, messages []map[string]interface{}, system string) map[string]interface{} {
	return map[string]interface{}{
		"model":      model,
		"max_tokens": maxTokens,
		"system":     system,
		"messages":   messages,
	}
}

func _doRequest(body map[string]interface{}, apiKey string) (map[string]interface{}, error) {
	// Convert the request body to JSON
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create an HTTP request with the necessary headers
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", version)

	// Send the HTTP request and read the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body as a string
	var respData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
