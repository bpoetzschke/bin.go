package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/twinj/uuid"

	"github.com/bpoetzschke/bin.go/logger"
	"github.com/bpoetzschke/bin.go/models"
	slack_middleware "github.com/bpoetzschke/bin.go/slack-middleware"
	"github.com/bpoetzschke/bin.go/storage"
	"github.com/bpoetzschke/bin.go/word_manager"
)

type GameLoop interface {
	Run()
}

func NewGameLoop(
	slackMW slack_middleware.Middleware,
	storage storage.Storage,
	wordManager word_manager.WordManager,
) (GameLoop, error) {
	gl := gameLoop{
		slackMW:     slackMW,
		storage:     storage,
		wordManager: wordManager,
	}

	if err := gl.init(); err != nil {
		return nil, err
	}

	return &gl, nil
}

type gameLoop struct {
	slackMW     slack_middleware.Middleware
	wordManager word_manager.WordManager
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
		switch message.Type {
		case slack_middleware.MessageTypeChannelMessage:
			gl.handleChannelMessage(message)
		}
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
	if err != nil {
		logger.Error("Error while loading word list. %s", err)
	}

	if len(words) == 0 {
		words = gl.wordManager.LoadInitialWords()
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

func (gl *gameLoop) handleChannelMessage(msg *slack_middleware.IncomingMessage) {
	foundWords := models.WordList{}

	for _, word := range gl.currentGame.RemainingWords {
		if strings.Contains(strings.ToLower(msg.Message), strings.ToLower(word.Value)) {
			if word.AddedBy != msg.UserID {
				logger.Debug("Found word: %q", word.Value)
				foundWords = append(foundWords, word)
			} else {
				logger.Debug("User %q added word %q. Not going to ack this.", msg.UserID, word.Value)
			}
		}
	}

	if len(foundWords) == 0 {
		logger.Debug("Didn't found any word in message:\n%s", msg.Message)
		if err := gl.react("speak_no_evil", msg); err != nil {
			logger.Error("Error while reacting to message %#v. Error: %s", msg, err)
		}
		return
	}

	if err := gl.react("boom", msg); err != nil {
		logger.Error("Error while reacting to message %#v. Error: %s", msg, err)
	}

	answer := slack_middleware.OutgoingMessage{
		BaseMessage: slack_middleware.BaseMessage{
			Message: fmt.Sprintf("Bingo! <@%s> said: %s.\n\nThere are %d more words to discover", msg.UserID, foundWords.Join(", "), len(gl.currentGame.RemainingWords)),
			Channel: msg.Channel,
		},
	}

	for _, found := range foundWords {
		answer.Attachments = append(answer.Attachments, found.GifUrl)
	}

	gl.slackMW.PostMessage(answer)
}

func (gl *gameLoop) react(emoji string, msg *slack_middleware.IncomingMessage) error {
	return msg.React(emoji)
}
