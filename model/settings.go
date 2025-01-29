package model

import (
	"bufio"
	"errors"
	"fmt"
	"jarbas-go/main/utils"
	"reflect"
	"strings"
)

const (
	ApiKey       = "ApiKey"
	Model        = "Model"
	Vendor       = "Vendor"
	SaveMessages = "SaveMessages"
)

type Settings struct {
	ApiKey       string `json:"api_key"`
	Model        string `json:"model"`
	Vendor       string `json:"vendor"`
	SaveMessages bool   `json:"save_messages"`
}

func GetSettings() (Settings, error) {
	apiKey, err := GetKey()
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
		ApiKey:       apiKey,
		Model:        model,
		Vendor:       vendor,
		SaveMessages: saveMessages,
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

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(ApiKey)) {
			apiKey := strings.TrimPrefix(line, GetJsonKey(ApiKey)+": ")
			return apiKey, nil
		}
	}

	return "", errors.New("api key not found")
}

func GetModel() (string, error) {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(Model)) {
			apiKey := strings.TrimPrefix(line, GetJsonKey(Model)+": ")
			return apiKey, nil
		}
	}

	return "", errors.New("model not found")
}

func GetSaveMessages() (bool, error) {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return false, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(SaveMessages)) {
			option := strings.TrimPrefix(line, GetJsonKey(SaveMessages)+": ")

			if option == "n" {
				return false, nil
			}

			if option == "y" {
				return true, nil
			}

			return false, fmt.Errorf("invalid option: %s", option)
		}
	}

	return false, errors.New("save messages not found")
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
			return strings.TrimPrefix(line, GetJsonKey(Vendor)+": "), nil
		}
	}

	// Return an error if the file is empty
	return "", errors.New("vendor not found")
}
