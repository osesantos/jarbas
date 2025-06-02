package model

type Chat struct {
	Messages []Message `json:"messages"`
	Title    string    `json:"title"`
}
