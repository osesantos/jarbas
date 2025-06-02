package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"jarbas-go/main/actions"
	"jarbas-go/main/model"
	"jarbas-go/main/settings"
	"jarbas-go/main/utils"
	"jarbas-go/main/vendors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/glamour"
)

const (
	DefaultPrompt  = "\u001B[1;34manswer:\u001B[0m "
	TokenPrompt    = "\u001B[1;35mtoken:\u001B[0m "
	QuestionPrompt = "\u001B[1;32mquestion:\u001B[0m "
	Separator      = "\u001B[1;33m-------------------------\u001B[0m"
)

func Chat(settings settings.Settings, messages []map[string]any, isOldConversation bool) error {
	if messages == nil {
		messages = []map[string]any{}
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithEmoji(),
		glamour.WithWordWrap(0),
		glamour.WithTableWrap(true),
	)
	if err != nil {
		return fmt.Errorf("error initializing glamour renderer: %w", err)
	}

	fmt.Println("Welcome to the chat! Write 'exit' or press Ctrl-C to close the chat.")
	fmt.Println("Write 'token' to activate and deactivate token information.")
	fmt.Println("Write 'editor' to open an editor to write your question.")
	fmt.Println("Write 'role' to change the role of the assistant.")
	if isOldConversation {
		fmt.Println("")
		fmt.Println("\u001B[1;33mYou are continuing an old conversation. You can still change the role.\u001B[0m")
		fmt.Println("Write 'previous' to see the previous messages in the conversation.")
	}
	fmt.Println("")

	withToken := false
	role := model.CloudEngineer
	systemPrompt := ""
	for {
		input, err := _getInput()
		if err != nil || input == "exit" {
			break
		}

		if input == "previous" {
			if len(messages) != 0 {
				fmt.Println("\u001B[1;33mPrevious messages in the conversation:\u001B[0m")
				fmt.Println("")
				for i, msg := range messages {
					content, ok := msg["content"].(string)
					if !ok {
						content = "No content"
					}
					out, err := renderer.Render(content)
					if err != nil {
						fmt.Println("Error rendering message:", err)
					} else {
						fmt.Printf("\u001B[1;33mMessage %d:\u001B[0m %s\n", i+1, out)
						fmt.Println(Separator)
					}
				}
			} else {
				fmt.Println("\u001B[1;31mNo previous messages in the conversation.\u001B[0m")
			}

			continue
		}

		if input == "editor" {
			input, err = _getEditor()
			if err != nil {
				return err
			}
		}

		if input == "role" {
			fmt.Println("\u001B[1;32mcurrent role is: " + role + "\u001B[0m")
			role, err = _listRoles()
			if err != nil {
				return err
			}

			messages = append(messages, map[string]any{
				"role":    "system",
				"content": vendors.MapToSystemPrompt(role),
			})

			continue
		}

		if input == "token" {
			withToken = !withToken
			if withToken {
				fmt.Println("\u001B[1;32mtoken information active!\u001B[0m")
			} else {
				fmt.Println("\u001B[1;31mtoken information deactivated!\u001B[0m")
			}
		} else {
			systemPrompt = vendors.MapToSystemPrompt(role)
			answer := actions.ChatQuestion(messages, input, settings, systemPrompt)
			messages = answer.PreviousMessages
			if withToken {
				fmt.Println(TokenPrompt + answer.TotalToken + " " + DefaultPrompt)
				fmt.Println(Separator)
				out, err := renderer.Render(answer.LastMessage)
				if err != nil {
					fmt.Println("Error rendering message:", err)
				} else {
					fmt.Println(out)
				}
				fmt.Println(Separator)
			} else {
				fmt.Println(DefaultPrompt)
				fmt.Println(Separator)
				out, err := renderer.Render(answer.LastMessage)
				if err != nil {
					fmt.Println("Error rendering message:", err)
				} else {
					fmt.Println(out)
				}
				fmt.Println(Separator)
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

func ContinueChat(settings settings.Settings) error {
	file, err := _listConversations()
	if err != nil {
		return err
	}

	// Parse conversation string to []map[string]interface{}
	messages, err := _loadConversation(file)
	if err != nil {
		return err
	}

	err = Chat(settings, messages, true)
	if err != nil {
		return err
	}

	return nil
}

func _listRoles() (string, error) {
	roles := []string{
		model.AIEngineer,
		model.SoftwareEngineer,
		model.CloudEngineer,
		model.Writer,
	}

	role := ""
	prompt := &survey.Select{
		Message: "Select a role:",
		Options: roles,
	}
	err := survey.AskOne(prompt, &role)
	if err != nil {
		return "", err
	}

	return role, nil
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

			files = utils.AddDateTimeToFiles(files)

			return files
		},
	}
	err := survey.AskOne(prompt, &file)
	if err != nil {
		return "", err
	}

	cleanedFile := utils.CleanFileName(file)

	return cleanedFile, nil
}

func _loadConversation(file string) ([]map[string]any, error) {
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

func _parse(data []byte) ([]map[string]any, error) {
	var respData []map[string]any
	err := json.Unmarshal(data, &respData)
	if err != nil {
		return nil, err
	}
	return respData, nil
}
