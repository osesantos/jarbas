package model

import (
	"bufio"
	"fmt"
	"jarbas-go/main/utils"
	"reflect"
	"strings"
)

const (
	ApiKey         = "ApiKey"
	ApiURL         = "ApiURL"
	RequestHeaders = "RequestHeaders"
	Model          = "Model"
	Vendor         = "Vendor"
	SaveMessages   = "SaveMessages"
)

type Settings struct {
	ApiKey         string            `json:"api_key"`
	ApiURL         string            `json:"api_url"`
	RequestHeaders map[string]string `json:"request_headers"`
	Model          string            `json:"model"`
	Vendor         string            `json:"vendor"`
	SaveMessages   bool              `json:"save_messages"`
}

func GetSettings() (Settings, error) {
	apiKey, err := GetKey()
	if err != nil {
		return Settings{}, err
	}

	apiURL, err := GetAPIUrl()
	if err != nil {
		return Settings{}, err
	}

	requestHeaders, err := GetRequestHeaders()
	if err != nil {
		return Settings{}, err
	}

	model, err := GetModel()
	if err != nil {
		return Settings{}, err
	}

	vendor, err := GetVendor()
	if err != nil {
		return Settings{}, err
	}

	saveMessages, err := GetSaveMessages()
	if err != nil {
		return Settings{}, err
	}

	return Settings{
		ApiKey:         apiKey,
		ApiURL:         apiURL,
		RequestHeaders: requestHeaders,
		Model:          model,
		Vendor:         vendor,
		SaveMessages:   saveMessages,
	}, nil
}

func GetJsonKey(fieldName string) string {
	t := reflect.TypeOf(Settings{})
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return ""
	}

	jsonKey := field.Tag.Get("json")
	if jsonKey == "" {
		jsonKey = field.Name
	}
	return jsonKey
}

func GetKey() (string, error) {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return "", err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(ApiKey)) {
			// Remove the "api-key: " prefix from the line
			apiKey := strings.TrimPrefix(line, GetJsonKey(ApiKey))
			return apiKey, nil
		}
	}

	// Return an error if the file is empty
	return "", bufio.ErrFinalToken
}

func GetModel() (string, error) {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return "", err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(Model)) {
			// Remove the "api-key: " prefix from the line
			apiKey := strings.TrimPrefix(line, GetJsonKey(Model))
			return apiKey, nil
		}
	}

	// Return an error if the file is empty
	return "", bufio.ErrFinalToken
}

func GetSaveMessages() (bool, error) {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return false, err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(SaveMessages)) {
			option := strings.TrimPrefix(line, GetJsonKey(SaveMessages))

			if option == "n" {
				return false, nil
			}

			if option == "y" {
				return true, nil
			}

			return false, fmt.Errorf("invalid option: %s", option)
		}
	}

	// Return an error if the file is empty
	return false, bufio.ErrFinalToken
}

func GetAPIUrl() (string, error) {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return "", err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(ApiURL)) {
			return strings.TrimPrefix(line, GetJsonKey(ApiURL)), nil
		}
	}

	// Return an error if the file is empty
	return "", bufio.ErrFinalToken
}

func GetRequestHeaders() (map[string]string, error) {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	headers := map[string]string{}

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(RequestHeaders)) {
			headersLine := strings.TrimPrefix(line, GetJsonKey(RequestHeaders))
			headersList := strings.Split(headersLine, ",")
			for _, header := range headersList {
				headerParts := strings.Split(header, ":")
				headers[headerParts[0]] = headerParts[1]
			}
			return headers, nil
		}
	}

	// Return an error if the file is empty
	return nil, bufio.ErrFinalToken
}

func GetVendor() (string, error) {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return "", err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(Vendor)) {
			return strings.TrimPrefix(line, GetJsonKey(Vendor)), nil
		}
	}

	// Return an error if the file is empty
	return "", bufio.ErrFinalToken
}
