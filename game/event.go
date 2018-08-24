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
	Player string
}

func (TurnEvent) sealed() {}

type HintEvent struct {
	Hint string
}

func (HintEvent) sealed() {}

type RevealEvent struct {
	Char rune
}

func (RevealEvent) sealed() {}
