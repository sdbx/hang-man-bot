package game

type Event interface {
	sealed()
}

type EndEvent struct {
}

func (EndEvent) sealed() {}
