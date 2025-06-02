package model

type Conversation struct {
	Messages []map[string]any `json:"messages"`
	Title    string           `json:"title"`
}
