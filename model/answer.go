package model

type Answer struct {
	PreviousMessages []Message `json:"previous_messages"`
	LastMessage      string    `json:"last_message"`
	PromptToken      string    `json:"prompt_token"`
	CompletionToken  string    `json:"completion_token"`
	TotalToken       string    `json:"total_token"`
}
