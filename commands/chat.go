package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"jarbas-go/main/actions"
	"jarbas-go/main/model"
	"jarbas-go/main/prompts"
	"jarbas-go/main/settings"
	"jarbas-go/main/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/glamour"
)

const (
	DefaultPrompt  = "\u001B[1;34manswer:\u001B[0m "
	TokenPrompt    = "\u001B[1;35mtoken:\u001B[0m "
	QuestionPrompt = "\u001B[1;32mquestion:\u001B[0m "
	Separator      = "\u001B[1;33m-------------------------\u001B[0m"
	Saving         = "\u001B[1;32mSaving conversation...\u001B[0m"
)

func Chat(settings settings.Settings, messages []model.Message, isOldConversation bool) error {
	if messages == nil {
		messages = []model.Message{}
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
					content := msg.Content
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

			messages = append(messages, model.Message{
				Role:    model.System,
				Content: prompts.MapToSystemPrompt(role),
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
			systemPrompt = prompts.MapToSystemPrompt(role)
			memories := GetMemories().Unwrap()
			if len(memories) > 0 {
				systemPrompt = prompts.AddMemory(systemPrompt, memories) // Add memories to the system prompt
			}

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

	defer func() {
		if settings.SaveMessages {
			titlePrompt :=
				`
			I need a title for the conversation.
			Please provide a short and descriptive title for the conversation.
			Keep in mind that the response will be used directly as the file name, so your response should not contain any special characters or spaces and should be concise.
			YOU MUST NOT RESPOND WITH ANYTHING OTHER THAN THE TITLE.
			Fomat MUST ALWAYS BE: "Title of the conversation"
			With no additional text or explanation.
			With no more than 5 words.
			`

			title := actions.ChatQuestion(messages, titlePrompt, settings, systemPrompt).LastMessage

			fmt.Println(Saving)
			fmt.Println("Saving conversation with title:", title)

			err := SaveConversation(model.Chat{
				Messages: messages,
				Title:    title,
			})
			if err != nil {
				return
			}
		}
	}()

	return nil
}

func ContinueChat(settings settings.Settings) error {
	file, err := _listConversations()
	if err != nil {
		return err
	}

	// Parse conversation string to []map[string]interface{}
	conversation, err := _loadConversation(file)
	if err != nil {
		return err
	}

	err = Chat(settings, conversation.Messages, true)
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
		model.Pentester,
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

			// order files by time
			files = utils.OrderFilesByTime(files)

			// add the title to the files
			files = utils.AddTitleToFiles(files, _loadConversation)

			// add the date time to the files
			files = utils.AddDateTimeToFiles(files)

			return files
		},
	}
	err := survey.AskOne(prompt, &file)
	if err != nil {
		return "", err
	}

	// Remove the "title (date time)" part from the filename
	cleanedFile := utils.CleanFileName(file)

	return cleanedFile, nil
}

func _loadConversation(file string) (model.Chat, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return model.Chat{}, err
	}

	conversation, err := _parse(data)
	if err != nil {
		return model.Chat{}, err
	}

	return conversation, nil
}

func _parse(data []byte) (model.Chat, error) {
	var respData model.Chat
	err := json.Unmarshal(data, &respData)
	if err != nil {
		// lets try to parse it as a map first
		var respMap []model.Message
		err = json.Unmarshal(data, &respMap)
		if err != nil {
			return model.Chat{}, fmt.Errorf("error parsing conversation: %w", err)
		}

		return model.Chat{
			Messages: respMap,
		}, err
	}
	return respData, nil
}
