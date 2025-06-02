package vendors

import (
	"jarbas-go/main/model"
	"jarbas-go/main/settings"
)

type Vendor interface {
	DoSingleQuestion(input string, settings settings.Settings) (string, error)
	DoChatQuestion(messages []model.Message, question string, settings settings.Settings) (model.Answer, error)
}
