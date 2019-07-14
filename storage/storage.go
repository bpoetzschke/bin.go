package storage

import (
	"github.com/bpoetzschke/bin.go/logger"
	"github.com/bpoetzschke/bin.go/models"
)

type StorageMethod string

const (
	StorageMethodInMemory StorageMethod = "in_memory"
	StorageMethodFile                   = "file"
)

type StorageFactory func() Storage

var supportedStorageMethods = map[StorageMethod]StorageFactory{
	StorageMethodInMemory: NewInMemoryStorage,
	StorageMethodFile:     NewFileStorage,
}

type Storage interface {
	WordStorage
	GameStorage
}

type WordStorage interface {
	LoadWordList() ([]models.Word, error)
	// This method adds a new word to the global list of words.
	// Return values:
	// - bool: defines whether the word which should be added already exists
	// - models.Word: If the word which should be added already exists then it is returned in this value
	// - error: In here you will find answers if something goes south.
	AddWord(models.Word) (bool, models.Word, error)
}

type GameStorage interface {
	LoadCurrentGame() (models.Game, bool, error)
	SaveGame(models.Game) error
}

func NewStorage(storageMethod string) Storage {
	factory, found := supportedStorageMethods[StorageMethod(storageMethod)]
	if !found {
		logger.Info("Storage method %q not found. Fallback to in memory", storageMethod)
		return NewInMemoryStorage()
	}

	return factory()
}
