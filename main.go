package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/sdbx/hang-man-bot/config"
	"github.com/sdbx/hang-man-bot/display"
	"github.com/sdbx/hang-man-bot/imgserv"
	"github.com/sdbx/hang-man-bot/mole"
)

var mol *mole.Mole

func main() {
	err := config.Load()
	if err != nil {
		log.Println(err)
		return
	}

	display.InitMans()
	imgserv.Start()

	discord, err := discordgo.New("Bot " + config.Conf.Token)

	if err != nil {
		log.Println(err)
		return
	}

	discord.AddHandler(newMessageCreate)

	err = discord.Open()
	if err != nil {
		log.Println(err)
		return
	}

	mol = mole.New(discord, config.Conf.ChannelID, config.Conf.LogChannelID)
	go mol.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
