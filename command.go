package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sdbx/hang-man-bot/config"
	"github.com/sdbx/hang-man-bot/mole"
	"github.com/sdbx/hang-man-bot/utils"
)

func cmdHelp(s *discordgo.Session, channelID string) {
	s.ChannelMessage(channelID, "도움: 도움말을 보여줍니다.\n등록 한글단어/영어단어 (힌트메세지) : 출제될 단어를 등록합니다.")
}

func cmdRegister(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "인수를 제대로 입력해주세요. 등록 (단어) (힌트메세지) ")
		return
	}

	iskorean := utils.IsKorean(args[1])
	if iskorean == utils.IsEnglish(args[1]) {
		s.ChannelMessageSend(m.ChannelID, "올바른 단어를 입력해주세요.")
		return
	}

	gm := mole.GameMeta{
		Korean:   iskorean,
		Creator:  m.Author.ID,
		Solution: []rune(args[1]),
	}

	if len(args) >= 3 {
		gm.Hint = args[2]
	} else {
		gm.Hint = ""
	}
	err := mol.Add(gm)
	if err != nil {
		if err == mole.ErrAlreadyExist {
			s.ChannelMessageSend(m.ChannelID, "이미 단어를 등록하셨습니다.")
		}
		return
	}

	s.ChannelMessageSend(m.ChannelID, "성공적으로 단어를 등록하였습니다.")
}

func gameChannelCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)

	if isTwoLetter(m.Content) {
		letter := []rune(m.Content)[0]
		if strings.HasSuffix(m.Content, config.Conf.Suffix) && isRightLetter(letter) {
			if utils.IsKoreanRune(letter) != mol.Game().Korean() {
				return
			}
			err := mol.Game().Play(m.Author.ID, letter)
			if err != nil {
			}
			return
		}
		return
	}
}
func userChannelCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if args[0] == "도움" {
		cmdHelp(s, m.ChannelID)
	}

	if args[0] == "등록" {
		cmdRegister(s, m, args)
	}
}
