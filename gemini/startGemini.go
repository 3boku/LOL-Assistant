package gemini

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"os"
)

type ChatSession struct {
	ChatSesstion *genai.ChatSession
}

var ChatHistory []*genai.Content

func NewGeminiClient() ChatSession {
	apiKey := os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()

	var err error
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	instructions := os.Getenv("GEMINI_INSTRUCTIONS")
	model := client.GenerativeModel("gemini-2.0-flash")
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(instructions)},
	}
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}

	cs := model.StartChat()

	return ChatSession{
		ChatSesstion: cs,
	}
}

func (cs ChatSession) ChatWithDiscord(ctx context.Context, text string) string {
	cs.ChatSesstion.History = ChatHistory

	resp, err := cs.ChatSesstion.SendMessage(ctx, genai.Text(text))
	if err != nil {
		log.Fatal(err)
	}

	var content genai.Part
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				content = part
			}
		}
	}

	ChatHistory = append(ChatHistory, &genai.Content{
		Parts: []genai.Part{
			genai.Text(text),
		},
		Role: "user",
	})
	ChatHistory = append(ChatHistory, &genai.Content{
		Parts: []genai.Part{
			content,
		},
		Role: "model",
	})

	return string(content.(genai.Text))
}
