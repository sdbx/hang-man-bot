package utils

import (
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

func IsEnglish(str string) bool {
	for _, r := range str {
		if !unicode.IsUpper(r) && !unicode.IsLower(r) {
			return false
		}
	}
	return true
}

func IsKorean(str string) bool {
	for _, r := range str {
		if !unicode.In(r, &koreanRange) {
			return false
		}
	}
	return true
}
