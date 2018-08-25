package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var Conf Config

type Config struct {
	Token        string `yaml:"token"`
	ChannelID    string `yaml:"channel_id"`
	LogChannelID string `yaml:"log_channel_id"`
	Suffix       string `yaml:"suffix"`

	Cool  int `yaml:"cool"`
	MaxHp int `yaml:"max_hp"`

	URL     string `yaml:"url"`
	MansDir string `yaml:"mans_dir"`
}

func Load() error {
	_, err := os.Stat("config.yaml")
	if os.IsNotExist(err) {
		err = createConfig()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	buf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, &Conf)
}

func createConfig() error {
	Conf = Config{
		URL:          "127.0.0.1:8053",
		Token:        "Write Token Here",
		ChannelID:    "Write ChannelID Here",
		LogChannelID: "Write LogChannelID Here",
		MansDir:      "images",
		Cool:         3,
		Suffix:       ".",
		MaxHp:        10,
	}

	buf, err := yaml.Marshal(Conf)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("config.yaml", buf, 0644)
	if err != nil {
		return err
	}

	return nil
}
