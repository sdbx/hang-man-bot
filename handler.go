package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sdbx/hang-man-bot/config"
)

func newMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Content == "" {
		return
	}

	//gameChannel
	if m.ChannelID == config.Conf.ChannelID {
		gameChannelCmd(s, m)
		mol.MessageHandler(s, m)
		return
	}

	//userChannel
	userChannel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return
	}

	if userChannel.ID == m.ChannelID {
		userChannelCmd(s, m)
	}
}
