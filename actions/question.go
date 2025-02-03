package actions

import (
	"errors"

	"jarbas-go/main/model"
	"jarbas-go/main/vendors/anthropic"
	"jarbas-go/main/vendors/openai"
)

func SingleQuestion(input string, settings model.Settings, system string) (string, error) {
	if settings.Vendor == model.OpenAI {
		return openai.DoSingleQuestion(input, settings)
	}

	if settings.Vendor == model.Anthropic {
		return anthropic.DoSingleQuestion(input, settings, system)
	}

	return "", errors.New("vendor not found")
}

func ChatQuestion(messages []map[string]interface{}, question string, settings model.Settings, system string) (model.Answer, error) {
	if settings.Vendor == model.OpenAI {
		return openai.DoChatQuestion(messages, question, settings)
	}

	if settings.Vendor == model.Anthropic {
		return anthropic.DoChatQuestion(messages, question, settings, system)
	}

	return model.Answer{}, errors.New("vendor not found")
}
