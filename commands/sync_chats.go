package commands

import (
	"fmt"
	"jarbas-go/main/vendors/gomind"

	"github.com/pterm/pterm"
)

func SyncChats() {
	spinner, _ := pterm.DefaultSpinner.Start("Syncing conversations with Gomind...")

	chats, err := GetAllConversationFiles()
	if err != nil {
		fmt.Println("Error getting conversation files:", err)
		return
	}

	for _, chat := range chats {
		spinner.UpdateText("Syncing conversation: " + chat)

		conversation, err := _loadConversation(chat)
		if err != nil {
			fmt.Println("Error loading conversation:", err)
			continue
		}

		title := conversation.Title
		messages := conversation.Messages

		// Save conversation to Gomind
		result := gomind.StoreChat(title, messages)
		if result.IsErr() {
			fmt.Println("Error saving conversation to Gomind:", result.Err)
			continue
		}
	}

	spinner.Success("All conversations synced with Gomind!")
}
