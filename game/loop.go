package game

import (
	"github.com/bpoetzschke/bin.go/logger"
	slack_middleware "github.com/bpoetzschke/bin.go/slack-middleware"
)

type GameLoop interface {
	Run()
}

func NewGameLoop(messageChan <-chan *slack_middleware.Message) GameLoop {
	return &gameLoop{
		messageChan: messageChan,
	}
}

type gameLoop struct {
	messageChan <-chan *slack_middleware.Message
}

func (gl *gameLoop) Run() {
	for message := range gl.messageChan {
		logger.Debug("Received message: %v", message)
	}
}
