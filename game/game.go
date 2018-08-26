package game

import (
	"log"
	"sort"
	"time"

	uuid "github.com/odeke-em/go-uuid"
	"github.com/sdbx/hang-man-bot/utils"
)

type runes []rune

func (a runes) Len() int           { return len(a) }
func (a runes) Less(i, j int) bool { return a[i] < a[j] }
func (a runes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a runes) Contains(r rune) bool {
	for _, c := range a {
		if c == r {
			return true
		}
	}
	return false
}

type GameMeta struct {
	Creator  string
	Hint     string
	Korean   bool
	Solution []rune
}

type Game struct {
	mu utils.LogLock

	creator string
	hint    string
	korean  bool

	solution []rune
	mask     []bool

	log runes

	hp    int
	maxHp int

	cool  time.Duration
	cools map[string]time.Time

	evs []chan Event
}

func New(creator string, hint string, korean bool, cool time.Duration, solution []rune, maxHp int) *Game {
	return &Game{
		creator:  creator,
		hint:     hint,
		korean:   korean,
		solution: solution,
		mask:     createMask(len(solution)),
		log:      runes{},
		hp:       maxHp,
		maxHp:    maxHp,
		cool:     cool,
		cools:    make(map[string]time.Time),
		evs:      []chan Event{},
	}
}

func createMask(n int) []bool {
	out := make([]bool, n)
	for i := range out {
		out[i] = true
	}
	return out
}

func (g *Game) Play(player string, char rune) error {
	if player == g.creator {
		return ErrYouCreator
	}

	g.mu.RLock()
	r, yes := g.isCool(player)
	l := g.log
	current := runes(g.current())
	hp := g.hp
	g.mu.RUnlock()

	if hp <= 0 {
		return ErrEnded
	}

	if yes {
		return &CoolError{
			Player: player,
			Remain: r,
		}
	}

	if l.Contains(char) || current.Contains(char) {
		return ErrAlreadyInput
	}

	g.mu.Lock()

	g.resetAllCool()
	g.setCool(player, g.cool)
	rv := g.reveal(char)
	if !rv {
		g.addLog(char)
		g.hp--
	}

	current = g.current()
	l = g.log

	hp = g.hp
	cr := g.cleared()

	g.mu.Unlock()

	if rv {
		g.sendEvent(&TurnEvent{
			Right:   true,
			Char:    char,
			Player:  player,
			Current: current,
			Log:     l,
			Hp:      hp,
		})

	} else {
		g.sendEvent(&TurnEvent{
			Right:   false,
			Char:    char,
			Player:  player,
			Current: current,
			Log:     l,
			Hp:      hp,
		})
	}

	if hp == int(g.maxHp/2) {
		g.sendEvent(&HintEvent{
			Hint: g.hint,
		})
	}

	if hp <= 0 {
		g.sendEvent(&EndEvent{
			Win: false,
		})
	} else if cr {
		g.sendEvent(&EndEvent{
			Win: true,
		})
	}

	return nil
}

func (g *Game) sendEvent(e Event) {
	for _, ev := range g.evs {
		if ev == nil {
			continue
		}
		log.Println(uuid.New(), "event send")
		ev <- e
		log.Println(uuid.New(), "event sent")
	}
}

func (g *Game) cleared() bool {
	for _, m := range g.mask {
		if m {
			return false
		}
	}
	return true
}

func (g *Game) reveal(char rune) bool {
	ok := false
	for i, c := range g.solution {
		if c == char && g.mask[i] {
			g.mask[i] = false
			ok = true
		}
	}
	return ok
}

func (g *Game) addLog(char rune) {
	if g.log.Contains(char) {
		return
	}

	g.log = append(g.log, char)
	sort.Sort(g.log)
}

func (g *Game) isCool(player string) (time.Duration, bool) {
	if c, ok := g.cools[player]; ok {
		if c.After(time.Now()) {
			return c.Sub(time.Now()), true
		}
		return 0, false
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

	return g.current()
}

func (g *Game) current() []rune {
	out := make([]rune, len(g.solution))
	for i, c := range g.solution {
		if !g.mask[i] {
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

func (g *Game) Log() []rune {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.log
}

func (g *Game) Creator() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.creator
}

func (g *Game) Hint() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.hint
}

func (g *Game) Korean() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.korean
}

func (g *Game) Listen() (func(), <-chan Event) {
	g.mu.Lock()
	defer g.mu.Unlock()

	c := make(chan Event)
	i := len(g.evs)
	g.evs = append(g.evs, c)
	fn := func() {
		g.mu.Lock()
		defer g.mu.Unlock()

		g.evs[i] = nil
	}

	return fn, c
}
