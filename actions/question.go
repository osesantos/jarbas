package actions

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Question(input string, apiKey string, model string) (string, error) {
	// Define the request body as a JSON object
	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": input},
		},
	}
	respData, err := request(body, apiKey)
	if err != nil {
		return "", err
	}

	text := respData["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	return text, nil
}

func ChatQuestion(messages []map[string]interface{}, question string, apiKey string, model string) ([]map[string]interface{}, string, error) {
	newQuestion := map[string]interface{}{
		"role": "user", "content": question,
	}

	finalMessage := append(messages, newQuestion)

	// Define the request body as a JSON object
	body := map[string]interface{}{
		"model":    model,
		"messages": finalMessage,
	}
	respData, err := request(body, apiKey)
	if err != nil {
		return nil, "", err
	}

	lastMessage := respData["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	messagesRequest := respData["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})

	return append(finalMessage, messagesRequest), lastMessage, nil
}

func request(body map[string]interface{}, apiKey string) (map[string]interface{}, error) {
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
