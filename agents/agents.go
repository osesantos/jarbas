package agents

import (
	"jarbas-go/main/agents/summarizer"
	"jarbas-go/main/settings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/osesantos/resulto"
)

type AgentType string

type Agents []AgentType

const (
	Summarizer AgentType = "summarizer"
)

var AgentTypes = Agents{
	Summarizer,
}

func (at Agents) ToStringArray() []string {
	var result []string
	for _, v := range at {
		result = append(result, string(v))
	}
	return result
}

func (a AgentType) String() string {
	return string(a)
}

func RunAgent(agent string, settings settings.Settings) {
	if agent == "" {
		agent = SelectAgent().Unwrap()
	}

	if agent == Summarizer.String() {
		summarizer.Run(summarizer.GetOptions(), settings)
	}
}

func SelectAgent() resulto.Result[string] {
	return _listAgents()
}

func _listAgents() resulto.Result[string] {
	agent := ""
	prompt := &survey.Input{
		Message: "agents:",
		Suggest: func(toComplete string) []string {
			a := AgentTypes
			return a.ToStringArray()
		},
	}
	err := survey.AskOne(prompt, &agent)
	if err != nil {
		return resulto.Failure[string](err)
	}

	return resulto.Success(agent)
}
