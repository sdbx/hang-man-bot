package main

import (
	"fmt"
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
		fmt.Println(err)
		return
	}

	discord, err := discordgo.New("Bot " + config.Conf.Token)

	if err != nil {
		fmt.Println(err)
		return
	}

	discord.AddHandler(newMessageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = display.InitMans(config.Conf.AssetsDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	imgserv.Start()

	mol = mole.New(discord)
	go mol.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
