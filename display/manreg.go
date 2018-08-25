package display

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/sdbx/hang-man-bot/config"
)

var (
	ErrContainsDirectory = errors.New("sdbx/hang-man-bot/display man directory contains a directory")
	ErrDuplicate         = errors.New("sdbx/hang-man-bot/display duplicate image")
	ErrMissing           = errors.New("sdbx/hang-man-bot/display missing image")
)

type Man []string

var mans map[int][]Man

func InitMans() {
	mans = make(map[int][]Man)
	err := initMans()
	if err != nil {
		panic(err)
	}
}

func PickManID(maxhp int) int {
	return rand.Intn(len(mans[maxhp]))
}

func GetManImage(maxhp int, id int, hp int) string {
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

func addMan(maxhp int, man Man) {
	if _, ok := mans[maxhp]; !ok {
		mans[maxhp] = []Man{}
	}
	mans[maxhp] = append(mans[maxhp], man)
}

func initMans() error {
	// [id][maxhp][]path
	assets := make(map[int]map[int][]string)
	log.Println("probing " + config.Conf.MansDir)
	err := filepath.Walk(config.Conf.MansDir, func(path string, info os.FileInfo, err error) error {
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

	for id, mans := range assets {
		for maxhp, man := range mans {
			m := make(Man, 0, maxhp)
			for hp, picture := range man {
				if picture == "" {
					return ErrMissing
				}
				log.Printf("adding man id:%d maxhp:%d hp:%d\n", id, maxhp, hp)
				m = append(m, picture)
			}
			addMan(maxhp, m)
		}
	}
	return nil
}
