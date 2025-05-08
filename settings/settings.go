package settings

import (
	"bufio"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"jarbas-go/main/utils"

	"github.com/osesantos/resulto"
)

const (
	APIKey       = "APIKey"
	Model        = "Model"
	Vendor       = "Vendor"
	SaveMessages = "SaveMessages"
)

type Settings struct {
	APIKey       string `json:"api_key"`
	Model        string `json:"model"`
	Vendor       string `json:"vendor"`
	SaveMessages bool   `json:"save_messages"`
}

func GetSettings() Settings {
	return Settings{
		APIKey:       GetKey().Unwrap(),
		Model:        GetModel().Unwrap(),
		Vendor:       GetVendor().Unwrap(),
		SaveMessages: GetSaveMessages().Unwrap(),
	}
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

func GetKey() resulto.Result[string] {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return resulto.Failure[string](err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(APIKey)) {
			apiKey := strings.TrimPrefix(line, GetJsonKey(APIKey)+": ")
			return resulto.Success(apiKey)
		}
	}

	return resulto.Failure[string](errors.New("api key not found"))
}

func GetModel() resulto.Result[string] {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return resulto.Failure[string](err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(Model)) {
			apiKey := strings.TrimPrefix(line, GetJsonKey(Model)+": ")
			return resulto.Success(apiKey)
		}
	}

	return resulto.Failure[string](errors.New("model not found"))
}

func GetSaveMessages() resulto.Result[bool] {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return resulto.Failure[bool](err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(SaveMessages)) {
			option := strings.TrimPrefix(line, GetJsonKey(SaveMessages)+": ")

			if option == "n" {
				return resulto.Success(false)
			}

			if option == "y" {
				return resulto.Success(true)
			}

			return resulto.Failure[bool](fmt.Errorf("invalid option: %s", option))
		}
	}

	return resulto.Failure[bool](errors.New("save messages not found"))
}

func GetVendor() resulto.Result[string] {
	file, err := utils.GetSettingsFile()
	if err != nil {
		return resulto.Failure[string](err)
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, GetJsonKey(Vendor)) {
			return resulto.Success(strings.TrimPrefix(line, GetJsonKey(Vendor)+": "))
		}
	}

	// Return an error if the file is empty
	return resulto.Failure[string](errors.New("vendor not found"))
}
