package gomind

import (
	"bytes"
	"encoding/json"
	"jarbas-go/main/model"
	"jarbas-go/main/settings"
	"net/http"

	"github.com/osesantos/resulto"
)

type GomindRequest struct {
	Query string `json:"query"`
}

type GomindResponse struct {
	Answer string `json:"answer"`
}

type ChatRequest struct {
	Title    string          `json:"title"`
	Messages []model.Message `json:"messages"`
}

func DoSingleQuestion(input string, settings settings.Settings) resulto.Result[string] {
	response, err := doMCPRequest(GomindRequest{Query: input})
	if err != nil {
		return resulto.Failure[string](err)
	}

	if response.Answer == "" {
		return resulto.Failure[string](nil)
	}

	return resulto.Success(response.Answer)
}

func doMCPRequest(body GomindRequest) (GomindResponse, error) {
	// Convert the request body to JSON
	jsonData, err := json.Marshal(body)
	if err != nil {
		return GomindResponse{}, err
	}

	// Create an HTTP request with the necessary headers
	req, err := http.NewRequest("POST", "http://gomind.home/mcp", bytes.NewBuffer(jsonData))
	if err != nil {
		return GomindResponse{}, err
	}

	req.Header.Set("content-type", "application/json")

	// Send the HTTP request and read the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GomindResponse{}, err
	}
	defer resp.Body.Close()

	// Read the response body as a string
	var respData GomindResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return GomindResponse{}, err
	}

	return respData, nil
}

func StoreChat(title string, messages []model.Message) resulto.Result[bool] {
	chatRequest := ChatRequest{
		Title:    title,
		Messages: messages,
	}

	// Convert the request body to JSON
	jsonData, err := json.Marshal(chatRequest)
	if err != nil {
		return resulto.Failure[bool](err)
	}

	// Create an HTTP request with the necessary headers
	req, err := http.NewRequest("POST", "http://gomind.home/store_chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return resulto.Failure[bool](err)
	}

	req.Header.Set("content-type", "application/json")

	// Send the HTTP request and read the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resulto.Failure[bool](err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resulto.Failure[bool](nil)
	}

	return resulto.Success(true)
}
