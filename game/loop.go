package game

import (
	"fmt"

	"github.com/nlopes/slack"
)

type GameLoop interface {
	Run()
}

func NewGameLoop(messageChan <-chan *slack.MessageEvent) GameLoop {
	return &gameLoop{
		messageChan: messageChan,
	}
}

type gameLoop struct {
	messageChan <-chan *slack.MessageEvent
}

func (gl *gameLoop) Run() {
	for message := range gl.messageChan {
		fmt.Printf("Received message: %+v", message)
	}
}
