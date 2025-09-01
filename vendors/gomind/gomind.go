package gomind

import (
	"bytes"
	"encoding/json"
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
