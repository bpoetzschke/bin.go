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
	"github.com/bpoetzschke/bin.go/helper"
)

type GameLoop interface {
	Run()
}

func NewGameLoop(
	slackMW slack_middleware.Middleware,
	storage storage.Storage,
	wordManager helper.WordManager,
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
	wordManager helper.WordManager
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
		case slack_middleware.MessageTypeDirectMessage:
			gl.handleDirectMessage(message)
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

	for _, word := range words {
		if _, _, err := gl.storage.AddWord(word); err != nil {
			logger.Error("Failed to add word %#v. Error: %s", word, err)
		}
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
		gl.react(msg, "speak_no_evil")
		return
	}

	gl.react(msg, "boom")

	gl.currentGame.RemainingWords = gl.currentGame.RemainingWords.Diff(foundWords)
	gl.currentGame.FoundWords = append(gl.currentGame.FoundWords, foundWords...)

	if err := gl.save(); err != nil {
		return
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

	if err := gl.slackMW.PostMessage(answer); err != nil {
		logger.Error("Failed to post message. Errors: %s", err)
		return
	}
}

func (gl *gameLoop) handleDirectMessage(msg *slack_middleware.IncomingMessage) {
	if strings.ToLower(msg.Message) == "cheat" {
		gl.react(msg, "see_no_evil")

		answer := slack_middleware.OutgoingMessage{
			BaseMessage: slack_middleware.BaseMessage{
				Message: gl.currentGame.RemainingWords.Join(", "),
				Channel: msg.Channel,
			},
		}

		if err := gl.slackMW.PostMessage(answer); err != nil {
			logger.Error("Failed to post message. Errors: %s", err)
			return
		}
	} else if strings.HasPrefix(strings.ToLower(msg.Message), "add") {
		rawWord := strings.Replace(msg.Message, "add", "", 1)
		rawWord = strings.TrimSpace(rawWord)

		word, err := gl.wordManager.GetWord(rawWord)
		if err != nil {
			return
		}

		word.AddedBy = msg.UserID

		added, existingWord, err := gl.storage.AddWord(word)
		if err != nil {
			logger.Error("Failed to add word %q. %s", rawWord, err)
			return
		}

		if !added {

			logger.Info("Word %q already exists.", rawWord)
			answerMessage := "Word %q was already added"
			parameters := []interface{}{
				rawWord,
			}
			if existingWord.AddedBy != "" {
				answerMessage += " by <@%s>"
				parameters = append(parameters, existingWord.AddedBy)
			}
			answerMessage += "."
			answer := slack_middleware.OutgoingMessage{
				BaseMessage: slack_middleware.BaseMessage{
					Message: fmt.Sprintf(answerMessage, parameters...),
					Channel: msg.Channel,
				},
			}

			if err := gl.slackMW.PostMessage(answer); err != nil {
				logger.Error("Failed to post message %#v. Error: %s", answer, err)
			}
			return
		}

		gl.currentGame.AddNewWord(word)
		gl.react(msg, "white_check_mark")
	} else {
		gl.react(msg, "question")
	}
}

func (gl *gameLoop) react(msg *slack_middleware.IncomingMessage, reactions ...string) {
	if err := msg.React(reactions...); err != nil {
		logger.Error("Error while reacting to message %#v. Error: %s", msg, err)
	}
}

func (gl *gameLoop) save() error {
	if err := gl.storage.SaveGame(*gl.currentGame); err != nil {
		logger.Error("Error while saving current game. %s", err)
		return err
	}

	return nil
}
