package utils

import (
	"strconv"
	"unicode"
)

var koreanRange = unicode.RangeTable{
	R32: []unicode.Range32{
		{
			Lo:     '\uAC00',
			Hi:     '\uD7AF',
			Stride: 1,
		},
	},
}

func IsEnglishRune(r rune) bool {
	return unicode.IsUpper(r) || unicode.IsLower(r)
}

func IsEnglish(str string) bool {
	for _, r := range str {
		if !IsEnglishRune(r) {
			return false
		}
	}
	return true
}

func IsKoreanRune(r rune) bool {
	return unicode.In(r, &koreanRange)
}

func IsKorean(str string) bool {
	for _, r := range str {
		if !IsKoreanRune(r) {
			return false
		}
	}
	return true
}

func MakeSpace(rs []rune) string {
	if len(rs) == 0 {
		return ""
	}
	str := ""
	for _, r := range rs {
		str += " " + string([]rune{r})
	}
	return str[1:]
}

func CreateURL(base string, maxhp int, id int, hp int) string {
	return base + "/img/" + strconv.Itoa(maxhp) + "/" + strconv.Itoa(id) + "/" + strconv.Itoa(hp)
}
