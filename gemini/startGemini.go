package gemini

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// ChatSession 구조체는 Gemini AI와의 채팅 세션을 관리합니다.
type ChatSession struct {
	ChatSesstion *genai.ChatSession
}

// ChatHistory는 채팅 기록을 저장하는 전역 변수입니다.
var ChatHistory []*genai.Content

// NewGeminiClient 함수는 새로운 Gemini AI 클라이언트를 생성합니다.
// API 키를 로드하고 모델을 초기화하며, 안전 설정을 구성하고
// PDF 파일을 로드하여 ChatHistory에 추가합니다.
// 마지막으로 채팅 세션을 시작하고 반환합니다.
func NewGeminiClient() ChatSession {
	// 환경 변수에서 API 키 로드
	apiKey := os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()

	// Gemini API 클라이언트 생성
	var err error
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	// 환경 변수에서 시스템 지시사항 로드
	instructions := os.Getenv("GEMINI_INSTRUCTIONS")
	// 모델 초기화 (gemini-2.5-flash-preview-04-17 버전 사용)
	model := client.GenerativeModel("gemini-2.5-flash-preview-04-17")
	// 시스템 지시사항 설정
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(instructions)},
	}

	// 안전 설정 구성 (모든 카테고리에서 필터링 없음)
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

	// pdfs 디렉토리에서 파일 목록 읽기
	files, err := os.ReadDir("pdfs")
	if err != nil {
		log.Fatal(err)
	}

	// 파일 이름 배열 생성
	var names []string
	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}

	// 각 PDF 파일 로드 및 ChatHistory에 추가
	for _, name := range names {
		wikiData, err := os.ReadFile(fmt.Sprintf("pdfs/%s", name))
		if err != nil {
			log.Fatal(err)
		}
		// 파일 MIME 타입 감지
		wikiMimeType := http.DetectContentType(wikiData)

		// 파일 데이터를 모델 역할로 ChatHistory에 추가
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

	// 채팅 세션 시작
	cs := model.StartChat()

	// ChatSession 객체 반환
	return ChatSession{
		ChatSesstion: cs,
	}
}

// ChatWithDiscord 메서드는 디스코드에서 받은 텍스트를 Gemini AI에 전송하고
// 응답을 받아 반환합니다. 대화 기록은 ChatHistory에 저장됩니다.
// 오류가 발생하면 빈 문자열과 오류를 반환합니다.
func (cs ChatSession) ChatWithDiscord(ctx context.Context, text string) (string, error) {
	// 현재 채팅 세션에 대화 기록 설정
	cs.ChatSesstion.History = ChatHistory

	// Gemini AI에 메시지 전송
	resp, err := cs.ChatSesstion.SendMessage(ctx, genai.Text(text))
	if err != nil {
		return "", err
	}

	// 응답에서 콘텐츠 추출
	var content genai.Part
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				content = part
			}
		}
	}

	// 사용자 메시지를 대화 기록에 추가
	ChatHistory = append(ChatHistory, &genai.Content{
		Parts: []genai.Part{
			genai.Text(text),
		},
		Role: "user",
	})

	// AI 응답을 대화 기록에 추가
	ChatHistory = append(ChatHistory, &genai.Content{
		Parts: []genai.Part{
			content,
		},
		Role: "model",
	})

	// 텍스트 형식으로 응답 반환
	return string(content.(genai.Text)), nil
}
