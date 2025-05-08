package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"jarbas-go/main/model"
	"jarbas-go/main/settings"

	"github.com/osesantos/resulto"
)

const (
	version   = "2023-06-01"
	maxTokens = 1024
)

func DoSingleQuestion(input string, settings settings.Settings, system string) resulto.Result[string] {
	body := _getBody(settings.Model, []map[string]any{{
		"role":    "user",
		"content": input,
	}}, system)

	respData, err := _doRequest(body, settings.APIKey)
	if err != nil {
		return resulto.Failure[string](err)
	}

	response, err := ParseResponse(respData)
	if err != nil {
		return resulto.Failure[string](err)
	}

	return resulto.Success(response.Content)
}

func DoChatQuestion(messages []map[string]any, question string, settings settings.Settings, system string) resulto.Result[model.Answer] {
	newQuestion := map[string]any{
		"role": "user", "content": question,
	}

	conversation := append(messages, newQuestion)

	body := _getBody(settings.Model, conversation, system)
	respData, err := _doRequest(body, settings.APIKey)
	if err != nil {
		return resulto.Failure[model.Answer](err)
	}

	response, err := ParseResponse(respData)
	if err != nil {
		return resulto.Failure[model.Answer](err)
	}

	answer := model.Answer{
		PreviousMessages: append(conversation, response.GetMessageRequest()),
		LastMessage:      response.Content,
		PromptToken:      fmt.Sprint(response.InputTokens),
		CompletionToken:  fmt.Sprint(response.OutputTokens),
		TotalToken:       fmt.Sprint(response.TotalTokens),
	}

	return resulto.Success(answer)
}

func _getBody(model string, messages []map[string]any, system string) map[string]any {
	return map[string]any{
		"model":      model,
		"max_tokens": maxTokens,
		"system":     system,
		"messages":   messages,
	}
}

func _doRequest(body map[string]any, apiKey string) (map[string]any, error) {
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
	var respData map[string]any
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
