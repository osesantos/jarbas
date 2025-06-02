package google

import (
	"context"
	"fmt"
	"jarbas-go/main/model"
	"jarbas-go/main/settings"

	"github.com/osesantos/resulto"
	"google.golang.org/genai"
)

func DoSingleQuestion(input string, settings settings.Settings) resulto.Result[string] {
	ctx := context.Background()
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  settings.APIKey,
		Backend: genai.BackendGeminiAPI,
	})

	var history []*genai.Content

	chat, _ := client.Chats.Create(ctx, settings.Model, nil, history)
	res, _ := chat.SendMessage(ctx, genai.Part{Text: input})

	if len(res.Candidates) == 0 {
		return resulto.Failure[string](fmt.Errorf("no response received"))
	}

	return resulto.Success(res.Candidates[0].Content.Parts[0].Text)
}

func DoChatQuestion(messages []model.Message, question string, settings settings.Settings) resulto.Result[model.Answer] {
	ctx := context.Background()
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  settings.APIKey,
		Backend: genai.BackendGeminiAPI,
	})

	var history []*genai.Content
	for _, msg := range messages {
		history = append(history, genai.NewContentFromText(msg.Content, genai.Role(msg.Role)))
	}

	chat, _ := client.Chats.Create(ctx, settings.Model, nil, history)
	res, _ := chat.SendMessage(ctx, genai.Part{Text: question})

	if len(res.Candidates) == 0 {
		return resulto.Failure[model.Answer](fmt.Errorf("no response received"))
	}

	answer := model.Answer{
		PreviousMessages: append(messages, model.Message{
			Role:    model.User,
			Content: question,
		}),
		LastMessage:     res.Candidates[0].Content.Parts[0].Text,
		PromptToken:     fmt.Sprint(res.UsageMetadata.PromptTokenCount),
		CompletionToken: fmt.Sprint(res.UsageMetadata.CandidatesTokenCount),
		TotalToken:      fmt.Sprint(res.UsageMetadata.TotalTokenCount),
	}

	return resulto.Success(answer)
}
