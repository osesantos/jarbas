package actions

import (
	"errors"

	"jarbas-go/main/model"
	"jarbas-go/main/settings"
	"jarbas-go/main/vendors/anthropic"
	"jarbas-go/main/vendors/openai"
)

func SingleQuestion(input string, settings settings.Settings, system string) (string, error) {
	if settings.Vendor == model.OpenAI {
		return openai.DoSingleQuestion(input, settings)
	}

	if settings.Vendor == model.Anthropic {
		return anthropic.DoSingleQuestion(input, settings, system)
	}

	return "", errors.New("vendor not found")
}

func ChatQuestion(messages []map[string]any, question string, settings settings.Settings, system string) model.Answer {
	if settings.Vendor == model.OpenAI {
		return openai.DoChatQuestion(messages, question, settings).Unwrap()
	}

	if settings.Vendor == model.Anthropic {
		return anthropic.DoChatQuestion(messages, question, settings, system).Unwrap()
	}

	return model.Answer{}
}
