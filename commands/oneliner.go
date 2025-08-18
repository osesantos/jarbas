package commands

import (
	"jarbas-go/main/actions"
	"jarbas-go/main/prompts"
	"jarbas-go/main/settings"
)

func GetOneLiner(settings settings.Settings, userIntput string) (string, error) {
	prompt := prompts.GetOneLiner(userIntput)
	response := actions.SingleQuestion(prompt, settings, "")
	return response, nil
}
