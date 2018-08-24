package display

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrContainsDirectory = errors.New("sdbx/hang-man-bot/display man directory contains a directory")
	ErrDuplicate         = errors.New("sdbx/hang-man-bot/display duplicate image")
	ErrMissing           = errors.New("sdbx/hang-man-bot/display missing image")
)

type Man []string

var mans map[int][]Man

func PickID(maxhp int) int {
	return rand.Intn(len(mans[maxhp]))
}

func GetPicture(maxhp int, id int, hp int) string {
	if maxhp < 0 {
		return ""
	}

	if _, ok := mans[maxhp]; !ok {
		return ""
	}

	if id >= len(mans[maxhp]) || id < 0 {
		return ""
	}

	if hp > maxhp || hp < 0 {
		return ""
	}

	return mans[maxhp][id][hp]
}

func init() {
	mans = make(map[int][]Man)
}

func addMan(maxhp int, man Man) {
	if _, ok := mans[maxhp]; !ok {
		mans[maxhp] = []Man{}
	}
	mans[maxhp] = append(mans[maxhp], man)
}

func InitMans(root string) error {

	// [id][maxhp][]path
	assets := make(map[int]map[int][]string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(info.Name())
		name := strings.TrimSuffix(info.Name(), ext)
		var (
			id    int
			maxhp int
			hp    int
		)
		fmt.Sscanf(name, "%d_%d_%d", &id, &maxhp, &hp)
		if _, ok := assets[id]; !ok {
			assets[id] = make(map[int][]string)
		}

		if _, ok := assets[id][maxhp]; !ok {
			assets[id][maxhp] = make([]string, maxhp+1)
		}

		if e := assets[id][maxhp][hp]; e != "" {
			return ErrDuplicate
		}

		assets[id][maxhp][hp] = path
		return nil
	})
	if err != nil {
		return err
	}

	for _, mans := range assets {
		for maxhp, man := range mans {
			m := make(Man, 0, maxhp)
			for _, picture := range man {
				if picture == "" {
					return ErrMissing
				}
				m = append(m, picture)
			}
			addMan(maxhp, m)
		}
	}
	return nil
}
