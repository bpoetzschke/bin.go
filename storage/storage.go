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
	LoadCurrentGame()
	SaveGame()
}

func NewStorage() Storage {
	return NewInMemoryStorage()
}