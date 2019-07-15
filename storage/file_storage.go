package storage

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/bpoetzschke/bin.go/logger"

	"github.com/bpoetzschke/bin.go/models"
)

const fileName = "bin.go_save.json"

func NewFileStorage() Storage {
	return &fileStorage{}
}

type fileModel struct {
	WordList []models.Word          `json:"word_list"`
	GameMap  map[string]models.Game `json:"game_map"`
}

type fileStorage struct {
}

func (fs fileStorage) LoadWordList() (words []models.Word, err error) {
	content := fs.loadFile()

	return content.WordList, nil
}

func (fs fileStorage) AddWord(newWord models.Word) (bool, models.Word, error) {
	content := fs.loadFile()

	for _, word := range content.WordList {
		if strings.ToLower(word.Value) == strings.ToLower(newWord.Value) {
			return false, word, nil
		}
	}

	content.WordList = append(content.WordList, newWord)
	if err := fs.saveFile(content); err != nil {
		return false, models.Word{}, err
	}

	return true, models.Word{}, nil
}

func (fs fileStorage) LoadCurrentGame() (models.Game, bool, error) {
	content := fs.loadFile()

	for _, game := range content.GameMap {
		if game.FinishedAt == nil {
			return game, true, nil
		}
	}

	return models.Game{}, false, nil
}

func (fs fileStorage) SaveGame(game models.Game) error {
	content := fs.loadFile()

	content.GameMap[game.ID] = game

	return fs.saveFile(content)
}

func (fs fileStorage) loadFile() fileModel {
	content := fileModel{
		GameMap:  map[string]models.Game{},
		WordList: []models.Word{},
	}

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return content
	}

	if err := json.Unmarshal(bytes, &content); err != nil {
		logger.Warning("Failed to unmarshal content of save file. Returning default content. Error: %s", err)
		return content
	}

	return content
}

func (fs fileStorage) saveFile(content fileModel) error {
	bytes, err := json.Marshal(content)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, bytes, 0644)
}
