package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	token     string
	channelID string
}

var config Config

func main() {

	err := loadConfig()
	if err != nil {
		fmt.Println(err.error())
	}

	discord, err := discordgo.New("Bot " + config.token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//handlers
	discord.AddHandler(newMessageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//sc
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()

}

func loadConfig() error {
	buf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			return createConfig()
		}
		return err
	}
	return yaml.Unmarshal(buf, &config)

}
func createConfig() error {
	config = config{
		token:     "Write Token Here",
		ChannelID: "Write ChannelID Here",
	}

	buf, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("config.yaml", buf, 0644)
	if err != nil {
		return err
	}
}
