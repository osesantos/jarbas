package vendors

import (
	"jarbas-go/main/model"
	"jarbas-go/main/settings"
)

type Vendor interface {
	DoSingleQuestion(input string, settings settings.Settings) (string, error)
	DoChatQuestion(messages []map[string]any, question string, settings settings.Settings) (model.Answer, error)
}
