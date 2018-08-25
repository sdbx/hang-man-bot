package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/sdbx/hang-man-bot/config"
)

func newMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//gameChannel
	if m.ChannelID == config.Conf.ChannelID {
		go mol.MessageHandler(s, m)
		gameChannelCmd(s, m)
		return
	}

	//userChannel
	userChannel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Println(err)
		return
	}

	if userChannel.ID == m.ChannelID {
		userChannelCmd(s, m)
	}
}
