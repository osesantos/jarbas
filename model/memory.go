package model

// Memory represents a collection of messages that can be used to store conversation history or context.
type Memory struct {
	// Title that describes the memory
	Title string `json:"title"`
	// Messages that are part of the memory
	Messages []string `json:"messages"`
}
