package prompts

func GetChatTitlePrompt(conversation string) string {
	return `
I need a title for a conversation.
Please provide me a short and descriptive title for the conversation.

rules:
	- The title should be concise and to the point.
	- Keep in mind that the response will be used directly as the file name, so your response should not contain any special characters or spaces and should be concise.
	- YOU MUST NOT RESPOND WITH ANYTHING OTHER THAN THE TITLE.
	- Fomat MUST ALWAYS BE: "Title_of_the_conversation"
	- With no additional text or explanation.
	- With no more than 5 words.
	- Use underscores instead of spaces.

The conversation is between the <input> and </input> tags.

<input>
` + conversation + `
</input>
`
}
