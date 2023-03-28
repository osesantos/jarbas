package commands

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"jarbas-go/main/actions"
)

func Chat(apKey string) error {
	var messages []map[string]interface{}
	for {
		input, err := getInput()
		if err != nil || input == "exit" {
			break
		}

		respMessages, answer, err := actions.ChatQuestion(messages, input, apKey)
		if err != nil {
			return err
		}

		messages = respMessages
		fmt.Println("\033[1;34manswer:\033[0m " + answer)
	}
	return nil
}

func getInput() (string, error) {
	question := ""
	prompt := &survey.Input{
		Message: "question: ",
		Help:    "To end the chat write 'exit' or press Ctrl-C",
	}
	err := survey.AskOne(prompt, &question)
	if err != nil {
		return "", err
	}

	return question, nil
}
