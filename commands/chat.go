package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"jarbas-go/main/actions"
	"jarbas-go/main/model"
	"jarbas-go/main/utils"
	"jarbas-go/main/vendors"

	"github.com/AlecAivazis/survey/v2"
)

const (
	DefaultPrompt  = "\u001B[1;34manswer:\u001B[0m "
	TokenPrompt    = "\u001B[1;35mtoken:\u001B[0m "
	QuestionPrompt = "\u001B[1;32mquestion:\u001B[0m "
)

func Chat(settings model.Settings, messages []map[string]interface{}) error {
	if messages == nil {
		messages = []map[string]interface{}{}
	}

	fmt.Println("Welcome to the chat! Write 'exit' or press Ctrl-C to close the chat.")
	fmt.Println("Write 'token' to activate and deactivate token information.")
	fmt.Println("Write 'editor' to open an editor to write your question.")
	fmt.Println("Write 'previous' to see previous chat messages.")
	fmt.Println("")

	withToken := false
	for {
		input, err := _getInput()
		if err != nil || input == "exit" {
			break
		}

		if input == "editor" {
			input, err = _getEditor()
			if err != nil {
				return err
			}
		}

		if input == "token" {
			withToken = !withToken
			if withToken {
				fmt.Println("\u001B[1;32mtoken information active!\u001B[0m")
			} else {
				fmt.Println("\u001B[1;31mtoken information deactivated!\u001B[0m")
			}
		} else {
			answer, err := actions.ChatQuestion(messages, input, settings, vendors.SoftwareEngineer())
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

	if settings.SaveMessages {
		err := SaveConversation(messages)
		if err != nil {
			return err
		}
	}

	return nil
}

func ContinueChat(settings model.Settings) error {
	file, err := _listConversations()
	if err != nil {
		return err
	}

	// parse conversation string to []map[string]interface{}
	messages, err := _loadConversation(file)
	if err != nil {
		return err
	}

	err = Chat(settings, messages)
	if err != nil {
		return err
	}

	return nil
}

func _getInput() (string, error) {
	question := ""
	prompt := &survey.Input{
		Message: "question: ",
		Help: "Write 'exit' or press Ctrl-C to close the chat.\n" +
			"Write 'token' to activate and deactivate token information.\n" +
			"Write 'editor' to open an editor to write your question.",
	}
	err := survey.AskOne(prompt, &question)
	if err != nil {
		return "", err
	}

	return question, nil
}

func _getEditor() (string, error) {
	editor := ""
	prompt := &survey.Editor{
		Message: "editor:",
	}
	err := survey.AskOne(prompt, &editor)
	if err != nil {
		return "", err
	}

	return editor, nil
}

func _listConversations() (string, error) {
	file := ""
	prompt := &survey.Input{
		Message: "conversation to open:",
		Suggest: func(toComplete string) []string {
			dir := GetCacheDir()
			files, _ := filepath.Glob(dir + "/" + toComplete + "*")

			files = utils.OrderFilesByTime(files)

			fmt.Println("ordered files:\n", files)

			return files
		},
	}
	err := survey.AskOne(prompt, &file)
	if err != nil {
		return "", err
	}

	return file, nil
}

func _loadConversation(file string) ([]map[string]interface{}, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	messages, err := _parse(data)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func _parse(data []byte) ([]map[string]interface{}, error) {
	var respData []map[string]interface{}
	err := json.Unmarshal(data, &respData)
	if err != nil {
		return nil, err
	}
	return respData, nil
}
