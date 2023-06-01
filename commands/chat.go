package commands

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"jarbas-go/main/actions"
)

const defaultPrompt = "\u001B[1;34manswer:\u001B[0m "
const tokenPrompt = "\u001B[1;35mtoken:\u001B[0m "

func Chat(apKey string, model string) error {
	var messages []map[string]interface{}
	var withToken = false
	for {
		input, err := getInput()
		if err != nil || input == "exit" {
			break
		}

		if input == "token" {
			withToken = !withToken
			if withToken {
				fmt.Println("\u001B[1;32mtoken information active!\u001B[0m")
			} else {
				fmt.Println("\u001B[1;31mtoken information deactivated!\u001B[0m")
			}
		} else {
			answer, err := actions.ChatQuestion(messages, input, apKey, model)
			if err != nil {
				return err
			}

			messages = answer.PreviousMessages
			if withToken {
				fmt.Println(tokenPrompt + answer.TotalToken + " " + defaultPrompt + answer.LastMessage)
			} else {
				fmt.Println(defaultPrompt + answer.LastMessage)
			}
		}
	}
	return nil
}

func getInput() (string, error) {
	question := ""
	prompt := &survey.Input{
		Message: "question: ",
		Help: "Write 'exit' or press Ctrl-C to close the chat.\n" +
			"Write 'token' to activate and deactivate token information.",
	}
	err := survey.AskOne(prompt, &question)
	if err != nil {
		return "", err
	}

	return question, nil
}
