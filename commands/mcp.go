package commands

import (
	"jarbas-go/main/settings"
	"jarbas-go/main/vendors/gomind"
)

func McpQuery(settings settings.Settings, query string) string {
	response := gomind.DoSingleQuestion(query, settings)
	return response.Unwrap()
}
