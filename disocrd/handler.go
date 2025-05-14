package disocrd

import (
	"LOL-Assistant/gemini"
	"context"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func Message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	cs := gemini.NewGeminiClient()

	resp := cs.ChatWithDiscord(context.Background(), m.Content)

	if strings.Contains(m.Content, "디코봇아") {
		msg, err := s.ChannelMessageSend(m.ChannelID, resp)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(msg)
		}
	}
}
