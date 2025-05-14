package disocrd

import (
	"LOL-Assistant/gemini"
	"context"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// 전역 변수로 Gemini 클라이언트 선언
var geminiClient gemini.ChatSession

// Initialize 함수 추가 - 봇 시작 시 한 번만 호출
func Initialize() {
	geminiClient = gemini.NewGeminiClient()
}

func Message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "디코봇아") {
		delMsg, err := s.ChannelMessageSend(m.ChannelID, "답변을 생성하고 있어요")
		if err != nil {
			log.Println("답변을 생성하고 있어요 실패", err)
			return
		}

		resp, err := geminiClient.ChatWithDiscord(context.Background(), m.Content)
		if err != nil {
			_, sendErr := s.ChannelMessageSend(m.ChannelID, "답변을 생성하지 못했어요")
			if sendErr != nil {
				log.Fatalln("답변생성 실패", err)
				return
			}
			log.Println("gemini api 에러", err)
			return
		}

		msg, err := s.ChannelMessageSend(m.ChannelID, resp)
		if err != nil {
			_, sendErr := s.ChannelMessageSend(m.ChannelID, "답변을 생성하지 못했어요")
			if sendErr != nil {
				log.Fatalln("답변생성 실패", err)
				return
			}
			log.Println("gemini api 에러", err)
			return
		} else {
			log.Println(msg)
		}

		err = s.ChannelMessageDelete(m.ChannelID, delMsg.ID)
		if err != nil {
			log.Println("메세지 삭제 실패", err)
			return
		}
	} else if strings.HasPrefix(m.Content, "내가 마지막으로 플레이한 게임을 분석해줘") {

	}
}
