package mole

import (
	"log"
	"math/rand"

	"github.com/sdbx/hang-man-bot/utils"
)

type preset struct {
	str  string
	hint string
	kor  bool
}

var presets []preset

func addPreset(str string, hint string) {
	if utils.IsKorean(str) == utils.IsEnglish(str) {
		panic("invalid str")
	}

	kor := utils.IsKorean(str)
	log.Printf("preset added str: %s hint: %s kor: %b\n", str, hint, kor)
	presets = append(presets, preset{
		str:  str,
		hint: hint,
		kor:  kor,
	})
}

func addPresets(strs ...string) {
	for i := 0; i < len(strs); i += 2 {
		addPreset(strs[i], strs[i+1])
	}
}

func getPreset() preset {
	return presets[rand.Intn(len(presets))]
}

func init() {
	addPresets(
		"으하하하", "",
		"우하하하", "으으으믐")
}
