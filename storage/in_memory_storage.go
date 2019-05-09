package storage

import (
	"fmt"

	"github.com/bpoetzschke/bin.go/models"
)

func NewInMemoryStorage() Storage {
	s := inMemoryStorage{}
	s.init()

	return &s
}

type inMemoryStorage struct {
	wordMap map[string]models.Word
	gameMap map[string]*models.Game
}

func (ims *inMemoryStorage) init() {
	ims.wordMap = make(map[string]models.Word)
	ims.gameMap = make(map[string]*models.Game)
}

func (ims *inMemoryStorage) LoadWordList() ([]models.Word, error) {
	var wordList = make([]models.Word, 0)

	for _, word := range ims.wordMap {
		wordList = append(wordList, word)
	}

	return wordList, nil
}

func (ims *inMemoryStorage) AddWord(word models.Word) (bool, error) {
	_, found := ims.wordMap[word.Value]
	if found {
		return false, fmt.Errorf("word with value %s already exists", word.Value)
	}

	ims.wordMap[word.Value] = word

	return true, nil
}

func (ims *inMemoryStorage) LoadCurrentGame() (models.Game, error) {

	for _, gamePtr := range ims.gameMap {
		if gamePtr.FinishedAt == nil {
			return *gamePtr, nil
		}
	}

	return models.Game{}, nil
}

func (ims *inMemoryStorage) SaveGame(game models.Game) error {
	ims.gameMap[game.ID] = &game

	return nil
}
