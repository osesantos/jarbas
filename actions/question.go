package actions

import (
	"jarbas-go/main/model"
	"jarbas-go/main/settings"
	"jarbas-go/main/vendors/anthropic"
	"jarbas-go/main/vendors/openai"
)

func SingleQuestion(input string, settings settings.Settings, system string) string {
	if settings.Vendor == model.OpenAI {
		return openai.DoSingleQuestion(input, settings).Unwrap()
	}

	if settings.Vendor == model.Anthropic {
		return anthropic.DoSingleQuestion(input, settings, system).Unwrap()
	}

	return ""
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
