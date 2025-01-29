package actions

import (
	"errors"
	"jarbas-go/main/model"
	"jarbas-go/main/vendors/openai"
)

func SingleQuestion(input string, settings model.Settings) (string, error) {
	if settings.Vendor == "openai" {
		return openai.DoSingleQuestion(input, settings)
	}

	return "", errors.New("Vendor not found")
}

func ChatQuestion(messages []map[string]interface{}, question string, settings model.Settings) (model.Answer, error) {
	if settings.Vendor == "openai" {
		return openai.DoChatQuestion(messages, question, settings)
	}

	return model.Answer{}, errors.New("Vendor not found")
}
