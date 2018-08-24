package display

import (
	"strconv"

	"github.com/sdbx/hang-man-bot/utils"
)

func DisplayDefault(maxHp int, hp int, log []rune, current []rune) string {
	return ("`" + utils.MakeSpace(current) +
		"`\n:heart: > " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxHp) +
		"\n:skull: > " + utils.MakeSpace(log))
}

func DisplayDie(creator string) string {
	return ""
}
