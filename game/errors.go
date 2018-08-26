package game

import (
	"errors"
	"time"
)

var (
	ErrYouCreator   = errors.New("sdbx/hang-man-bot/game you are the creator")
	ErrAlreadyInput = errors.New("sdbx/hang-man-bot/game already inputed value")
	ErrEnded        = errors.New("sdbx/hang-man-bot/game already ended game")
)

type CoolError struct {
	Player string
	Remain time.Duration
}

func (c CoolError) Error() string {
	return "sdbx/hang-man-bot/game cooling down: " + c.Remain.String()
}
