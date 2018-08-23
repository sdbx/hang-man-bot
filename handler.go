package main

import "github.com/bwmarrin/discordgo"

func newMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != config.channelID {
		return
	}
	if m.Content == "" {
		return
	}

	if()
}
