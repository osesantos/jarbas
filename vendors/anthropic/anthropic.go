package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jarbas-go/main/model"
	"net/http"
)

const version = "2023-06-01"
const maxTokens = 1024

func DoSingleQuestion(input string, settings model.Settings) (string, error) {
	body := _getBody(settings.Model, []map[string]interface{}{{
		"role":    "user",
		"content": input,
	}})

	respData, err := _doRequest(body, settings.ApiKey)
	if err != nil {
		return "", err
	}

	response, err := _validateResponse(respData["content"])
	if err != nil {
		return response, err
	}

	// TODO: Implement a better way to handle the response, by using a struct and parsing the JSON
	messageType := respData["content"].([]interface{})[0].(map[string]interface{})["type"].(string)
	if messageType != "text" {
		return "", fmt.Errorf("response is not of type text")
	}

	// TODO: Implement a better way to handle the response, by using a struct and parsing the JSON
	text := respData["content"].([]interface{})[0].(map[string]interface{})["text"].(string)

	return text, nil
}

func DoChatQuestion(messages []map[string]interface{}, question string, settings model.Settings) (model.Answer, error) {
	newQuestion := map[string]interface{}{
		"role": "user", "content": question,
	}

	conversation := append(messages, newQuestion)

	body := _getBody(settings.Model, conversation)
	respData, err := _doRequest(body, settings.ApiKey)
	if err != nil {
		return model.Answer{}, err
	}

	// TODO: Implement a better way to handle the response, by using a struct and parsing the JSON
	responseType := respData["content"].([]interface{})[0].(map[string]interface{})["type"].(string)
	if responseType != "text" {
		return model.Answer{}, fmt.Errorf("response is not of type text")
	}

	// TODO: Implement a better way to handle the response, by using a struct and parsing the JSON
	response := respData["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	role := respData["role"].(string)
	// This is the message that will be sent as assistant response
	messageRequest := map[string]interface{}{
		"role":    role,
		"content": response,
	}

	// TODO: Implement a better way to handle the response, by using a struct and parsing the JSON
	inputTokens := respData["usage"].(map[string]interface{})["input_tokens"].(float64)
	outputTokens := respData["usage"].(map[string]interface{})["output_tokens"].(float64)
	totalTokens := inputTokens + outputTokens

	answer := model.Answer{
		PreviousMessages: append(conversation, messageRequest),
		LastMessage:      response,
		PromptToken:      fmt.Sprint(inputTokens),
		CompletionToken:  fmt.Sprint(outputTokens),
		TotalToken:       fmt.Sprint(totalTokens),
	}

	return answer, nil
}

func _getBody(model string, messages []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"model":      model,
		"max_tokens": maxTokens,
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
