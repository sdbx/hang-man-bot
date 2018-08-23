package game

import "math/rand"

var mans map[int][][]string

func registerMan(maxhp int, man []string) {
	mans[maxhp] = append(mans[maxhp], man)
}

func GetMan(maxhp int) []string {
	return mans[maxhp][rand.Intn(len(mans[maxhp]))]
}

func init() {
	registerMan(7, []string{
		"주거따",
		"1",
		"2",
	})
}
