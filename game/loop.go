package game

import (
	"github.com/bpoetzschke/bin.go/logger"
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
		logger.Debug("Received message: %v", message)
	}
}
