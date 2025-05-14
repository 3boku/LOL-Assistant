package gemini

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"net/http"
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

	files, err := os.ReadDir("pdfs")
	if err != nil {
		log.Fatal(err)
	}

	var names []string
	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}

	for _, name := range names {
		wikiData, err := os.ReadFile(fmt.Sprintf("pdfs/%s", name))
		if err != nil {
			log.Fatal(err)
		}
		wikiMimeType := http.DetectContentType(wikiData)

		ChatHistory = append(ChatHistory, &genai.Content{
			Parts: []genai.Part{
				genai.Blob{
					MIMEType: wikiMimeType,
					Data:     wikiData,
				},
			},
			Role: "model",
		})
	}

	cs := model.StartChat()

	return ChatSession{
		ChatSesstion: cs,
	}
}

func (cs ChatSession) ChatWithDiscord(ctx context.Context, text string) (string, error) {
	cs.ChatSesstion.History = ChatHistory

	resp, err := cs.ChatSesstion.SendMessage(ctx, genai.Text(text))
	if err != nil {
		return "", err
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

	return string(content.(genai.Text)), nil
}
