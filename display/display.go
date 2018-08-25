package display

import (
	"path/filepath"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/sdbx/hang-man-bot/config"
	"github.com/sdbx/hang-man-bot/utils"
)

func DisplayDefault(manID int, creator string, hint string, hp int, current []rune, log []rune) *discordgo.MessageEmbed {
	maxhp := config.Conf.MaxHp

	str := "`" + utils.MakeSpace(current) +
		"`\n:heart: > " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) +
		"\n:skull: > " + utils.MakeSpace(log)
	if hint != "" {
		str += "\n힌트:" + hint
	}
	str += "\n출제자: " + creator

	return makeEmbed(str, maxhp, manID, hp)
}
func DisplayDied(answer []rune, hint string, manID int, creator string) *discordgo.MessageEmbed {
	maxhp := config.Conf.MaxHp

	str := "목숨을 모두 소모하여 행맨을 죽였습니다!" +
		"\n정답: " + string(answer)
	if hint != "" {
		str += "\n힌트:" + hint
	}
	str += "\n출제자: " + creator

	return makeEmbed(str, maxhp, manID, 0)
}

func DisplayWin(answer []rune, creator string, hint string, manID int) *discordgo.MessageEmbed {
	maxhp := config.Conf.MaxHp

	str := "정답을 모두 맞추어 행맨을 살렸습니다!" +
		"\n정답: " + string(answer)
	if hint != "" {
		str += "\n힌트:" + hint
	}
	str += "\n출제자: " + creator
	return makeEmbed(str, maxhp, manID, maxhp)
}

func makeEmbed(desc string, maxhp int, manID int, hp int) *discordgo.MessageEmbed {
	pic := GetManImage(maxhp, manID, hp)
	ext := filepath.Ext(pic)
	return &discordgo.MessageEmbed{
		Description: desc,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: utils.CreateManURL(config.Conf.URL, maxhp, manID, hp, ext),
		},
	}
}
