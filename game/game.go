package game

import (
	"sort"
	"sync"
	"time"
)

type runes []rune

func (a runes) Len() int           { return len(a) }
func (a runes) Less(i, j int) bool { return a[i] < a[j] }
func (a runes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type Game struct {
	mu sync.RWMutex

	creator string

	solution []rune
	mask     []bool

	log []rune

	hp    int
	maxHp int

	cools map[string]time.Time

	evs []chan Event
}

func (g *Game) Play(player string, char rune) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if r, yes := g.isCool(player); yes {
		return &CoolError{
			Player: player,
			Remain: r,
		}
	}

	if player == g.creator {
		return ErrYouCreator
	}

	if test(char) {

	} else {

	}

}

func (g *Game) test(char rune) bool {
	ok := false
	for i, c := range g.solution {
		if c == char && g.mask[i] {
			g.mask[i] = false
			ok = true
		}
	}
	if !ok {
		g.addLog(char)
		return false
	}

}

func (g *Game) addLog(char rune) {
	for _, c := range g.log {
		if char == c {
			return
		}
	}

	g.log = append(g.log, char)
	sort.Sort(g.log)
}

func (g *Game) isCool(player string) (time.Duration, bool) {
	if c, ok := g.cools[player]; ok {
		if c.After(time.Now()) {
			return 0, false
		}
		return c.Sub(time.Now()), true
	}
	return 0, false
}

func (g *Game) setCool(player string, d time.Duration) {
	g.cools[player] = time.Now().Add(d)
}

func (g *Game) resetAllCool() {
	g.cools = make(map[string]time.Time)
}

func (g *Game) Current() []rune {
	g.mu.RLock()
	defer g.mu.RUnlock()

	out := make([]rune, len(g.solution))
	for i, c := range g.solution {
		if g.mask[i] {
			out[i] = c
		} else {
			out[i] = '_'
		}
	}
	return out
}

func (g *Game) MaxHP() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.maxHp
}

func (g *Game) Hp() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.hp
}

func (g *Game) Log() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return string(g.log)
}

func (g *Game) Listen() <-chan Event {
	g.mu.Lock()
	defer g.mu.Unlock()

	c := make(chan Event)
	g.evs = append(g.evs, c)
	return c
}
