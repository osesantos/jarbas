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

		conversation, _ := _loadConversation(chat)

		title := conversation.Title
		messages := conversation.Messages

		// Save conversation to Gomind
		gomind.StoreChat(title, messages)
	}

	spinner.Success("All conversations synced with Gomind!")
}
