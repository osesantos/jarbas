package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"jarbas-go/main/model"
	"jarbas-go/main/settings"

	"github.com/osesantos/resulto"
)

func DoSingleQuestion(input string, settings settings.Settings) resulto.Result[string] {
	// Define the request body as a JSON object
	body := map[string]any{
		"model": settings.Model,
		"messages": []map[string]string{
			{"role": "user", "content": input},
		},
	}

	respData, err := _doRequest(body, settings.APIKey)
	if err != nil {
		return resulto.Failure[string](err)
	}

	_, err = _validateResponse(respData["choices"])
	if err != nil {
		return resulto.Failure[string](err)
	}

	// TODO: Implement a better way to handle the response, by using a struct and parsing the JSON
	text := respData["choices"].([]any)[0].(map[string]any)["message"].(map[string]any)["content"].(string)

	return resulto.Success(text)
}

func DoChatQuestion(messages []map[string]any, question string, settings settings.Settings) resulto.Result[model.Answer] {
	newQuestion := map[string]any{
		"role": "user", "content": question,
	}

	finalMessage := append(messages, newQuestion)

	// Define the request body as a JSON object
	body := map[string]any{
		"model":    settings.Model,
		"messages": finalMessage,
	}
	respData, err := _doRequest(body, settings.APIKey)
	if err != nil {
		return resulto.Failure[model.Answer](err)
	}

	// TODO: Implement a better way to handle the response, by using a struct and parsing the JSON
	lastMessage := respData["choices"].([]any)[0].(map[string]any)["message"].(map[string]any)["content"].(string)
	messagesRequest := respData["choices"].([]any)[0].(map[string]any)["message"].(map[string]any)

	// TODO: Implement a better way to handle the response, by using a struct and parsing the JSON
	promptTokens := respData["usage"].(map[string]any)["prompt_tokens"]
	completionTokens := respData["usage"].(map[string]any)["completion_tokens"]
	totalTokens := respData["usage"].(map[string]any)["total_tokens"]

	answer := model.Answer{
		PreviousMessages: append(finalMessage, messagesRequest),
		LastMessage:      lastMessage,
		PromptToken:      fmt.Sprint(promptTokens),
		CompletionToken:  fmt.Sprint(completionTokens),
		TotalToken:       fmt.Sprint(totalTokens),
	}

	return resulto.Success(answer)
}

func _doRequest(body map[string]any, apiKey string) (map[string]any, error) {
	// Convert the request body to JSON
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create an HTTP request with the necessary headers
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

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

func _validateResponse(response interface{}) (string, error) {
	if response == nil {
		return "", fmt.Errorf("response is nil")
	}

	if response == "" {
		return "", fmt.Errorf("response is empty")
	}

	if response.([]interface{})[0] == nil {
		return response.(string), fmt.Errorf("response is not the expected type")
	}

	return fmt.Sprint(response), nil
}
