package game

type Event interface {
	sealed()
}

type EndEvent struct {
	Win bool
}

func (EndEvent) sealed() {}

type TurnEvent struct {
	Right  bool
	Char   rune
	Player string

	Current []rune
	Log     []rune
	Hp      int
}

func (TurnEvent) sealed() {}

type HintEvent struct {
	Hint string
}

func (HintEvent) sealed() {}

type RevealEvent struct {
	Char rune

	Current []rune
	Log     []rune
	Hp      int
}

func (RevealEvent) sealed() {}
