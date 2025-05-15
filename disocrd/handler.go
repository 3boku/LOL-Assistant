package disocrd

import (
	"LOL-Assistant/gemini"
	"LOL-Assistant/league"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// 전역 변수로 Gemini 클라이언트 선언
var geminiClient gemini.ChatSession

// Initialize 함수는 봇 시작 시 한 번만 호출되어
// Gemini AI 클라이언트를 초기화합니다.
func Initialize() {
	geminiClient = gemini.NewGeminiClient()
}

// Message 함수는 디스코드 메시지 이벤트 핸들러입니다.
// 사용자가 보낸 메시지를 받아 처리하고, 조건에 따라 응답을 전송합니다.
// 봇 자신이 보낸 메시지는 무시합니다.
func Message(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 봇 자신의 메시지는 무시
	if m.Author.ID == s.State.User.ID {
		return
	}

	// "디코봇아"가 포함된 메시지에 대한 처리
	if strings.Contains(m.Content, "디코봇아") {
		// "답변을 생성하고 있어요" 메시지 전송
		delMsg, err := s.ChannelMessageSend(m.ChannelID, "답변을 생성하고 있어요")
		if err != nil {
			log.Println("답변을 생성하고 있어요 실패", err)
			return
		}
		// Gemini AI로 사용자 메시지에 대한 응답 생성
		resp, err := geminiClient.ChatWithDiscord(context.Background(), m.Content)
		if err != nil {
			// 응답 생성 실패 시 에러 메시지 전송
			_, sendErr := s.ChannelMessageSend(m.ChannelID, "답변을 생성하지 못했어요")
			if sendErr != nil {
				log.Fatalln("답변생성 실패", err)
				return
			}
			log.Println("gemini api 에러", err)
			return
		}

		// Gemini AI의 응답을 디스코드 채널에 전송
		msg, err := s.ChannelMessageSend(m.ChannelID, resp)
		if err != nil {
			// 응답 전송 실패 시 에러 메시지 전송
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

		// 이전에, "답변을 생성하고 있어요" 메시지 삭제
		err = s.ChannelMessageDelete(m.ChannelID, delMsg.ID)
		if err != nil {
			log.Println("메세지 삭제 실패", err)
			return
		}
	} else if strings.HasPrefix(m.Content, "내가 마지막으로 플레이한 게임을 분석해줘") {
		// ex: "내가 마지막으로 플레이한 게임을 분석해줘|닉네임#태그
		delMsg, err := s.ChannelMessageSend(m.ChannelID, "답변을 생성하고 있어요")
		if err != nil {
			log.Println("답변을 생성하고 있어요 실패", err)
			return
		}

		split := strings.Split(m.Content, "|")
		nickName := strings.Split(split[1], "#")[0]
		tag := strings.Split(split[1], "#")[1]

		gameInfo, puuid, err := league.GetMatch(nickName, tag)
		if err != nil {
			if err != nil {
				// 응답 전송 실패 시 에러 메시지 전송
				_, sendErr := s.ChannelMessageSend(m.ChannelID, "답변을 생성하지 못했어요")
				if sendErr != nil {
					log.Fatalln("답변생성 실패", err)
					return
				}
				log.Println("gemini api 에러", err)
				return
			}
			err = s.ChannelMessageDelete(m.ChannelID, delMsg.ID)
			if err != nil {
				log.Println("메세지 삭제 실패", err)
				return
			}
		}
		matchReq := fmt.Sprintf("%s | 나의 puuid: %s, 게임정보: %s", m.Content, puuid, gameInfo)

		resp, err := geminiClient.ChatWithDiscord(context.Background(), matchReq)
		if err != nil {
			// 응답 생성 실패 시 에러 메시지 전송
			_, sendErr := s.ChannelMessageSend(m.ChannelID, "답변을 생성하지 못했어요")
			if sendErr != nil {
				log.Fatalln("답변생성 실패", err)
				return
			}
			log.Println("gemini api 에러", err)
			return
		}
		log.Println(resp)

		_, err = s.ChannelMessageSend(m.ChannelID, resp)
		if err != nil {
			_, sendErr := s.ChannelMessageSend(m.ChannelID, "답변을 생성하지 못했어요")
			if sendErr != nil {
				log.Fatalln("답변생성 실패", err)
				return
			}
			log.Println("gemini api 에러", err)
			return
		}
		err = s.ChannelMessageDelete(m.ChannelID, delMsg.ID)
		if err != nil {
			log.Println("메세지 삭제 실패", err)
			return
		}
	}
}
