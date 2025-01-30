package agents

import (
	"jarbas-go/main/agents/summarizer"
	"jarbas-go/main/model"

	"github.com/AlecAivazis/survey/v2"
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

func RunAgent(agent string, settings model.Settings) error {
	if agent == "" {
		agent, _ = SelectAgent()
	}

	if agent == Summarizer.String() {
		options, err := summarizer.GetOptions()
		if err != nil {
			return err
		}

		return summarizer.Run(options, settings)
	}

	return nil
}

func SelectAgent() (string, error) {
	return _listAgents()
}

func _listAgents() (string, error) {
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
		return "", err
	}

	return agent, nil
}
