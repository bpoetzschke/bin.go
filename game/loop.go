package game

import (
	"time"

	"github.com/twinj/uuid"

	"github.com/bpoetzschke/bin.go/logger"
	"github.com/bpoetzschke/bin.go/models"
	slack_middleware "github.com/bpoetzschke/bin.go/slack-middleware"
	"github.com/bpoetzschke/bin.go/storage"
)

type GameLoop interface {
	Run()
}

func NewGameLoop(
	slackMW slack_middleware.Middleware,
	storage storage.Storage,
) (GameLoop, error) {
	gl := gameLoop{
		slackMW: slackMW,
		storage: storage,
	}

	if err := gl.init(); err != nil {
		return nil, err
	}

	return &gl, nil
}

type gameLoop struct {
	slackMW     slack_middleware.Middleware
	storage     storage.Storage
	currentGame *models.Game
}

func (gl *gameLoop) init() error {
	game, found, err := gl.storage.LoadCurrentGame()
	if err != nil {
		logger.Error("Error while loading current game. %s", err)
		return err
	}

	if !found {
		return gl.startNewGame()
	}

	gl.currentGame = &game

	return nil
}

func (gl *gameLoop) Run() {
	for message := range gl.slackMW.Connect() {
		logger.Debug("Received message: %v", message)
	}
}

func (gl *gameLoop) startNewGame() error {
	if gl.currentGame != nil {
		logger.Debug("Finishing old game before starting new one.")
		now := time.Now().UTC()
		gl.currentGame.FinishedAt = &now
		if err := gl.storage.SaveGame(*gl.currentGame); err != nil {
			logger.Error("Error while saving current game. %s", err)
			return err
		}
	}

	words, err := gl.storage.LoadWordList()
	// TODO if word list is empty load initial word list
	if err != nil {
		logger.Error("Error while loading word list. %s", err)
	}

	logger.Debug("Starting new game.")

	gl.currentGame = &models.Game{
		ID:             uuid.NewV4().String(),
		RemainingWords: words,
		StartedAt:      time.Now().UTC(),
	}
	if err := gl.storage.SaveGame(*gl.currentGame); err != nil {
		logger.Error("Error while saving current game. %s", err)
		return err
	}

	return nil
}
