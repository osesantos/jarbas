package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jarbas-go/main/model"
	"net/http"
)

const version = "2023-06-01"

func DoSingleQuestion(input string, settings model.Settings) (string, error) {
	// Define the request body as a JSON object
	body := map[string]interface{}{
		"model": settings.Model,
		"messages": []map[string]string{
			{"role": "user", "content": input},
		},
	}

	respData, err := _doRequest(body, settings.ApiKey)
	if err != nil {
		return "", err
	}

	response, err := _validateResponse(respData["choices"])
	if err != nil {
		return response, err
	}

	text := respData["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	return text, nil
}

func DoChatQuestion(messages []map[string]interface{}, question string, settings model.Settings) (model.Answer, error) {
	return model.Answer{}, nil
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
