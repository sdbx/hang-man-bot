package game

type Event interface {
	sealed()
}

type EndEvent struct {
}

type WrongEvent struct {
	Player string
}

type RightEvent struct {
	Player string
}

func (EndEvent) sealed() {}
