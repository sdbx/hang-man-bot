package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var Conf Config

type Config struct {
	URL       string `yaml:"url"`
	Token     string `yaml:"token"`
	ChannelID string `yaml:"channel_id"`
	AssetsDir string `yaml:"assets_dir"`
	Suffix    string `yaml:"suffix"`
	Cool      int    `yaml:"cool"`
	MaxHp     int    `yaml:"max_hp"`
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
		URL:       "127.0.0.1:8053",
		Token:     "Write Token Here",
		ChannelID: "Write ChannelID Here",
		AssetsDir: "assets",
		Cool:      3,
		Suffix:    ".",
		MaxHp:     7,
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
