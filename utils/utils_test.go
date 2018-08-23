package utils_test

import (
	"testing"

	"github.com/sdbx/hang-man-bot/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsKorean(t *testing.T) {
	a := assert.New(t)
	a.Equal(true, utils.IsKorean("안녕하세요"))
	a.Equal(true, utils.IsKorean("가나다바마바사"))
	a.Equal(false, utils.IsKorean("asdads"))
	a.Equal(false, utils.IsKorean("ㄱㄴㄷㄱㄹㅎㄴ"))
	a.Equal(false, utils.IsKorean("가 나 다 라"))
	a.Equal(false, utils.IsKorean("$@$@#$@$"))
	a.Equal(true, utils.IsKorean("갉낡"))
}
