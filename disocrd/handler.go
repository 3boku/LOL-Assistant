package disocrd

import (
	"LOL-Assistant/gemini"
	"context"
	"github.com/bwmarrin/discordgo"
)

func Message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	cs := gemini.NewGeminiClient()

	resp := cs.ChatWithDiscord(context.Background(), m.Content)

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, resp)
	}
}
