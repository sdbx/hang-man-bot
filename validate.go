package main

import "github.com/sdbx/hang-man-bot/utils"

func isRightLetter(r rune) bool {
	return !(utils.IsEnglishRune(r) == utils.IsKoreanRune(r))
}

func isTwoLetter(str string) bool {
	return len([]rune(str)) == 2
}
