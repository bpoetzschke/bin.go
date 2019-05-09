package storage

import "github.com/bpoetzschke/bin.go/models"

type Storage interface {
	WordStorage
	GameStorage
}

type WordStorage interface {
	LoadWordList() ([]models.Word, error)
	AddWord(models.Word) (bool, error)
}

type GameStorage interface {
	LoadCurrentGame() (models.Game, error)
	SaveGame(models.Game) error
}

func NewStorage() Storage {
	return NewInMemoryStorage()
}
