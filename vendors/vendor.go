package vendors

import "jarbas-go/main/model"

type Vendor interface {
	DoSingleQuestion(input string, settings model.Settings) (string, error)
	DoChatQuestion(messages []map[string]interface{}, question string, settings model.Settings) (model.Answer, error)
}
