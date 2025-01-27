package commands

import (
	"fmt"
	"jarbas-go/main/actions"

	"github.com/AlecAivazis/survey/v2"
)

const DefaultPrompt = "\u001B[1;34manswer:\u001B[0m "
const TokenPrompt = "\u001B[1;35mtoken:\u001B[0m "
const QuestionPrompt = "\u001B[1;32mquestion:\u001B[0m "

func Chat(apKey string, model string, save_messages bool) error {
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
				fmt.Println(TokenPrompt + answer.TotalToken + " " + DefaultPrompt + answer.LastMessage)
			} else {
				fmt.Println(DefaultPrompt + answer.LastMessage)
			}
		}
	}

	if save_messages {
		err := actions.SaveMessages(messages)
		if err != nil {
			return err
		}
	}

	return nil
}

func getInput() (string, error) {
	question := ""
	prompt := &survey.Multiline{
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
